package service_rest

import (
	"context"
	"e-learning/src/database/collection"
	model_attendance "e-learning/src/database/model/attendance"
	model_tuition "e-learning/src/database/model/tuition"
	face_config "e-learning/src/face-config"
	service_user "e-learning/src/service/user"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Kagami/go-face"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func CheckFace(w http.ResponseWriter, r *http.Request) {
	// Connect WebSocket
	ctx := context.Background()
	classID := r.FormValue("class_id")

	faceDesc, err := service_user.GetImagesDescByClassID(ctx, classID)
	if err != nil {
		log.Println("Error getting face descriptions:", err)
		http.Error(w, "Failed to get face descriptions", http.StatusInternalServerError)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()

	// Thiết lập mẫu khuôn mặt cho nhận diện
	var samples []face.Descriptor
	var labels []int32
	for i, p := range faceDesc {
		samples = append(samples, face.Descriptor(p.DescriptorAvatar))
		labels = append(labels, int32(i))
	}
	face_config.Recognizer.SetSamples(samples, labels)

	// Nhận diện khuôn mặt từ tệp
	tonyImage := filepath.Join("/home/ad/Documents/e-learning/e-learning-be/src/service/service.rest/imggg/hl.jpg")
	tonyFace, err := face_config.Recognizer.RecognizeSingleFile(tonyImage)
	if err != nil {
		log.Fatalf("Can't recognize tonyFace: %v", err)
	}
	if tonyFace == nil {
		log.Fatalf("Not a single face detected in tonyFace")
	}

	// Phân loại khuôn mặt kiểm tra
	faceID := face_config.Recognizer.Classify(tonyFace.Descriptor)
	if faceID < 0 {
		log.Fatalf("Can't classify the face")
	}

	// Kiểm tra và cập nhật thời gian điểm danh
	today := time.Now().Truncate(24 * time.Hour)
	filter := bson.M{
		"user_id": faceDesc[faceID].ID,
		"time_check_in": bson.M{
			"$gte": today,
			"$lt":  today.Add(24 * time.Hour),
		},
	}

	res := &model_attendance.Attendance{}
	condition := collection.Attendance().Collection().FindOne(ctx, filter).Decode(&res)
	if condition != mongo.ErrNoDocuments {
		message := fmt.Sprintf("%s has already checked in\n", faceDesc[faceID].Name)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Println("Write:", err)
		}
		return
	}

	// Điểm danh mới
	res = &model_attendance.Attendance{
		ID:          primitive.NewObjectID().Hex(),
		UserID:      faceDesc[faceID].ID,
		TimeCheckIn: time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err = collection.Attendance().Collection().InsertOne(ctx, res)
	if err != nil {
		log.Println("Error inserting attendance:", err)
		http.Error(w, "Failed to insert attendance", http.StatusInternalServerError)
		return
	}
	message := fmt.Sprintf("%s check-in success.\n", faceDesc[faceID].Name)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Println("Write:", err)
		return
	}

	// Cập nhật thông tin học phí
	var tuition model_tuition.Tuition
	err = collection.Tuition().Collection().FindOne(ctx, bson.M{"user_id": faceDesc[faceID].ID}).Decode(&tuition)
	if err == mongo.ErrNoDocuments {
		tuition = model_tuition.Tuition{
			ID:           primitive.NewObjectID().Hex(),
			UserID:       faceDesc[faceID].ID,
			LessonsCount: 1,
			Price:        30,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		// Chèn học phí mới vào DB
		_, err = collection.Tuition().Collection().InsertOne(ctx, tuition)
		if err != nil {
			log.Println("Error inserting tuition:", err)
			http.Error(w, "Failed to insert tuition", http.StatusInternalServerError)
			return
		}
	} else {
		newTuition := tuition.LessonsCount + 1
		newPrice := tuition.Price + 30
		updateTuition := bson.M{
			"$set": bson.M{
				"lessons_count": newTuition,
				"price":         newPrice,
			},
		}
		_, err = collection.Tuition().Collection().UpdateOne(ctx, bson.M{"user_id": faceDesc[faceID].ID}, updateTuition)
		if err != nil {
			log.Println("Error updating tuition:", err)
			http.Error(w, "Failed to update tuition", http.StatusInternalServerError)
			return
		}
	}

	// Check-out
	time.Sleep(2 * time.Second)
	tonyFace1, err := face_config.Recognizer.RecognizeSingleFile(tonyImage)
	if err != nil {
		log.Fatalf("Can't recognize tonyFace: %v", err)
	}
	if tonyFace1 == nil {
		log.Fatalf("Not a single face detected in tonyFace")
	}

	faceID1 := face_config.Recognizer.Classify(tonyFace1.Descriptor)
	if faceID1 < 0 {
		log.Fatalf("Can't classify the face")
	}

	filter2 := bson.M{
		"user_id":       faceDesc[faceID].ID,
		"time_checkout": bson.M{"$exists": false},
	}

	var existingAttendance model_attendance.Attendance
	err = collection.Attendance().Collection().FindOne(ctx, filter2).Decode(&existingAttendance)
	if err != mongo.ErrNoDocuments {
		update := bson.M{
			"$set": bson.M{
				"time_check_out": time.Now(),
				"updated_at":     time.Now(),
			},
		}
		_, err = collection.Attendance().Collection().UpdateOne(ctx, filter2, update)
		if err != nil {
			log.Println("Error updating attendance:", err)
			http.Error(w, "Failed to update attendance", http.StatusInternalServerError)
			return
		}
		message = fmt.Sprintf("%s Check-out success\n", faceDesc[faceID].Name)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Println("Write:", err)
			return
		}
	}

	// Đọc tin nhắn từ WebSocket
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read:", err)
			break
		}
		log.Printf("Received: %s", message)
		// Xử lý thêm các tin nhắn nếu cần
	}
}
