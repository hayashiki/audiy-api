package ffmpeg

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"os"

	//"log"
	"os/exec"
	"sync"
)

type Transcoder struct {}

func (t *Transcoder) Transcode(input io.Reader) (output io.ReadCloser, done *sync.WaitGroup, err error) {
	tempFileName := uuid.New().String()
	tmpFile, err := os.CreateTemp(os.TempDir(), tempFileName)
	// TODO: どうにかする
	//defer os.Remove(tmpFile.Name())
	tmpFile.ReadFrom(input)

	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, nil, fmt.Errorf("ffmpeg was not found in PATH. Please install ffmpeg")
	}

	// TODO: 要調整
	cmd := exec.Command(ffmpeg,
		"-i", tmpFile.Name(),
		//"-c:a", "mp3",
		//"-acodec", "copy",
		//"-movflags", "faststart",
		//"-c", "copy",
		//"-bsf:a", "aac_adtstoasc",
		"-f", "mp3",
		"-",
	)
	output, err = cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to pipe output from audio converter, err: %v", err)
	}
	cmd.Stderr = os.Stderr

	done = &sync.WaitGroup{}
	done.Add(1)
	if err = cmd.Start(); err != nil { //Use start, not run
		return nil, nil, fmt.Errorf("failed to start conversion, err: %v", err)
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("ffmpeg encountered an error while decoding: %v\n", err)
		}
		done.Done()
	}()

	return output, done, nil
}
