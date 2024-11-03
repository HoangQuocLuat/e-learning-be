package service_rest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"e-learning/src/database/collection"
	model_attendance "e-learning/src/database/model/attendance"
	face_config "e-learning/src/face-config"
	service_user "e-learning/src/service/user"

	"github.com/Kagami/go-face"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gocv.io/x/gocv"
)

// TestOpenCam handles the face detection and check-in process
func TestOpenCam(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Image   string `json:"image"`
		ClassID string `json:"class_id"`
	}

	// Decode JSON from request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Decode the base64 image
	imageData := req.Image
	imageData = strings.TrimPrefix(imageData, "data:image/jpeg;base64,")
	imgBytes, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		http.Error(w, "Error decoding image", http.StatusInternalServerError)
		return
	}

	// Convert bytes to gocv Mat
	img, err := gocv.IMDecode(imgBytes, gocv.IMReadColor)
	if err != nil {
		log.Println("Error decoding image:", err)
		http.Error(w, "Error decoding image", http.StatusInternalServerError)
		return
	}
	defer img.Close()

	// Get all user images for the specified class
	faceDesc, err := service_user.GetImagesDescByClassID(r.Context(), req.ClassID)
	if err != nil {
		log.Println("Error getting face descriptions:", err)
		http.Error(w, "Failed to get face descriptions", http.StatusInternalServerError)
		return
	}

	// Set up face recognition samples
	var samples []face.Descriptor
	var labels []int32
	for i, p := range faceDesc {
		samples = append(samples, face.Descriptor(p.DescriptorAvatar))
		labels = append(labels, int32(i))
	}
	face_config.Recognizer.SetSamples(samples, labels)

	// Load classifier
	classifier := gocv.NewCascadeClassifier()
	if !classifier.Load("/home/ad/Documents/e-learning/e-learning-be/haarcascade_frontalface_default.xml") {
		log.Println("Error loading cascade file")
		http.Error(w, "Failed to load cascade file", http.StatusInternalServerError)
		return
	}
	defer classifier.Close()

	// Detect faces
	rects := classifier.DetectMultiScale(img)
	if len(rects) == 0 {
		message := "No face detected"
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": message})
		return
	}

	// Process detected face
	faceImg := img.Region(rects[0])
	defer faceImg.Close()

	// Convert Mat to byte slice
	imgBytesEncoded, err := gocv.IMEncode(gocv.JPEGFileExt, faceImg)
	if err != nil {
		log.Println("Error encoding image:", err)
		http.Error(w, "Error encoding image", http.StatusInternalServerError)
		return
	}

	// Recognize face
	faceDescriptor, err := face_config.Recognizer.RecognizeSingle(imgBytesEncoded.GetBytes())
	if err != nil {
		log.Println("Face recognition error:", err)
		http.Error(w, "Face recognition error", http.StatusInternalServerError)
		return
	}

	// Classify face
	faceID := face_config.Recognizer.Classify(faceDescriptor.Descriptor)
	if faceID < 0 {
		http.Error(w, "Incorrect face detected", http.StatusUnauthorized)
		return
	}

	// Check-in logic
	today := time.Now().Truncate(24 * time.Hour)
	filter := bson.M{
		"user_id": faceDesc[faceID].ID,
		"time_check_in": bson.M{
			"$gte": today,
			"$lt":  today.Add(24 * time.Hour),
		},
	}

	// Check if already checked in
	res := &model_attendance.Attendance{}
	condition := collection.Attendance().Collection().FindOne(r.Context(), filter).Decode(&res)
	if condition != mongo.ErrNoDocuments {
		message := fmt.Sprintf("%s has already checked in\n", faceDesc[faceID].Name)
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": message})
		return
	}

	// Save new check-in
	newAttendance := &model_attendance.Attendance{
		ID:            primitive.NewObjectID().Hex(),
		UserID:        faceDesc[faceID].ID,
		TimeCheckIn:   time.Now(),
		StatusCheckIn: "1",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	_, err = collection.Attendance().Collection().InsertOne(r.Context(), newAttendance)
	if err != nil {
		log.Println("Error inserting attendance:", err)
		http.Error(w, "Failed to insert attendance", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	message := fmt.Sprintf("%s check-in success.\n", faceDesc[faceID].Name)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
