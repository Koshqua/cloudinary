package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/komfy/cloudinary"
)

type downloadHandler struct {
	cs *cloudinary.Service
}

func (h downloadHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method must be POST", http.StatusMethodNotAllowed)
		return
	}
	// Parsing multipart form is not necessary, FormFile invokes it if form isn't parsed.
	file, fh, err := req.FormFile("file")
	upResp, err := h.cs.Upload(fh.Filename, file, false)
	if err != nil {
		http.Error(res, err.Error(), 505)
		return
	}
	url := upResp.SecureURL
	res.Write([]byte(url))
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
	// Sending local file.
	file, err := os.Open("example.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	upResp, err := s.Upload(file.Name(), file, false)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(upResp.URL)
	http.Handle("/download", downloadHandler{cs: s})
	http.ListenAndServe(":3000", nil)
}
