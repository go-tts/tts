# go-tts

Converts text-to-speech using google translate API, so it needs an internet connection to do so.

The application can get text stream from stdin and create audio stream into stdout.

Use `.` or `\n` (new line) to over a statement, including the last one.

<h1>Installation</h1>

```
go install github.com/go-tts/tts/cmd/tts@latest
```
This command will build `tts` tool on your machine, so you could use it as `tts` further.

or run from source code:

```
git clone github.com/go-tts/tts
cd tts
go run cmd/tts/main.go "Hello world"
```

<h2>Play from terminal</h2>

```
tts "Hello world"
```
Such call suppose to pronounce a text passed to `tts`. If `-l` is not defined, `en-US` language will be used for pronunciation.

<h2>Write to file</h2>

```
tts -l=en -i=text_file.txt -o=audio_file.mp3
tts -l en-US < text_file.txt > audio_file.mp3
```

<h2>Use with pipe</h2>

```
tts -l it "Muy bien. Chao." | ffmpeg -i - -filter:a "atempo=1.5" audio_file.mp3
```

As an example, here `ffmpeg` increases the speech speed by 1.5x.

<h2>Use in your projects</h2>

```
import "github.com/go-tts/tts/pkg/speech"
```

Just import speech package and use it's functions.

```
audioIn, err := speech.FromText(text, speech.LangEn)
```
```
audioIn := speech.FromTextStream(textIn, speech.LangUs)
```
```
err := speech.WriteToAudioStream(textIn, audioOut, "it")
```
