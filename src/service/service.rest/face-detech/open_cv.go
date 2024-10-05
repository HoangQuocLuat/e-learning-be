package face_detech

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

func DetectFace(deviceID int, xmlFile string) error {
	// open webcam
	webcam, err := gocv.VideoCaptureDevice(deviceID)
	if err != nil {
		return fmt.Errorf("error opening video capture device: %v", err)
	}
	defer webcam.Close()

	// open display window
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		return fmt.Errorf("error reading cascade file: %v", xmlFile)
	}

	fmt.Printf("start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			return fmt.Errorf("cannot read device %d", deviceID)
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		for _, r := range rects {
			// Cắt phần hình ảnh trong vùng hình chữ nhật 'r'
			face := img.Region(r)
			defer face.Close() // Giải phóng vùng nhớ của `Region`

			// Hiển thị khuôn mặt đã cắt
			window2 := gocv.NewWindow("Face")
			defer window2.Close()
			window2.IMShow(face)

			// Tô màu hình chữ nhật lên ảnh gốc
			gocv.Rectangle(&img, r, blue, 3)

			// Tạo text "Human"
			size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
		}

		// show the image in the window, and wait 1 millisecond
		window.IMShow(img)
		key := window.WaitKey(1)
		if key == 27 { // 27 là mã ASCII của phím Esc
			break
		}
	}
	return nil
}
