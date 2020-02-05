package speech

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/dmulholl/mp3lib"
)

func FromText(textIn io.Reader, audioOut io.Writer, lang string) error {
	input := bufio.NewScanner(textIn)
	for input.Scan() {
		audioIn, err := readText(input.Text(), lang)
		if err != nil {
			return err
		}

		if err := appendMp3(audioIn, audioOut); err != nil {
			return err
		}

		if err := audioIn.Close(); err != nil {
			return err
		}
	}

	return nil
}

func readText(text, lang string) (io.ReadCloser, error) {
	url := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), lang)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return response.Body, err
}

func appendMp3(in io.Reader, out io.Writer) error {
	isFirstFrame := true

	for {
		frame := mp3lib.NextFrame(in)
		if frame == nil {
			break
		}

		if isFirstFrame {
			isFirstFrame = false
			if mp3lib.IsXingHeader(frame) || mp3lib.IsVbriHeader(frame) {
				continue
			}
		}

		if _, err := out.Write(frame.RawBytes); err != nil {
			return errors.New("failed to write file")
		}
	}

	return nil
}
