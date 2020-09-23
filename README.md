# cloudinary

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/komfy/cloudinary?style=flat-square)
[![GoDoc](https://godoc.org/github.com/komfy/cloudinary?status.svg)](https://pkg.go.dev/github.com/komfy/cloudinary)

Cloudinary API wrapper in Go.

## Installation

```
go get -u -v github.com/komfy/cloudinary
```

## Example

```go 
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
//Creating a handler to donwload from form.
//I prefer to use handlers, because it's easier to add some external services into it's logic.
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
```

## Uploading mutiple form files

```go

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


```


## License

MIT
