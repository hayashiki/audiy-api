package ffmpeg

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"gopkg.in/vansante/go-ffprobe.v2"
	"io"
	"log"
	"os"
	"os/exec"
)

func init() {
	_, err := exec.LookPath("ffprobe")

	if err != nil {
		log.Fatal("Could not find ffprobe. Please install ffmpeg and ffprobe.")
	}
}

type Service struct {
	inputPath string
	outputPath string
}

func (s *Service) GetProbe(ctx context.Context, r io.Reader) (*ffprobe.ProbeData, error) {
	tempFileName := uuid.New().String()
	tmpFile, err := os.CreateTemp(os.TempDir(), tempFileName)
	defer os.Remove(tmpFile.Name())
	tmpFile.ReadFrom(r)

	// ffprove.ProbeReader()をつかって、引数のio.Readerをよみこませたかったが、m4a形式だとdurationがとれない
	data, err := ffprobe.ProbeURL(ctx, tmpFile.Name())
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) probeDebug(ctx context.Context, r io.Reader) error {
	probe, err := s.GetProbe(ctx, r)
	if err != nil {
		return err
	}
	buf, err := json.MarshalIndent(probe, "", "  ")
	if err != nil {
		log.Panicf("Error unmarshalling: %v", err)
	}
	log.Print(string(buf))
	return nil
}
