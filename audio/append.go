package audio

import (
	"errors"
	"io"

	"github.com/dmulholl/mp3lib"
)

func Append(in io.Reader, out io.Writer) error {
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
