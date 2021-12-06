package transcript

import (
	"context"
	"os"
	"testing"
)

func TestNewSpeechRecogniser(t *testing.T) {
	t.Skip()
	f, err := os.Open("testdata/audio_voice_sample.mp3")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	ctx := context.Background()
	spCli := NewSpeechRecogniser()
	transcriptResp, err := spCli.Recognize(ctx, f)
	if err != nil {
		t.Error(err)
	}
	t.Log(transcriptResp)
//	私たちは1990年という IT 分野の幕開けともいえる時期から it ビジネス支援事業としてデジタルコンテンツを中心としたユニークなサービスを提供しています現在数多くの実績を誇る国内屈指のインフォメーションインテグレータで企業に求められるさまざまな要件について、顧客企業様の経営戦略に基づく包括的なサービスを提供しています。
}
