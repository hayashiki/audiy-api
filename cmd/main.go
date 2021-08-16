package main

import (
	"context"
	"log"
	"os"

	"github.com/hayashiki/audiy-api/infrastructure/gcs"
)

func main() {
	fileID := "F023EJ55R2A"
	filePath := gcs.StorageObjectFilePath(fileID, "m4a")
	fileType := "" // this is unnecessary when GET request
	bucketName := os.Getenv("GCS_INPUT_AUDIO_BUCKET")
	url, err := gcs.GetGCSSignedURL(context.Background(), bucketName, filePath, "GET", fileType)
	log.Println(url, err)
}
