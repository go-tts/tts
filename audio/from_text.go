package audio

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func FromText(textIn io.Reader, audioOut io.Writer, lang string) error {
	input := bufio.NewScanner(textIn)
	for input.Scan() {
		audioIn, err := readText(input.Text(), lang)
		if err != nil {
			return err
		}

		if err := Append(audioIn, audioOut); err != nil {
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
