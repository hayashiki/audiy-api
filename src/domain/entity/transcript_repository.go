package entity

import "context"

// TranscriptRepository interface
type TranscriptRepository interface {
	Save(context.Context, *Transcript) error
}

