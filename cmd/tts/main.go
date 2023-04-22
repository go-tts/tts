package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-tts/tts/pkg/audio"
	"github.com/go-tts/tts/pkg/speech"
)

func main() {
	in, player, lang := readParams()

	speech := speech.FromTextStream(in, lang)
	player.Play(speech)
}

func readParams() (io.Reader, audio.Player, string) {
	var (
		inputFile, outputFile, lang string
		help                        bool
	)
	langDescription := fmt.Sprintf("language, '%s' by default", speech.LangUs)
	flag.StringVar(&inputFile, "i", "", "input file")
	flag.StringVar(&inputFile, "input", "", "input file")
	flag.StringVar(&outputFile, "o", "", "output file name, reads text if empty")
	flag.StringVar(&outputFile, "output", "", "output file name, reads text if empty")
	flag.StringVar(&lang, "l", speech.LangUs, langDescription)
	flag.StringVar(&lang, "lang", speech.LangUs, langDescription)
	flag.BoolVar(&help, "h", false, "Help")
	flag.BoolVar(&help, "help", false, "Help")

	flag.Parse()

	if help {
		fmt.Printf(`Usage: %s\n "Hello world"`, os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	return inputTextReader(inputFile, flag.Args()), outputAudioPlayer(outputFile), lang
}

func exit(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func generateOutputFileName() string {
	return fmt.Sprintf("speech-%s.mp3", time.Now().Format("2006-01-02-15:04:05"))
}

func inputTextReader(inputFile string, args []string) io.ReadCloser {
	switch {
	case inputFile != "":
		f, err := os.Open(inputFile)
		if err != nil {
			exit("sorry, failed to read an input file %q: %s", inputFile, err.Error())
		}
		return f
	case len(args) > 0:
		return io.NopCloser(strings.NewReader(strings.Join(args, " ")))
	default:
		return os.Stdin
	}
}

func outputAudioPlayer(outputFile string) audio.Player {
	if outputFile != "" {
		dir, _ := path.Split(outputFile)
		if dir != "" {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				exit("sorry, failed to create directories by file path %q: %s", outputFile, err.Error())
			}
		}

		info, err := os.Stat(outputFile)
		if err != nil && !os.IsNotExist(err) {
			exit("sorry, failed to check an output file %q: %s", outputFile, err.Error())
		}

		if info != nil && info.IsDir() {
			outputFile = path.Join(outputFile, generateOutputFileName())
		}

		outFile, err := os.Create(outputFile)
		if err != nil {
			exit("sorry, failed to create an output file %q: %s", outputFile, err.Error())
		}
		return audio.NewRecorder(outFile)
	} else {
		stdoutStat, err := os.Stdout.Stat()
		if err != nil {
			exit("sorry, failed to check stdout stream: %s", err.Error())
		}
		if (stdoutStat.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
			return audio.NewSpeaker()
		}
	}
	return audio.NewRecorder(os.Stdout)
}
