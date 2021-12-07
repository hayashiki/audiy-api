package config

import (
	"cloud.google.com/go/compute/metadata"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GCSInputAudioBucket string `envconfig:"GCS_INPUT_AUDIO_BUCKET" required:"true"`
	SlackBotToken       string `envconfig:"SLACK_BOT_TOKEN" required:"true"`
	IsDev               bool   `envconfig:"IS_DEV" default:"true"`
}

func NewConfig() (Config, error) {
	env := Config{}
	err := envconfig.Process("", &env)
	return env, err
}

// GetProject on Google Cloud
func GetProject() string {
	var (
		project string
		err     error
	)
	if project = os.Getenv("GCP_PROJECT"); project == "" {
		project, err = metadata.ProjectID()
		if err != nil || project == "" {
			log.Fatal("project id can't be empty")
		}
	}

	return project
}
