# cloudinary

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/komfy/cloudinary?style=flat-square)
![GoDoc](https://godoc.org/github.com/komfy/cloudinary?status.svg)](https://godoc.org/github.com/komfy/cloudinary)

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
	"io"
	"log"
	"net/http"
	"github.com/komfy/api/cloudinary"
)

func main() {
    //Initialising the service struct
    service, err := cloudinary.Dial("cloudinary://api_key:api_secret@cloud_name")
	if err != nil {
        log.Fatalln(err)
    }
	http.HandleFunc("/", handle)
	http.ListenAndServe(":3000", nil)
}
//Just an example handle function 
func handle(res http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(32 << 20)
	_, fh, err := req.FormFile("file")
	if err != nil {
		log.Fatalln(err)
	}
	upResp, err := service.UploadFile(fh, false)
	if err != nil {
		io.WriteString(res, err.Error())
    }
    //Do something with response
	fmt.Println(upResp)
}

```

## License

MIT
