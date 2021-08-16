package usecase

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/hayashiki/audiy-api/domain/entity"
	"github.com/hayashiki/audiy-api/infrastructure/gcs"
	"github.com/hayashiki/audiy-api/infrastructure/slack"
)

type IntegrateAudioUsecase interface {
	Do(ctx context.Context, input *AudioInput) error
}

type integrateAudioUsecase struct {
	slackSvc  slack.Slack
	radioRepo entity.AudioRepository
	gcsSvc    gcs.Client
}

type MockAudioUsecase struct {
	WantErr bool
}

func (mock MockAudioUsecase) Do(ctx context.Context, input *AudioInput) error {
	if mock.WantErr == true {
		return errors.New("usecase error")
	}
	return nil
}

type AudioInput struct {
	Name               string `json:"name"`
	ID                 string `json:"id"`
	Title              string `json:"title"`
	URLPrivateDownload string `json:"url_private_download"`
	Created            int64  `json:"created"`
	Mimetype           string `json:"mimetype"`
}

func (i AudioInput) Validate() error {
	if i.URLPrivateDownload == "" {
		return errors.New("empty download url")
	}
	return nil
}

func NewAudio(
	slackSvc slack.Slack,
	radioRepo entity.AudioRepository,
	gcsSvc gcs.Client,
) IntegrateAudioUsecase {
	return &integrateAudioUsecase{
		slackSvc:  slackSvc,
		radioRepo: radioRepo,
		gcsSvc:    gcsSvc,
	}
}

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// Do is ラジオ情報を保存して、コンバートしてストレージに保存
func (au *integrateAudioUsecase) Do(ctx context.Context, input *AudioInput) error {

	// TODO:
	// input.IDがrepoにある、storageにある場合はスキップする
	//fileName := filepath.Base(payload.File.URLPrivateDownload)
	b := bytes.Buffer{}
	err := au.slackSvc.Download(input.URLPrivateDownload, &b)
	if err != nil {
		log.Printf("failed to get a slack file err=%v", err)
		return err
	}
	ext := filepath.Ext(input.Name)
	if ext != ".m4a" {
		fmt.Println("must be file .m4a: " + input.Name)
		return nil
	}
	if err := au.gcsSvc.Put(ctx, fmt.Sprintf("%s%s", input.ID, ext), b.Bytes()); err != nil {
		log.Printf("failed to put gcs client")
		return err
	}

	//data, err := au.gcsSvc.Get(ctx, input.Name)
	//if err != nil || data == nil {
	//	return fmt.Errorf("failed to read gcsSvc data %w", err)
	//}
	// file check extension only m4a

	//size := getSize(data)
	getFilePath := getFilePath(au.gcsSvc.Bucket(), fmt.Sprintf("%s%s", input.ID, ext))
	ut := time.Unix(input.Created, 0)
	newRadio := entity.NewAudio(input.ID, input.Name, int(100), getFilePath, input.Mimetype, ut)
	err = au.radioRepo.Save(ctx, newRadio)
	log.Printf("newRadio %+v", newRadio.GetKey())
	if err != nil {
		return fmt.Errorf("fail to create radios record err: %w", err)
	}
	return nil
}

func getSize(data []byte) int32 {
	return int32(binary.Size(data))
}

func getFilePath(bucket, name string) string {
	return "https://storage.cloud.google.com/" + bucket + "/" + name + "?authuser=1"
}
