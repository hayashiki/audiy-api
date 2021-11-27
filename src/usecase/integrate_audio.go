package usecase

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/hayashiki/audiy-api/src/infrastructure/ffmpeg"
	"log"
	"path/filepath"
	"time"

	entity "github.com/hayashiki/audiy-api/src/domain/entity"
	gcs "github.com/hayashiki/audiy-api/src/infrastructure/gcs"
	slack "github.com/hayashiki/audiy-api/src/infrastructure/slack"
)

type IntegrateAudioUsecase interface {
	Do(ctx context.Context, input *AudioInput) error
}

type integrateAudioUsecase struct {
	slackSvc  slack.Service
	gcsSvc    gcs.Service
	proverSvc ffmpeg.Service
	audioRepo entity.AudioRepository
	feedRepo  entity.FeedRepository
	userRepo  entity.UserRepository
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
	slackSvc slack.Service,
	gcsSvc gcs.Service,
	proverSvc ffmpeg.Service,
	audioRepo entity.AudioRepository,
	feedRepo entity.FeedRepository,
	userRepo entity.UserRepository,
) IntegrateAudioUsecase {
	return &integrateAudioUsecase{
		slackSvc:  slackSvc,
		gcsSvc:    gcsSvc,
		proverSvc: proverSvc,
		audioRepo: audioRepo,
		feedRepo:  feedRepo,
		userRepo:  userRepo,
	}
}

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// Do is audioを保存して、コンバートしてストレージに保存
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
	r := bytes.NewReader(b.Bytes())

	if err := au.gcsSvc.Write(ctx, fmt.Sprintf("%s%s", input.ID, ext), r); err != nil {
		log.Printf("failed to put gcs client")
		return err
	}

	newFile, err := au.gcsSvc.Read(ctx, fmt.Sprintf("%s%s", input.ID, ext))
	if err != nil {
		log.Printf("failed to read gcs client")
		return err
	}

	data, err := au.proverSvc.GetProbe(ctx, newFile)
	log.Println("ffprobe data", data)
	if err != nil {
		return err
	}
	// file check extension only m4a
	getFilePath := getFilePath(au.gcsSvc.Bucket(), fmt.Sprintf("%s%s", input.ID, ext))
	ut := time.Unix(input.Created, 0)
	newAudio := entity.NewAudio(input.ID, input.Name, data.Format.DurationSeconds, getFilePath, input.Mimetype, ut)
	err = au.audioRepo.Save(ctx, newAudio)
	log.Printf("newAudio %+v", newAudio.GetKey())
	if err != nil {
		return fmt.Errorf("fail to create radios record err: %w", err)
	}

	users, _ := au.userRepo.GetAll(ctx)
	feeds := make([]*entity.Feed, len(users))
	userIDs := make([]string, len(users))
	newFeed := entity.NewFeed(newAudio.Key.Name, newAudio.PublishedAt)
	newFeed.PublishedAt = newAudio.PublishedAt

	for i, u := range users {
		userIDs[i] = u.ID
		feeds[i] = newFeed
	}
	au.feedRepo.SaveAll(ctx, userIDs, feeds)

	return nil
}

func getSize(data []byte) int32 {
	return int32(binary.Size(data))
}

func getFilePath(bucket, name string) string {
	return "https://storage.cloud.google.com/" + bucket + "/" + name + "?authuser=1"
}
