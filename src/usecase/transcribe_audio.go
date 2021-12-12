package usecase

import (
	"context"
	"fmt"
	"github.com/hayashiki/audiy-api/src/domain/repository"
	"log"
	"sync"

	"github.com/hayashiki/audiy-api/src/domain/model"
	"github.com/hayashiki/audiy-api/src/infrastructure/ffmpeg"
	"github.com/hayashiki/audiy-api/src/infrastructure/gcs"
	"github.com/hayashiki/audiy-api/src/infrastructure/transcript"
)

type TranscriptAudioUsecase interface {
	Do(ctx context.Context, input model.CreateTranscriptInput) error
}

type transcriptAudioUsecase struct {
	gcsSvc         gcs.Service
	audioRepo      repository.AudioRepository
	transcriptRepo repository.TranscriptRepository
	ffmpegProveSvc ffmpeg.Service
	transcoder     ffmpeg.Transcoder
	transcriptSvc  transcript.SpeechRecogniser
}

func (t transcriptAudioUsecase) Do(ctx context.Context, input model.CreateTranscriptInput) error {
	audio, err := t.audioRepo.Get(ctx, input.AudioID)
	if err != nil {
		return nil
	}
	if audio.Transcribed {
		log.Printf("already transcribed")
		return nil
	}
	log.Println(0)
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
		log.Println(1)
		err := t.gcsSvc.Write(ctx, input.AudioID+".mp3", convOutput)
		if err != nil {
			fmt.Errorf("failed to write to gcs %w", err)
			return
		}
		// transcribe
		log.Println(2, t.gcsSvc.Bucket(), input.AudioID+".mp3")

		ts, err := t.transcriptSvc.RecognizeGCS(ctx, "gs://"+t.gcsSvc.Bucket()+"/"+input.AudioID+".mp3")
		if err != nil {
			log.Println(2.5, err)
			fmt.Errorf("failed to recognize to gcs %w", err)
			return
		}
		// TODO: transaction
		log.Println(3)
		newTranscript := model.NewTranscript(audio.ID, ts)
		if err := t.transcriptRepo.Put(ctx, newTranscript); err != nil {
			log.Println("ds save err", err)
			fmt.Errorf("failed to transcribe to gcs %w", err)
			return
		}
		log.Println(4)
		audio.SetTranscribed()
		if err := t.audioRepo.Put(ctx, audio); err != nil {
			fmt.Errorf("failed to save audio %w", err)
			return
		}
		log.Println(5)
		convOutput.Close()
	}()
	convProgress.Wait()
	wg.Wait()
	return nil
}

func NewTranscriptAudioUsecase(
	gcsSvc gcs.Service,
	audioRepo repository.AudioRepository,
	transcriptRepo repository.TranscriptRepository,
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
