package face_detech

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"gocv.io/x/gocv"
)

func DetectFace(deviceID int, xmlFile string) ([]byte, error) {
	// mở camera
	webcam, err := gocv.VideoCaptureDevice(deviceID)
	if err != nil {
		return nil, fmt.Errorf("error opening video capture device: %v", err)
	}
	defer webcam.Close()

	// mở cửa sổ hiển thị
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	// chuẩn bị ma trận ảnh
	img := gocv.NewMat()
	defer img.Close()

	// màu cho hình chữ nhật khi phát hiện khuôn mặt
	blue := color.RGBA{0, 0, 255, 0}

	// load classifier nhận diện khuôn mặt
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		return nil, fmt.Errorf("error reading cascade file: %v", xmlFile)
	}

	countFrame := 0
	fmt.Printf("Start reading camera device: %v\n", deviceID)

	for {
		if ok := webcam.Read(&img); !ok {
			return nil, fmt.Errorf("cannot read device %d", deviceID)
		}
		if img.Empty() {
			continue
		}

		// Tăng biến đếm khung hình
		countFrame++
		if countFrame < 20 { // Đợi 20 khung hình trước khi bắt đầu nhận diện
			gocv.PutText(&img, "Please wait...", image.Pt(10, 30), gocv.FontHersheyPlain, 1.2, blue, 2)
			window.IMShow(img)
			window.WaitKey(1)
			continue
		}

		// phát hiện khuôn mặt
		rects := classifier.DetectMultiScale(img)

		if len(rects) > 0 { // Khi có khuôn mặt
			for _, r := range rects {
				// Cắt phần khuôn mặt
				face := img.Region(r)
				defer face.Close()

				// Hiển thị khuôn mặt đã cắt
				window2 := gocv.NewWindow("Face")
				defer window2.Close()
				window2.IMShow(face)

				// Tô màu chữ nhật lên ảnh gốc
				gocv.Rectangle(&img, r, blue, 3)

				// Viết chữ "Human"
				size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
				pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
				gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)

				// Chuyển hình ảnh thành []byte và trả về
				imageBytes := img.ToBytes()
				return imageBytes, nil
			}
		} else { // Không phát hiện khuôn mặt, tiếp tục
			gocv.PutText(&img, "No face detected, waiting...", image.Pt(10, 30), gocv.FontHersheyPlain, 1.2, blue, 2)
			window.IMShow(img)
			window.WaitKey(1)
			time.Sleep(1 * time.Second) // Chờ 1 giây trước khi kiểm tra tiếp
		}

		// Hiển thị hình ảnh và đợi 1ms
		window.IMShow(img)
		if window.WaitKey(1) == 27 { // Dừng khi nhấn phím Esc
			break
		}
	}

	return nil, fmt.Errorf("stopped manually")
}
