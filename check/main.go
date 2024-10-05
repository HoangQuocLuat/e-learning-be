package main

import (
	"context"
	service_user "e-learning/src/service/user"
	"fmt"
)

func main() {
	// face image to face-detech
	// image into vector
	// so sanh
	ctx := context.Background()
	i, err := service_user.GetImagesByClassID(ctx, "66ed346e58fcd8e1b6d8cea8")
	if err != nil {
		panic(err)
	}
	fmt.Println(i)
}

// func chuyen ve duoi dang vector
func ImageInToVector() {

}
