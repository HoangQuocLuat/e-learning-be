package face_recogiton

import (
	"fmt"
	"log"

	"github.com/Kagami/go-face"
)

const dataDir = "train"

func NewRec() face.Recognizer{
	fmt.Println("Facial Recognition System v0.01")

	rec, err := face.NewRecognizer(dataDir)
	if err != nil {
		log.Println("Cannot initialize recognizer:", err)
		return
	}
	defer rec.Close()

	fmt.Println("Recognizer Initialized")

	tonyImage := filepath.Join(dataDir, "tony-stark.jpg")
	tonyFace, err := rec.RecognizeSingleFile(tonyImage)
	if err != nil {
		log.Fatalf("Can't recognize tonyFace: %v", err)
	}
	if tonyFace == nil {
		log.Fatalf("Not a single face on the tonyFace")
	}
}
func CompareFace(inputImage string, dataImages []string) {
	fmt.Println("Facial Recognition System v0.01")

	rec, err := face.NewRecognizer(dataDir)
	if err != nil {
		log.Println("Cannot initialize recognizer:", err)
		return
	}
	defer rec.Close()

	fmt.Println("Recognizer Initialized")

	//nhan dien khuon mat tu anh dau vao
	inputFace, err := rec.RecognizeSingleFile(inputImage)
	if err != nil {
		log.Fatalf("Can't recognize input image: %v", err)
	}
	if inputFace == nil {
		log.Fatalf("Not a single face on the input image")
	}

	inputDescriptor := inputFace.Descriptor

	for i, j := range dataImages {
		dataFace, err := rec.RecognizeSingleFile(j)
		if err != nil {
			log.Printf("Can't recognize image %s: %v", j, err)
			continue
		}
		if dataFace == nil {
			log.Printf("Not a single face on image %s", j)
			continue
		}

		//so sanh khuon mat
		dataFace.Descriptor.Compare()
	}
}
