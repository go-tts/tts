package audio

import (
	"io"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

const (
	numOfChannels = 2
	audioBitDepth = 2
)

type Player interface {
	Play(audioIn io.Reader) error
}

func NewRecorder(audioOut io.WriteCloser) Recorder {
	return Recorder{audioOut: audioOut}
}

// Recorder is player that records to io writer: file, stdout etc.
type Recorder struct {
	audioOut io.WriteCloser
}

func (r Recorder) Play(audioIn io.Reader) error {
	_, err := io.Copy(r.audioOut, audioIn)
	return err
}

func NewSpeaker() Speaker {
	return Speaker{}
}

// Speaker is player that plays audio to speakers.
type Speaker struct {
}

func (s Speaker) Play(audioIn io.Reader) error {
	decodedMp3, err := mp3.NewDecoder(audioIn)
	if err != nil {
		return err
	}

	otoCtx, readyCh, err := oto.NewContext(decodedMp3.SampleRate(), numOfChannels, audioBitDepth)
	if err != nil {
		return err
	}
	<-readyCh

	player := otoCtx.NewPlayer(decodedMp3)
	player.Play()

	for player.IsPlaying() {
		time.Sleep(10 * time.Millisecond)
	}

	return player.Close()
}
