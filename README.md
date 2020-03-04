# cloudinary

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/komfy/cloudinary?style=flat-square)
[![GoDoc](https://godoc.org/github.com/komfy/cloudinary?status.svg)](pkg.go.dev/github.com/komfy/cloudinary)

Cloudinary API wrapper in Go.

## Installation

```
go get -u -v github.com/komfy/cloudinary
```

## Examples

### Create a service from enviromental variable

``` go 
s, err := cloudinary.NewServiceFromEnv("CLOUDINARY_URL")
if err != nil {
	log.Fatalln(err)
}
```
### Create a service from cloudinary url

```go 
s, err := cloudinary.NewService("cloudinary://api_key:api_secret@cloud_name")

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
```
### Download from the os.File 
``` go 
file, err := os.Open("example.jpg")

if err != nil {
	log.Fatalln(err)
}

upResp, err := s.UploadLocalFile(file, false)

if err != nil {
	log.Fatalln(err)
}
```
### Raw upload from file data []byte and filename
``` go 
file, err := os.Open("example.jpg")

if err != nil {
	log.Fatalln(err)
}

buffer, err := ioutil.ReadAll(file)

if err != nil {
	log.Fatalln(err)
}
// There is no need to take care of file extensions. Upload will take trim it by it's own.
upResp, err = s.Upload(file.Name(), buffer, false)

if err != nil {
	log.Fatalln(err)
}
```

## License

MIT
