package gcs

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/hayashiki/audiy-api/etc/config"

	"google.golang.org/api/option"

	"github.com/pkg/errors"

	"google.golang.org/api/iam/v1"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

type client struct {
	gcsClient *storage.Client
	bucket    string
}

type Client interface {
	Bucket() string
	Put(ctx context.Context, objName string, data []byte) error
}

func NewGCSClient(ctx context.Context, bucket string) (Client, error) {
	gcsClient, err := storage.NewClient(ctx)

	if err != nil {
		panic(err)
	}

	return &client{
		gcsClient: gcsClient,
		bucket:    bucket,
	}, nil
}

func (c *client) Put(ctx context.Context, objName string, data []byte) error {
	w := c.gcsClient.Bucket(c.bucket).Object(objName).NewWriter(ctx)
	defer w.Close()

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

func (c *client) Bucket() string {
	return c.bucket
}

func init() {
	cred, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		log.Printf("failed to initialize the Google client.\n")
		log.Printf("%v\n", errors.WithStack(err.(error)).Error())
		return
	}

	iamService, err = iam.NewService(context.Background(), option.WithHTTPClient(cred))
	if err != nil {
		log.Printf("failed to initialize the IAM.")
		return
	}
}

// ServiceAccountName returns email address format of google service account.
func ServiceAccountName() string {
	return config.GetProject() + "@appspot.gserviceaccount.com"
	//return "audiy-adminapi-sa@bulb-audiy.iam.gserviceaccount.com"
}

// ServiceAccountID returns full account id.
func ServiceAccountID() string {
	return "projects/" + config.GetProject() + "/serviceAccounts/" + ServiceAccountName()
}

var iamService *iam.Service

func StorageObjectFilePath(id string, extension string) string {
	return fmt.Sprintf("%s.%s", id, extension)
}

func GetGCSSignedURL(ctx context.Context, bucket string, key string, method string, contentType string) (string, error) {
	expire := time.Now().AddDate(0, 0, 7) // expire after 3 days.
	url, err := storage.SignedURL(bucket, key, &storage.SignedURLOptions{
		GoogleAccessID: ServiceAccountName(),
		SignBytes: func(b []byte) ([]byte, error) {
			resp, err := iamService.Projects.ServiceAccounts.SignBlob(
				ServiceAccountID(),
				&iam.SignBlobRequest{BytesToSign: base64.StdEncoding.EncodeToString(b)},
			).Context(ctx).Do()
			if err != nil {
				return nil, err
			}
			return base64.StdEncoding.DecodeString(resp.Signature)
		},
		Method:      method,
		ContentType: contentType,
		Expires:     expire,
	})

	if err != nil {
		return url, err
	}

	return url, nil
}
