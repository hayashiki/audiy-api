package config

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"

	"cloud.google.com/go/compute/metadata"
)

type SlackConf struct {
	GCSInputAudioBucket string `envconfig:"GCS_INPUT_AUDIO_BUCKET" required:"true"`
	SlackBotToken       string `envconfig:"SLACK_BOT_TOKEN" required:"true"`
}

func NewConf() (SlackConf, error) {
	env := SlackConf{}
	err := envconfig.Process("", &env)
	return env, err
}

// GetProject on Google Cloud
func GetProject() string {
	var (
		project string
		err     error
	)

	project, err = metadata.ProjectID()
	if err != nil {
		if project = os.Getenv("GCP_PROJECT"); project == "" {
			log.Fatal("project id can't be empty")
		}
	}

	return project
}
