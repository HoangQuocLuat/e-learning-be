package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"

	"gocv.io/x/gocv"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("How to run:\n\tfacedetect [camera ID] [classifier XML file]")
		return
	}

	// parse args
	deviceID, _ := strconv.Atoi(os.Args[1])
	xmlFile := os.Args[2]

	// open webcam
	webcam, err := gocv.VideoCaptureDevice(int(deviceID))
	if err != nil {
		fmt.Println(err)
		return
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
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	fmt.Printf("start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		// draw a rectangle around each face on the original image,
		// along with text identifying as "Human"
		// for _, r := range rects {
		// 	gocv.Rectangle(&img, r, blue, 3)

		// 	size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
		// 	pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
		// 	gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
		// }

		//chat gen
		for _, r := range rects {
			// Cắt phần hình ảnh trong vùng hình chữ nhật 'r'
			face := img.Region(r)
			defer face.Close() // Giải phóng vùng nhớ của `Region`

			// Hiển thị khuôn mặt đã cắt
			window2 := gocv.NewWindow("Face")
			defer window2.Close()
			window2.IMShow(face)
			// gocv.WaitKey(0)

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
}

// get images in database
// func getImages() []string {

// 	return images
// }