package usecase

import (
	"context"
	"fmt"
	"github.com/hayashiki/audiy-api/src/domain/entity"
	"github.com/hayashiki/audiy-api/src/infrastructure/ffmpeg"
	"github.com/hayashiki/audiy-api/src/infrastructure/gcs"
	"github.com/hayashiki/audiy-api/src/infrastructure/transcript"
	"gopkg.in/vansante/go-ffprobe.v2"
	"log"
)

type TranscriptAudioUsecase interface {
	Do(ctx context.Context, input *AudioInput) error
}

type transcriptAudioUsecase struct {
	gcsSvc         gcs.Service
	audioRepo      entity.AudioRepository
	transcriptRepo entity.TranscriptRepository
	ffmpegProveSvc ffmpeg.Service
	transcoder     ffmpeg.Transcoder
	transcriptSvc  transcript.SpeechRecogniser
}

func (t transcriptAudioUsecase) Do(ctx context.Context, input *AudioInput) error {
	audio, err := t.audioRepo.Find(ctx, input.ID)
	if err != nil {
		return nil
	}
	audioReader, err := t.gcsSvc.Read(ctx, input.ID)
	if err != nil {
		return err
	}
	// TODO: mp4だったらconvertする
	convOutput, convProgress, err := t.transcoder.Transcode(audioReader)
	if err != nil {
		fmt.Errorf("testMp3ToWav failed due to an error: %v", err)
		return err
	}
	go func() {
		data, err := ffprobe.ProbeReader(context.Background(), convOutput)
		if err != nil {
			panic(err)
		}
		log.Println(data.Format.FormatName)

		// transcribe
		ts, err := t.transcriptSvc.Recognize(ctx, convOutput)
		if err != nil {
			return
		}
		newTranscript := entity.NewTranscript(audio.ID, ts)
		if err := t.transcriptRepo.Save(ctx, newTranscript); err != nil {
			return
		}
		convOutput.Close()
	}()
	convProgress.Wait()

	return nil
}

func NewTranscriptAudioUsecase(
	gcsSvc gcs.Service,
	audioRepo entity.AudioRepository,
	transcriptRepo entity.TranscriptRepository,
	proveSvc ffmpeg.Service,
	transcoderSvc ffmpeg.Transcoder,
	transcriptSvc transcript.SpeechRecogniser,

) TranscriptAudioUsecase {
	return &transcriptAudioUsecase{
		gcsSvc:         gcsSvc,
		audioRepo:      audioRepo,
		transcriptRepo: transcriptRepo,
		ffmpegProveSvc: proveSvc,
		transcoder:     transcoderSvc,
		transcriptSvc:  transcriptSvc,
	}
}

