package ffmpeg

import (
	"context"
	"github.com/hayashiki/audiy-api/src/infrastructure/gcs"
	"gopkg.in/vansante/go-ffprobe.v2"
	"io"
	"log"
	"os"
	"testing"
)

func TestTranscoder_Transcode(t *testing.T) {
	//inputFile, _ := os.Open("/testdata/audio_voice_sample.m4a")
	//defer inputFile.Close()

	ctx := context.Background()

	gcsSvc := gcs.NewService("bulb-audiy-audio-bucket")
	input, err := gcsSvc.Read(ctx, "F01L6BH32JH.m4a")
	if err != nil {
		panic(err)
	}


	outputFile, _ := os.Create("test_clip.mp3")
	//defer outputFile.Close()

	ts := Transcoder{}
	convOutput, convProgress, err := ts.Transcode(input)
	if err != nil {
		t.Errorf("testMp3ToWav failed due to an error: %v", err)
		return
	}

	go func() {
		svc := gcs.NewService("staging.bulb-audiy.appspot.com")

		data, err := ffprobe.ProbeReader(context.Background(), convOutput)
		if err != nil {
			t.Error(err)
		}
		log.Println(data.Format.FormatName)

		if err := svc.Write(context.Background(), "out9.mp3", convOutput); err != nil {
			t.Error(err)
		}
		_, err = io.Copy(outputFile, convOutput)
		if err != nil {
			t.Errorf("failed to copy audio from converter into output: %v", err)
		}

		convOutput.Close()
	}()
	convProgress.Wait()

	svc := gcs.NewService("staging.bulb-audiy.appspot.com")
	if err := svc.Write(context.Background(), "out8.mp3", outputFile); err != nil {
		t.Error(err)
	}

}
