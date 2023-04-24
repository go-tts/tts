package speech

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// queryTextLimit is the limit the query to google handles fine.
// Usually response from google is Bad Request if the text size is above this limit.
const queryTextLimit = 200

func FromText(text, lang string) (io.ReadCloser, error) {
	url := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&textlen=%d&client=tw-ob&q=%s&tl=%s", len(text), url.QueryEscape(text), lang)
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to query google: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		bodyText := readBody(response.Body)
		return nil, fmt.Errorf("failed to query google: response status %d - %s: %s", response.StatusCode, http.StatusText(response.StatusCode), bodyText)
	}
	return response.Body, nil
}

func FromTextStream(textIn io.Reader, lang string) io.ReadCloser {
	r, w := io.Pipe()
	go func() {
		var err error
		defer func() {
			w.CloseWithError(err)
		}()
		err = WriteToAudioStream(textIn, w, lang)
	}()
	return r
}

func WriteToAudioStream(textIn io.Reader, audioOut io.Writer, lang string) error {
	readBuf := make([]byte, queryTextLimit)
	var textBuf []byte
	for {
		n, err := textIn.Read(readBuf)
		if errors.Is(err, io.EOF) {
			for len(textBuf) > 0 {
				textBuf = trimSpaces(textBuf)
				textBuf, err = writeToAudioStream(textBuf, audioOut, lang)
				if err != nil {
					return fmt.Errorf("failed to write audio: %w", err)
				}
			}
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to read input text: %w", err)
		}

		textBuf = append(textBuf, readBuf[:n]...)
		textBuf = trimSpaces(textBuf)
		for len(textBuf) >= queryTextLimit {
			textBuf, err = writeToAudioStream(textBuf, audioOut, lang)
			if err != nil {
				return fmt.Errorf("failed to write audio: %w", err)
			}
			textBuf = trimSpaces(textBuf)
		}
	}
}

func writeToAudioStream(textBuf []byte, audioOut io.Writer, lang string) ([]byte, error) {
	if len(textBuf) == 0 {
		return textBuf, nil
	}

	var text string
	textBuf, text = scanText(textBuf)

	audioIn, err := FromText(text, lang)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(audioOut, audioIn); err != nil {
		audioIn.Close()
		return nil, err
	}

	if err := audioIn.Close(); err != nil {
		return nil, err
	}
	return textBuf, nil
}

// readBody returns body from response if it exists,
// or error if there is an error.
func readBody(responseBody io.Reader) string {
	body, err := io.ReadAll(responseBody)
	if errors.Is(err, io.EOF) {
		return "<empty response>"
	}
	if err != nil {
		return err.Error()
	}
	if len(body) == 0 {
		return "<empty response>"
	}
	return string(body)
}

// scanText scans what part of buffer can be sent to request,
// and which part needs to be added before sending.
func scanText(buf []byte) ([]byte, string) {
	if len(buf) < queryTextLimit {
		text := string(buf)
		return buf[:0], text
	}

	separatorPlace := -1
	for i := queryTextLimit - 1; i >= 0; i-- {
		switch buf[i] {
		case '.', ',', '!', '?', ';', ':':
			separatorPlace = i
			goto found
		case '\n':
			// exclude new line byte
			separatorPlace = i - 1
			if i > 1 && buf[i-2] == '\r' {
				separatorPlace = i - 2
			}
			goto found
		case ' ':
			if separatorPlace < 0 {
				// case, when there is no sign of statement over
				separatorPlace = i
			}
		}
	}
found:
	if separatorPlace < 0 {
		// case when buffer contains one huge word
		separatorPlace = queryTextLimit
	}
	text := string(buf[:separatorPlace+1])
	copy(buf, buf[separatorPlace+1:])
	return buf[:len(buf)-separatorPlace-1], text
}

// trimSpaces removes spaces in the beginning of the buffer
// to avoid sending spaces via the internet.
func trimSpaces(buf []byte) []byte {
	cutFrom := -1
	for i := 0; i < len(buf); i++ {
		switch buf[i] {
		case ' ', '\r', '\n':
			cutFrom = i + 1
		default:
			goto out
		}
	}
out:
	if cutFrom < 0 {
		return buf
	}
	if cutFrom == len(buf) {
		return nil
	}
	return buf[cutFrom:]
}
