package service_user

import (
	"bytes"
	"context"
	"e-learning/src/database/collection"
	model_user "e-learning/src/database/model/user"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	cloud "cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/Kagami/go-face"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/api/option"
)

const dataDir = "/home/ad/Documents/e-learning/e-learning-be/train"

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
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln("aaaa", err)
	}
	storage, err := cloud.NewClient(ctx, sa)
	if err != nil {
		log.Fatalln(err)
	}

	//lấy file ảnh từ người dùng
	file, handler, err := r.FormFile("image")
	r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	// Đọc dữ liệu ảnh vào một biến bytes.Buffer
	var imgData bytes.Buffer
	_, err = io.Copy(&imgData, file)
	if err != nil {
		log.Println("Error copying image data:", err)
		return
	}

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
	userId := r.FormValue("user_id")

	rec, err := face.NewRecognizer(dataDir)
	if err != nil {
		log.Println("Cannot initialize recognizer:", err)
		return
	}

	userFace, err := rec.RecognizeSingle(imgData.Bytes())
	if err != nil {
		log.Fatalf("Can't recognize Face: %v", err)
	}
	if userFace == nil {
		log.Fatalf("Not a single face on the Face")
	}

	descriptorInterface := make([]interface{}, len(userFace.Descriptor))
	for i, v := range userFace.Descriptor {
		descriptorInterface[i] = v
	}
	update := bson.M{
		"$set": bson.M{
			"descriptor_avatar": descriptorInterface,
		},
	}
	_, err = collection.User().Collection().UpdateOne(ctx, bson.M{"_id": userId}, update)

	if err != nil {
		log.Println("Error updating user")
		return
	}

	// Lưu URL hình ảnh vào cơ sở dữ liệu
	err = SaveUrlImageIntoDb(ctx, res, userId)
	if err != nil {
		log.Println("Failed to save image URL into DB", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))

	log.Println("Successfully uploaded and saved image", res)
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

func SaveUrlImageIntoDb(ctx context.Context, urlImage string, userId string) error {
	var user *model_user.User
	err := collection.User().Collection().FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		log.Println("Error decoding find user_id", err)
	}

	update := bson.M{
		"$set": bson.M{
			"avatar": urlImage,
		},
	}

	_, err = collection.User().Collection().UpdateOne(ctx, bson.M{"_id": userId}, update)

	if err != nil {
		log.Println("Error updating user")
		return err
	}

	return nil
}
