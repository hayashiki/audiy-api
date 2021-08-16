package gcs

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
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

func (c *client) GenerateV4GetObjectSignedURL(bucketName, obj string) (string, error) {
	// [START storage_generate_signed_url_v4]
	jsonKey, err := ioutil.ReadFile("./credentials/admin-service-account.json")
	if err != nil {
		return "", fmt.Errorf("cannot read the JSON key file, err: %v", err)
	}

	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}

	// 有効期間の最大値は 604,800 秒（7 日間）
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(604800 * time.Second),
	}

	url, err := storage.SignedURL(bucketName, obj, opts)
	if err != nil {
		return "", fmt.Errorf("Unable to generate a signed URL: %v", err)
	}

	// [END storage_generate_signed_url_v4]
	return url, nil
}

// ServiceAccountName returns email address format of google service account.
func ServiceAccountName() string {
	return config.GetProject() + "@appspot.gserviceaccount.com"
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
