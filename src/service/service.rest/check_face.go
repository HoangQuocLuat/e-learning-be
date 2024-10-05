package service_rest

import (
	"context"
	face_detech "e-learning/src/service/service.rest/face-detech"
	service_user "e-learning/src/service/user"
	"fmt"
	"net/http"
)

func CheckFace(w http.ResponseWriter, r *http.Request) {
	//get face in database
	ctx := context.Background()
	class_id := r.FormValue("class_id")
	i, err := service_user.GetImagesByClassID(ctx, class_id)
	if err != nil {
		panic(err)
	}
	fmt.Println(i)

	//opencv
	go func() {
		err = face_detech.DetectFace(0, "/home/ad/Documents/e-learning/e-learning-be/src/service/service.rest/face-detech/haarcascade_frontalface_default.xml")
		if err != nil {
			panic(err)
		}
	}()

	//
	
}
