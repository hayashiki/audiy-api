package transcript

import (
	"strings"

	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1p1beta1"
)

type Transcript struct {
	Monologues []Monologue
	Body       string
}

type Monologue struct {
	Elements []MonologueElement
}

type MonologueElement struct {
	StartTime  float64
	EndTime    float64
	Word       string
	WordKana   string
	Confidence float32
}

func NewTranscriptionResult(resp *speechpb.LongRunningRecognizeResponse) *Transcript {
	t := Transcript{}

	for _, res := range resp.Results {

		if len(res.Alternatives) == 0 {
			continue
		}
		alt := res.Alternatives[0]

		ws := make([]MonologueElement, len(alt.Words))
		for i, w := range alt.Words {
			// TODO: Wordがpipeされないケースはある？
			parts := strings.Split(w.Word,
				"|")
			ws[i] = MonologueElement{
				Word:      parts[0],
				WordKana:  parts[1],
				StartTime: float64(w.StartTime.Seconds) + float64(w.StartTime.Nanos)*1e-9,
				EndTime:   float64(w.EndTime.Seconds) + float64(w.EndTime.Nanos)*1e-9,
			}
		}

		t.Body = t.Body + alt.Transcript
		m := Monologue{Elements: ws}
		t.Monologues = append(t.Monologues, m)
	}

	return &t
}
