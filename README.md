# cloudinary

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/komfy/cloudinary?style=flat-square)
[![GoDoc](https://godoc.org/github.com/komfy/cloudinary?status.svg)](https://pkg.go.dev/github.com/komfy/cloudinary)

Cloudinary API wrapper in Go.

## Installation

```
go get -u -v github.com/komfy/cloudinary
```

## Examples


### Create a service from environmental variable

```go 
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		log.Fatalln("there is no env variable with give name")
	}
	s, err := cloudinary.NewService(cloudinaryURL)
	if err != nil {
		log.Fatalln(err)
	}
```


### Upload from a form
```go 
type downloadHandler struct {
	cs *cloudinary.Service
}

func (h downloadHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method must be POST", http.StatusMethodNotAllowed)
		return
	}
	//Parsing multipart form is not necessary, FormFile invokes it if form isn't parsed.
	file, fh, err := req.FormFile("file")
	upResp, err := h.cs.Upload(fh.Filename, file, false)
	if err != nil {
		http.Error(res, err.Error(), 505)
		return
	}
	url := upResp.SecureURL
	res.Write([]byte(url))
}
```
### Download from the os.File 
``` go 
	// Sending local file.
	file, err := os.Open("example.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	upResp, err := s.Upload(file.Name(), file, false)
	if err != nil {
		log.Fatalln(err)
	}
```


## License

MIT
