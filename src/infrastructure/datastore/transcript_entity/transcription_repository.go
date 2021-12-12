package transcript_entity
//
//import (
//	"context"
//	"log"
//
//	"cloud.google.com/go/datastore"
//	"github.com/hayashiki/audiy-api/src/domain/model"
//)
//
//type transcriptRepository struct {
//	client *datastore.Client
//}
//
//// Save saves transcription
//func (repo *transcriptRepository) Save(ctx context.Context, transcript *model.Transcript) error {
//	log.Println("db save", transcript)
//	// TODO: if exists
//	key, err := repo.client.Put(ctx, datastore.IncompleteKey(model.TranscriptKind, nil), transcript)
//	if err != nil {
//		log.Println("db err", err)
//		return err
//	}
//
//	transcript.Key = key
//	transcript.ID = key.ID
//	return err
//}
//
//// Get user
//func (repo *transcriptRepository) Get(ctx context.Context, id int64) (*model.Transcript, error) {
//	var tran model.Transcript
//	err := repo.client.Get(ctx, datastore.IDKey(model.TranscriptKind, id, nil), &tran)
//	if err != nil {
//		return nil, err
//	}
//	tran.ID = tran.Key.ID
//	return &tran, err
//}
