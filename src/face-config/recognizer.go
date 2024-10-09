package face_config

import (
	"log"

	"github.com/Kagami/go-face"
)

var Recognizer *face.Recognizer

// InitRecognizer khởi tạo Recognizer
func InitRecognizer(dataDir string) {
	var err error
	Recognizer, err = face.NewRecognizer(dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize Recognizer: %v", err)
	}
	log.Println("Recognizer initialized")
}

// CloseRecognizer đóng Recognizer
func CloseRecognizer() {
	if Recognizer != nil {
		Recognizer.Close()
		log.Println("Recognizer closed")
	}
}
