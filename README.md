# go-tts
Text-to-speech API. 

Simple application that reads text from stdin and writes mp3 into stdout. It needs an internet connection to reach google translate API.

<h1>Install</h1>

```
go get github.com/go-tts/tts
```

<h2>Use from terminal</h2>

<h3>Play from terminal</h3>

```
tts "Hello world"
```

<h3>Store</h3>

```
tts -l=en -i=text_file.txt -o=audio_file.mp3
tts -l en-US < text_file.txt > audio_file.mp3
```

<h3>Use in pipe</h3>

```
tts -l it "Muy bien. Chao." | ffmpeg -i - -filter:a "atempo=1.5" audio_file.mp3
```

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

<h2>Run on your machine using golang</h2>

```
git clone github.com/go-tts/tts
cd tts
go run cmd/tts/main.go "Hello world"
```