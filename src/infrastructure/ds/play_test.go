package ds

//func TestAudioUser(t *testing.T) {
//	log.Println(os.Getenv("GCP_PROJECT"))
//	ctx := context.Background()
//	dsDataSource := Connect()
//	audioRepo := NewAudioRepository(dsDataSource)
//	audio, _ := audioRepo.Find(ctx, "")
//
//	audioUser := entity.NewPlay(111111, audio.ID)
//
//	audioUserRepo := audioUserRepository{dsDataSource}
//	t.Log(audioUser)
//	err := audioUserRepo.Save(ctx, audioUser)
//	if err != nil {
//		t.Fatal(err)
//	}
//}
//
//func TestAnotherAudioUser(t *testing.T) {
//	log.Println(os.Getenv("GCP_PROJECT"))
//	ctx := context.Background()
//	dsDataSource := Connect()
//	audioRepo := NewAudioRepository(dsDataSource)
//	audio, _ := audioRepo.Find(ctx, "")
//
//	audioUser := entity.NewPlay(111111, audio.ID)
//
//	playRepo := playRepository{dsDataSource}
//	t.Log(audioUser)
//	err := playRepo.Save(ctx, audioUser)
//	if err != nil {
//		t.Fatal(err)
//	}
//}
//
//func TestExistsAudioUser(t *testing.T) {
//	log.Println(os.Getenv("GCP_PROJECT"))
//	ctx := context.Background()
//	dsDataSource := Connect()
//	audioRepo := NewAudioRepository(dsDataSource)
//	audio, _ := audioRepo.Find(ctx, "")
//	playRepo := playRepository{dsDataSource}
//	_, err := audioUserRepo.Exists(ctx, 111111, audio.ID)
//	if err != nil {
//		t.Fatal(err)
//	}
//}

