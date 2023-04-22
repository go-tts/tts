package speech

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

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
	input := bufio.NewScanner(textIn)
	for input.Scan() {
		text := input.Text()
		if text == "" {
			continue
		}

		audioIn, err := FromText(input.Text(), lang)
		if err != nil {
			return err
		}

		if _, err := io.Copy(audioOut, audioIn); err != nil {
			audioIn.Close()
			return err
		}

		if err := audioIn.Close(); err != nil {
			return err
		}
	}

	return nil
}

func FromText(text, lang string) (io.ReadCloser, error) {
	url := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=%d&client=tw-ob&q=%s&tl=%s", len(text), url.QueryEscape(text), lang)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return response.Body, err
}
