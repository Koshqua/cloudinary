package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/komfy/cloudinary"
)

type downloadHandler struct {
	cs *cloudinary.Service
}

//UploadFiles func
func (h downloadHandler) UploadMultipleFiles(files []*multipart.FileHeader) []string {
	urls := []string{}
	for _, file := range files {
		fmt.Println(file.Filename)
		data, err := file.Open()
		if err == nil {
			upResp, err := h.cs.Upload(file.Filename, data, false)
			if err != nil {
				panic(err)
			} else {
				urls = append(urls, upResp.URL)
			}
		}
		data.Close()
	}
	return urls
}

func (h downloadHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method must be POST", http.StatusMethodNotAllowed)
		return
	}

	req.ParseMultipartForm(32 << 20)
	files := req.MultipartForm.File["files"] //name of file input
	urls := h.UploadMultipleFiles(files)
	resBody, _ := json.Marshal(urls)
	res.Write(resBody)
}

func main() {
	// Creates new Service from environmental variable
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		log.Fatalln("there is no env variable with given name")
	}
	s, err := cloudinary.NewService(cloudinaryURL)
	if err != nil {
		log.Fatalln(err)
	}

	http.Handle("/download", downloadHandler{cs: s})
	http.ListenAndServe(":3000", nil)
}
