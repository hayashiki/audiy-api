package slack

import (
	"bytes"
	"io"
	"log"

	"github.com/slack-go/slack"
)

//go:generate mockgen -source ./slack.go -destination ./mock/mock_slack.go
type Service interface {
	Upload(title, name, channel, ts string, r io.Reader) error
	Download(url string, b *bytes.Buffer) error
}

type client struct {
	cli *slack.Client
}

func NewClient(token string) Service {
	return &client{
		cli: slack.New(token),
	}
}

func (c *client) Download(url string, b *bytes.Buffer) error {

	err := c.cli.GetFile(url, b)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (c *client) Upload(title, name, channel, ts string, r io.Reader) error {

	params := slack.FileUploadParameters{
		Title:  title,
		File:   name,
		Reader: r,
		//InitialComment: TODO: need??,
		Channels:        []string{channel},
		ThreadTimestamp: ts,
	}

	_, err := c.cli.UploadFile(params)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
