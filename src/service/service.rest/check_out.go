package service_rest

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"e-learning/src/cronjob"
	"e-learning/src/database/collection"
	model_attendance "e-learning/src/database/model/attendance"
	face_config "e-learning/src/face-config"
	service_rest_resp "e-learning/src/service/service.rest/response"
	service_user "e-learning/src/service/user"

	"github.com/Kagami/go-face"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gocv.io/x/gocv"
)

//

// TestOpenCam handles the face detection and check-in process
func CheckOut(w http.ResponseWriter, r *http.Request) {
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
		mes := service_rest_resp.Response{
			Status:  0,
			Message: "Không tìm thấy khuôn mặt",
			Data:    nil,
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"message": mes})
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
		mes := service_rest_resp.Response{
			Status:  3,
			Message: "khuôn mặt không tồn tại",
			Data:    nil,
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"message": mes})
		return
	}

	// Check-out logic
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatalf("Lỗi khi tải múi giờ: %v", err)
	}
	today := time.Now().Truncate(24 * time.Hour)
	fil := bson.M{
		"user_id": faceDesc[faceID].ID,
		"time_check_in": bson.M{
			"$gte": today,
			"$lt":  today.Add(24 * time.Hour),
		},
		"status_check_in": "Điểm danh thành công",
	}
	var existingAttendance model_attendance.Attendance
	ctx := context.Background()
	err = collection.Attendance().Collection().FindOne(ctx, fil).Decode(&existingAttendance)
	if err == mongo.ErrNoDocuments {
		mes := service_rest_resp.Response{
			Status:  1,
			Message: "Bạn chưa điểm danh vào",
			Data:    nil,
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"message": mes})
		return
	}
	if err != nil {
		log.Println("Error decoding user:", err)
		return
	}
	update := bson.M{
		"$set": bson.M{
			"time_check_out": time.Now().In(location),
			"updated_at":     time.Now().In(location),
		},
	}
	_, err = collection.Attendance().Collection().UpdateOne(ctx, fil, update)
	if err != nil {
		log.Println("Error updating attendance:", err)
		http.Error(w, "Failed to update attendance", http.StatusInternalServerError)
		return
	}
	go func() {
		subject := "Thông báo điểm danh"
		body := fmt.Sprint("Bạn đã điểm danh vào lúc: ", existingAttendance.TimeCheckOut)
		if err := cronjob.SendMail(faceDesc[faceID].Email, subject, body); err != nil {
			log.Println("Error sending email:", err)
		}
	}()

	mes := service_rest_resp.Response{
		Status:  http.StatusOK,
		Message: "Check-in success",
		Data:    existingAttendance,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mes)
}
