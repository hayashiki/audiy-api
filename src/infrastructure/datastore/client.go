package datastore

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/hayashiki/audiy-api/src/config"
	mdatastore "go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
	"go.mercari.io/datastore/clouddatastore"
	"google.golang.org/api/option"
	"log"
	"os"
)

func NewClient(ctx context.Context, projectID string, options ...option.ClientOption) (*datastore.Client, error) {
	var opts []option.ClientOption

	if os.Getenv("DATASTORE_EMULATOR_HOST") == "" {
		if credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); credentialsFile != "" {
			opts = append(opts, option.WithCredentialsFile(credentialsFile))
		}
	}

	opts = append(opts, options...)

	return datastore.NewClient(ctx, projectID, opts...)
}

func FromContext(ctx context.Context) *boom.Boom {
	var opts []option.ClientOption

	if os.Getenv("DATASTORE_EMULATOR_HOST") == "" {
		if credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); credentialsFile != "" {
			opts = append(opts, option.WithCredentialsFile(credentialsFile))
		}
	}

	cli, err := datastore.NewClient(ctx, config.GetProject(), opts...)
	if err != nil {
		log.Println("cli", err)
		panic(err)
	}
	ds, err := clouddatastore.FromClient(ctx, cli)
	if err != nil {
		panic(err)
	}
	return boom.FromClient(ctx, ds)
}

func FromContext2(ctx context.Context) (*boom.Boom, mdatastore.Client) {
	var opts []option.ClientOption

	if os.Getenv("DATASTORE_EMULATOR_HOST") == "" {
		if credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); credentialsFile != "" {
			opts = append(opts, option.WithCredentialsFile(credentialsFile))
		}
	}

	cli, err := datastore.NewClient(ctx, config.GetProject(), opts...)
	if err != nil {
		log.Println("cli", err)
		panic(err)
	}
	ds, err := clouddatastore.FromClient(ctx, cli)

	if err != nil {
		panic(err)
	}
	return boom.FromClient(ctx, ds), ds
}
