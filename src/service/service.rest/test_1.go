package service_rest

import (
	"log"
	"net/http"
	"time"

	"gocv.io/x/gocv"
)

// StreamHandler phát luồng video từ camera
func StreamHandler(w http.ResponseWriter, r *http.Request) {
	// Mở camera mặc định
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Printf("Cannot open camera: %v", err)
		http.Error(w, "Cannot open camera", http.StatusInternalServerError)
		return
	}
	defer webcam.Close()

	// Cấu hình header cho response để phát luồng video
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")

	img := gocv.NewMat()
	defer img.Close()

	// Đặt một timeout cho ResponseWriter
	// Điều này giúp tránh việc giữ kết nối mở quá lâu
	w.WriteHeader(http.StatusOK)
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for {
		// Đọc khung hình từ camera
		if ok := webcam.Read(&img); !ok {
			log.Println("Cannot read frame from camera.")
			continue
		}
		if img.Empty() {
			continue
		}

		// Chuyển đổi khung hình thành JPEG
		buf, err := gocv.IMEncode(".jpg", img)
		if err != nil {
			log.Printf("Cannot encode frame: %v", err)
			continue
		}

		// Gửi khung hình dưới dạng `multipart/x-mixed-replace`
		if _, err := w.Write([]byte("--frame\nContent-Type: image/jpeg\n\n")); err != nil {
			log.Println("Client disconnected:", err)
			return // Ngắt vòng lặp nếu client ngắt kết nối
		}
		if _, err := w.Write(buf.GetBytes()); err != nil {
			log.Println("Error writing frame:", err)
			return // Ngắt vòng lặp nếu có lỗi khi ghi
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			log.Println("Error writing end frame:", err)
			return // Ngắt vòng lặp nếu có lỗi khi ghi
		}
		buf.Close() // Đóng buf sau khi gửi

		flusher.Flush() // Đảm bảo dữ liệu được gửi ngay lập tức

		// Thêm một khoảng thời gian để giảm tải CPU
		time.Sleep(30 * time.Millisecond) // Tùy chỉnh khoảng thời gian này nếu cần
	}
}
