package transcript

import (
	"bytes"
	speech "cloud.google.com/go/speech/apiv1p1beta1"
	"context"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1p1beta1"
	"io"
)

const (
	AudioLang = "ja-JP"
	AudioRateHertz = 44100
	RecognitionConfig = speechpb.RecognitionConfig_MP3
)

type SpeechRecogniser interface {
	Recognize(ctx context.Context, file io.Reader) (*speechpb.LongRunningRecognizeResponse, error)
}

type speechRecogniser struct {
	client *speech.Client
}

func NewSpeechRecogniser() *speechRecogniser {
	return &speechRecogniser{}
}

func (s *speechRecogniser) new(ctx context.Context) (*speech.Client, error) {
	return speech.NewClient(ctx)
}

func (s *speechRecogniser) Recognize(ctx context.Context, file io.Reader) (*speechpb.LongRunningRecognizeResponse, error) {
	client, err := s.new(ctx)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return nil, err
	}

	req := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:                   RecognitionConfig,
			SampleRateHertz:            AudioRateHertz,
			LanguageCode:               AudioLang,
			EnableAutomaticPunctuation: true,
			EnableWordTimeOffsets: true,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: buf.Bytes()},
		},
	}

	op, err := client.LongRunningRecognize(ctx, req)
	if err != nil {
		return nil, err
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return nil, err
	}
	//result := NewTranscriptionResult(resp)
	return resp, nil
}
