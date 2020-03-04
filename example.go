package example

import (
	"fmt"
	"io/ioutil"
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
	//UploadFromForm handles all the work
	//You dont need to parse form or open any files.
	//In case if RandomPublicId is false - cloudinary will generate pub. ID from file name.
	upResp, err := h.cs.UploadFromForm(req, "file", false)
	if err != nil {
		http.Error(res, err.Error(), 505)
		return
	}
	url := upResp.SecureURL
	res.Write([]byte(url))
}

func main() {
	//Creates new Service from enviromental variable
	//In case you want to Create Service from a string
	//Use cloudinary.NewService(url string)
	s, err := cloudinary.NewServiceFromEnv("CLOUDINARY_URL")
	if err != nil {
		log.Fatalln(err)
	}
	file, err := os.Open("example.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	upResp, err := s.UploadLocalFile(file, false)
	if err != nil {
		log.Fatalln(err)
	}
	//The same but with the raw cloudinary.UploadFunction
	//So here it's up to you how you wanna to handle files.
	//You just need to provide the filename and []byte as a file data.
	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	//There is no need to take care of file extensions. Upload will take trim it by it's own.
	upResp, err = s.Upload(file.Name(), buffer, false)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(upResp.URL)
	http.Handle("/download", downloadHandler{cs: s})
	http.ListenAndServe(":3000", nil)
}
