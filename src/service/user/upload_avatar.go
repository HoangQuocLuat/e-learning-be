package service_user

import (
	"bytes"
	"context"
	"e-learning/src/database/collection"
	model_user "e-learning/src/database/model/user"
	face_config "e-learning/src/face-config"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	cloud "cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/api/option"
)

type ImageStructure struct {
	ImagePath string `json:"image-path"`
	URLBucket string `json:"url-bucket"`
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("../serviceAcc.json")
	bucket := "golang-upload.appspot.com"

	// Lấy file ảnh từ người dùng
	file, handler, err := r.FormFile("image")
	if err != nil {
		log.Println("Error retrieving the file:", err)
		http.Error(w, "Unable to upload image", http.StatusBadRequest)
		return
	}
	defer file.Close()
	imagePath := handler.Filename

	// Tạo một buffer để lưu trữ dữ liệu ảnh và sử dụng TeeReader để copy vào Firebase đồng thời lưu vào buffer
	var imgData bytes.Buffer
	tee := io.TeeReader(file, &imgData)

	// Khởi tạo Firebase app
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln("Error initializing Firebase app:", err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln("Error initializing Firestore client:", err)
	}
	storage, err := cloud.NewClient(ctx, sa)
	if err != nil {
		log.Fatalln("Error initializing Storage client:", err)
	}

	// Ghi ảnh vào Firebase Storage
	wc := storage.Bucket(bucket).Object(imagePath).NewWriter(ctx)
	_, err = io.Copy(wc, tee)
	if err != nil {
		log.Println("Error copying to Firebase:", err)
		http.Error(w, "Error uploading image", http.StatusInternalServerError)
		return
	}
	if err := wc.Close(); err != nil {
		log.Println("Error closing Firebase writer:", err)
		http.Error(w, "Error finalizing image upload", http.StatusInternalServerError)
		return
	}

	// Tạo URL cho image avatar từ Firebase
	imageURL, err := CreateImageUrl(imagePath, bucket, ctx, client)
	if err != nil {
		log.Println("Error creating image URL:", err)
		http.Error(w, "Unable to generate image URL", http.StatusInternalServerError)
		return
	}

	// Nhận diện khuôn mặt từ dữ liệu ảnh
	userFace, err := face_config.Recognizer.RecognizeSingle(imgData.Bytes())
	if err != nil {
		log.Fatalf("Can't recognize face: %v", err)
	}
	if userFace == nil {
		log.Fatalf("Not a single face detected")
	}

	// Chuyển đổi descriptor khuôn mặt thành slice interface{} để lưu vào DB
	descriptorInterface := make([]interface{}, len(userFace.Descriptor))
	for i, v := range userFace.Descriptor {
		descriptorInterface[i] = v
	}

	// Lấy user_id từ form
	userId := r.FormValue("user_id")
	if userId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Lưu URL hình ảnh và descriptor vào cơ sở dữ liệu
	err = SaveUrlImageIntoDb(ctx, imageURL, descriptorInterface, userId)
	if err != nil {
		log.Println("Error saving image URL to DB:", err)
		http.Error(w, "Unable to save image information", http.StatusInternalServerError)
		return
	}

	// Trả về URL ảnh cho người dùng
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(imageURL))

	log.Println("Successfully uploaded and saved image:", imageURL)
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

func SaveUrlImageIntoDb(ctx context.Context, urlImage string, imgDesc []interface{}, userId string) error {
	var user *model_user.User
	err := collection.User().Collection().FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		log.Println("Error decoding user:", err)
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"avatar":            urlImage,
			"descriptor_avatar": imgDesc,
		},
	}

	_, err = collection.User().Collection().UpdateOne(ctx, bson.M{"_id": userId}, update)
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}

	return nil
}
