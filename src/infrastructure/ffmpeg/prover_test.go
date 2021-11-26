package ffmpeg

import (
	"context"
	"os"
	"testing"
)

func TestService_Get(t *testing.T) {
	ctx := context.Background()
	dir, _ := os.Getwd()
	f, err := os.Open(dir +"/testdata/file_example_MP3_700KB.mp3")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	svc := Service{}
	data, err := svc.GetProbe(ctx, f)
	if err != nil {
		t.Error(err)
	}
	if actual, want := data.Format.DurationSeconds, 27.252; actual != want {
		t.Errorf("want: %v, actual: %v", actual, want)
	}
}
