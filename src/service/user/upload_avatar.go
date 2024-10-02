package service_user

import (
	"context"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	cloud "cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

type ImageStructure struct {
	ImagePath string `json:"image-path"`
	URLBucket string `json:"url-bucket"`
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("../serviceAcc.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln("sssss", err)
	}
	//khởi tạo firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln("aaaa", err)
	}
	//khởi tạo clound clientchat
	storage, err := cloud.NewClient(ctx, sa)
	if err != nil {
		log.Fatalln(err)
	}

	file, handler, err := r.FormFile("image")
	r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	imagePath := handler.Filename

	bucket := "golang-upload.appspot.com"

	wc := storage.Bucket(bucket).Object(imagePath).NewWriter(ctx)
	_, err = io.Copy(wc, file)
	if err != nil {
		log.Println(err)
		return

	}
	if err := wc.Close(); err != nil {
		log.Println(err)
		return
	}

	res, err := CreateImageUrl(imagePath, bucket, ctx, client)
	if err != nil {
		log.Println("dđ", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))

	log.Println("successfully upload", res)
}

func CreateImageUrl(imagePath string, bucket string, ctx context.Context, client *firestore.Client) (string, error) {
	imageStructure := ImageStructure{
		ImagePath: imagePath,
		URLBucket: "https://storage.cloud.google.com/" + bucket + "/" + imagePath,
	}

	_, _, err := client.Collection("image").Add(ctx, imageStructure)
	if err != nil {
		return "", err
	}

	return imageStructure.URLBucket, nil
}

func SaveUrlImageInDb() {

}
