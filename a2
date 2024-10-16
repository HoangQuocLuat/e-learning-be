package service_rest

import (
	"bytes"
	"context"
	"e-learning/src/database/collection"
	model_attendance "e-learning/src/database/model/attendance"
	face_config "e-learning/src/face-config"
	service_user "e-learning/src/service/user"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"net/http"
	"time"

	"gocv.io/x/gocv"

	"github.com/Kagami/go-face"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

func CheckOutFace(w http.ResponseWriter, r *http.Request) {
	//
	ctx := context.Background()
	// nâng http -> tcp
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()

	//lấy toàn bộ ảnh user trong 1 class
	classID := r.FormValue("class_id")
	faceDesc, err := service_user.GetImagesDescByClassID(ctx, classID)
	if err != nil {
		log.Println("Error getting face descriptions:", err)
		http.Error(w, "Failed to get face descriptions", http.StatusInternalServerError)
		return
	}
	// thiết lập mẫu khuôn mặt nhận diện
	var samples []face.Descriptor
	var labels []int32
	for i, p := range faceDesc {
		samples = append(samples, face.Descriptor(p.DescriptorAvatar))
		labels = append(labels, int32(i))
	}
	face_config.Recognizer.SetSamples(samples, labels)

	//=== dùng gocv  => detech face
	// Open camera
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		log.Printf("Error opening video capture device: %v\n", err)
		http.Error(w, "Error opening video capture device", http.StatusInternalServerError)
		return
	}
	defer webcam.Close()
	//mở cửa sổ hiển thị
	window := gocv.NewWindow("Face Detect")
	defer window.Close()
	//chuẩn bị ma trận ảnh
	img := gocv.NewMat()
	defer img.Close()
	//màu cho khung hình
	blue := color.RGBA{0, 0, 255, 0}
	countFrame := 0
	for {
		// nếu không tìm thấy ảnh
		if ok := webcam.Read(&img); !ok || img.Empty() {
			log.Println("Cannot read device or empty frame")
			continue
		}

		// Delay để đảm bảo chương trình không check liên tục, chỉ check sau vài khung hình
		countFrame++
		if countFrame < 20 {
			gocv.PutText(&img, "Please wait...", image.Pt(10, 30), gocv.FontHersheyPlain, 1.2, blue, 2)
			window.IMShow(img)
			window.WaitKey(1)
			continue
		}

		// Phát hiện khuôn mặt
		classifier := gocv.NewCascadeClassifier()
		if !classifier.Load("/home/ad/Documents/e-learning/e-learning-be/haarcascade_frontalface_default.xml") {
			log.Println("Error loading cascade file")
			http.Error(w, "Failed to load cascade file", http.StatusInternalServerError)
			return
		}
		defer classifier.Close()

		// nhận diện khuôn mặt
		rects := classifier.DetectMultiScale(img)
		if len(rects) == 0 {
			gocv.PutText(&img, "No face detected", image.Pt(10, 30), gocv.FontHersheyPlain, 1.2, blue, 2)
			window.IMShow(img)
			window.WaitKey(1)
			continue
		}

		// Khi có khuôn mặt, cắt phần ảnh khuôn mặt
		faceImg := img.Region(rects[0])
		defer faceImg.Close()

		// Chuyển đổi gocv.Mat thành image.Image
		imgBytes, err := gocv.IMEncode(gocv.JPEGFileExt, faceImg)
		if err != nil {
			log.Println("Error encoding image:", err)
			continue
		}
		imgReader := bytes.NewReader(imgBytes.GetBytes())
		imgDecoded, err := jpeg.Decode(imgReader)
		if err != nil {
			log.Println("Error decoding image:", err)
			continue
		}

		// Chuyển đổi image thành byte slice
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, imgDecoded, nil); err != nil {
			log.Println("Error encoding to byte slice:", err)
			continue
		}
		faceBytes := buf.Bytes()

		// Nhận diện khuôn mặt
		faceDescriptor, err := face_config.Recognizer.RecognizeSingle(faceBytes)
		if err != nil || faceDescriptor == nil {
			// Nếu không nhận diện được khuôn mặt
			message := "Unable to recognize face"
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("Write:", err)
				return
			}
			continue
		}

		// Phân loại khuôn mặt
		faceID := face_config.Recognizer.Classify(faceDescriptor.Descriptor)
		if faceID < 0 {
			// Khuôn mặt không trùng khớp
			message := "Incorrect face detected"
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("Write:", err)
				return
			}
			continue
		}
		// Khuôn mặt đúng, thực hiện các logic khác
		today := time.Now().Truncate(24 * time.Hour)
		fil := bson.M{
			"user_id": faceDesc[faceID].ID,
			"time_check_in": bson.M{
				"$gte": today,
				"$lt":  today.Add(24 * time.Hour),
			},
			"status_check_in": "1",
		}

		var existingAttendance model_attendance.Attendance
		err = collection.Attendance().Collection().FindOne(ctx, fil).Decode(&existingAttendance)
		fmt.Println(existingAttendance)
		if err != nil {
			log.Println("Error decoding user:", err)
			return
		}
		update := bson.M{
			"$set": bson.M{
				"time_check_out": time.Now(),
				"updated_at":     time.Now(),
			},
		}
		_, err = collection.Attendance().Collection().UpdateOne(ctx, fil, update)
		if err != nil {
			log.Println("Error updating attendance:", err)
			http.Error(w, "Failed to update attendance", http.StatusInternalServerError)
			return
		}
		message := fmt.Sprintf("%s Check-out success\n", faceDesc[faceID].Name)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Println("Write:", err)
			return
		}

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read:", err)
				break
			}
			log.Printf("Received: %s", message)
		}
	}
}
