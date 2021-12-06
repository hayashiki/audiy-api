package gcs

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/hayashiki/audiy-api/src/config"

	"google.golang.org/api/option"

	"github.com/pkg/errors"

	"google.golang.org/api/iam/v1"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

//GetFileContentType retrieves the content type of files
func getFileContentType(out io.Reader) (string, error) {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
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
