package main

import (
	"flag"
	"os"

	"automation/tts/audio"
)

func main() {
	var lang string
	flag.StringVar(&lang, "l", "ru", "language")
	flag.Parse()

	textIn := os.Stdin
	audioOut := os.Stdout

	if err := audio.FromText(textIn, audioOut, lang); err != nil {
		os.Stderr.WriteString("failed to read text: " + err.Error())
	}
}
