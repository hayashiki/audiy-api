package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/hayashiki/audiy-api/src/domain/entity"
	"github.com/hayashiki/audiy-api/src/infrastructure/ffmpeg"
	"github.com/hayashiki/audiy-api/src/infrastructure/gcs"
	"github.com/hayashiki/audiy-api/src/infrastructure/transcript"
)

type TranscriptAudioUsecase interface {
	Do(ctx context.Context, input entity.CreateTranscriptInput) error
}

type transcriptAudioUsecase struct {
	gcsSvc         gcs.Service
	audioRepo      entity.AudioRepository
	transcriptRepo entity.TranscriptRepository
	ffmpegProveSvc ffmpeg.Service
	transcoder     ffmpeg.Transcoder
	transcriptSvc  transcript.SpeechRecogniser
}

func (t transcriptAudioUsecase) Do(ctx context.Context, input entity.CreateTranscriptInput) error {
	audio, err := t.audioRepo.Find(ctx, input.AudioID)
	if err != nil {
		return nil
	}
	// TODO: ".m4a" handle
	audioReader, err := t.gcsSvc.Read(ctx, input.AudioID+".m4a")
	if err != nil {
		fmt.Errorf("failed to read gcs %w", err)
		return err
	}
	// TODO: mp4だったらconvertする
	convOutput, convProgress, err := t.transcoder.Transcode(audioReader)
	if err != nil {
		fmt.Errorf("failed to convert mp3: %w", err)
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// TODO: ".mp3" handle
		err := t.gcsSvc.Write(ctx, input.AudioID+".mp3", convOutput)
		if err != nil {
			fmt.Errorf("failed to write to gcs %w", err)
			return
		}
		// transcribe
		ts, err := t.transcriptSvc.RecognizeGCS(ctx, "gs://"+t.gcsSvc.Bucket()+"/"+input.AudioID+".mp3")
		if err != nil {
			fmt.Errorf("failed to recognize to gcs %w", err)
			return
		}
		newTranscript := entity.NewTranscript(audio.ID, ts)
		if err := t.transcriptRepo.Save(ctx, newTranscript); err != nil {
			log.Println("ds save err", err)
			fmt.Errorf("failed to transcribe to gcs %w", err)
			return
		}
		convOutput.Close()
	}()
	convProgress.Wait()
	wg.Wait()
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
