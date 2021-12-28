package model

import (
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1p1beta1"
	"log"
	"strings"
	"time"
)

const TranscriptKind = "Transcript"

type Transcript struct {
	ID         int64          `json:"id"`
	//AudioKey   *datastore.Key `json:"audio_key"`
	AudioID   string `json:"audio_key" datastore:"audio_key"`
	Body       string         `json:"body"`
	Monologues []Monologue `json:"monologues"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Monologue struct {
	Elements []MonologueElement `json:"elements" datastore:"elements"`
}

type MonologueElement struct {
	StartTime  float64 `json:"start_time" datastore:"start_time"`
	EndTime    float64 `json:"end_time" datastore:"end_time"`
	Word       string `json:"word" datastore:"word"`
	WordKana   string `json:"word_kana" datastore:"word_kana"`
	Confidence float32 `json:"confidence" datastore:"confidence"`
}

func (Transcript) IsNode() {}

func NewTranscript(audioID string, resp *speechpb.LongRunningRecognizeResponse) *Transcript {
	t := &Transcript{
		ID:        NewID(TranscriptKind),
		AudioID:   audioID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, res := range resp.Results {
		if len(res.Alternatives) == 0 {
			continue
		}
		alt := res.Alternatives[0]

		ws := make([]MonologueElement, len(alt.Words))
		for i, w := range alt.Words {
			// TODO: Wordがpipeされないケースはある？
			parts := strings.Split(w.Word, "|")

			if len(parts) != 2 {
				log.Println(parts)
			}

			ws[i] = MonologueElement{
				Word:      parts[0],
				//WordKana:  parts[1],
				StartTime: float64(w.StartTime.Seconds) + float64(w.StartTime.Nanos)*1e-9,
				EndTime:   float64(w.EndTime.Seconds) + float64(w.EndTime.Nanos)*1e-9,
			}
		}

		t.Body = t.Body + alt.Transcript
		m := Monologue{Elements: ws}
		t.Monologues = append(t.Monologues, m)
	}

	return t
}
