//Package cloudinary provides an easy way of connection between go and cloudinary
package cloudinary

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"
)

//Service represents cloudinary service
type Service struct {
	cloudName string
	apiKey    string
	apiSecret string
	uploadURL *url.URL
	adminURL  *url.URL
	verbose   bool
	simulate  bool
	resType   int
}

//UploadResponse ...
type UploadResponse struct {
	PublicID     string    `json:"public_id"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	Format       string    `json:"format"`
	ResourseType string    `json:"resource_type"`
	CreatedAt    time.Time `json:"created_at"`
	SecureURL    string    `json:"secure_url"`
	URL          string    `json:"url"`
}

// {
// 	"public_id":"Screen Shot 2020-01-12 at 1.20.21 PM",
// 	"version":1580912870,
// 	"signature":"3137036a0f6a742dbde09074d44e7d8528409682",
// 	"width":1440,
// 	"height":900,
// 	"format":"png",
// 	"resource_type":"image",
// 	"created_at":"2020-02-05T14:27:50Z",
// 	"tags":[

// 	],
// 	"bytes":1135076,
// 	"type":"upload",
// 	"etag":"14f65dc9ab21f5868a0f2b29c02f326a",
// 	"placeholder":false,
// 	"url":"http://res.cloudinary.com/mlvni/image/upload/v1580912870/Screen%20Shot%202020-01-12%20at%201.20.21%20PM.png",
// 	"secure_url":"https://res.cloudinary.com/mlvni/image/upload/v1580912870/Screen%20Shot%202020-01-12%20at%201.20.21%20PM.png",
// 	"original_filename":"Screen Shot 2020-01-12 at 1.20.21 PM"
//  }

const (
	baseUploadURL string = "https://api.cloudinary.com/v1_1"
	imageType     int    = 0
)

//Dial configurates cloudinary service
// cloudinary://api_key:api_secret@cloud_name
func Dial(uri string) (*Service, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "cloudinary" {
		return nil, errors.New("URL scheme is not cloudinary")
	}
	secret, exists := u.User.Password()
	if !exists {
		return nil, errors.New("There is no api secret provided in URL")
	}
	s := &Service{
		cloudName: u.Hostname(),
		apiKey:    u.User.Username(),
		apiSecret: secret,
		resType:   imageType,
		simulate:  false,
		verbose:   false,
	}
	up, err := url.Parse(fmt.Sprintf("%s/%s/image/upload/", baseUploadURL, s.cloudName))
	if err != nil {
		return nil, err
	}
	s.uploadURL = up
	admURL, err := url.Parse(fmt.Sprintf("%s/%s", baseUploadURL, s.cloudName))
	if err != nil {
		return nil, err
	}
	admURL.User = url.UserPassword(s.apiKey, s.apiSecret)
	s.adminURL = admURL

	return s, nil
}

//UploadFile receives file, most probably from Multipart Form, uploads to cloud and returns
//a link to the resource
func (s *Service) UploadFile(fh *multipart.FileHeader, randomPublicID bool) (*UploadResponse, error) {
	var publicID string
	file, err := fh.Open()
	if err != nil {
		return nil, err
	}
	fileBuf, err := ioutil.ReadAll(file)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	if len(fileBuf) == 0 {
		return nil, fmt.Errorf("Not allowed to upload empty files: %s", fh.Filename)
	}
	filename := trimExt(fh.Filename)

	//Creating a form body for request
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	//Writing a public_id field for request
	if !randomPublicID {
		publicID = filename
		pi, err := mw.CreateFormField("public_id")
		if err != nil {
			return nil, err
		}
		pi.Write([]byte(publicID))
	}
	//Writing an API key
	ak, err := mw.CreateFormField("api_key")
	if err != nil {
		return nil, err
	}
	ak.Write([]byte(s.apiKey))

	//Writing timestamp
	ts, err := mw.CreateFormField("timestamp")
	if err != nil {
		return nil, err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	ts.Write([]byte(timestamp))

	//Writing signature
	si, err := mw.CreateFormField("signature")
	if err != nil {
		return nil, err
	}
	hash := sha1.New()
	part := fmt.Sprintf("timestamp=%s%s", timestamp, s.apiSecret)
	if !randomPublicID {
		part = fmt.Sprintf("public_id=%s&%s", publicID, part)
	}
	io.WriteString(hash, part)
	signature := fmt.Sprintf("%x", hash.Sum(nil))
	si.Write([]byte(signature))

	fi, err := mw.CreateFormFile("file", fh.Filename)
	if err != nil {
		return nil, err
	}
	fi.Write(fileBuf)

	err = mw.Close()
	if err != nil {
		return nil, err
	}

	uploadURL := s.uploadURL.String()

	req, err := http.NewRequest("POST", uploadURL, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	upResp := new(UploadResponse)
	err = json.Unmarshal(body, upResp)
	if err != nil {
		return nil, err
	}
	return upResp, err
}

func trimExt(filename string) string {
	fileExt := filepath.Ext(filename)
	return filename[0 : len(filename)-len(fileExt)]
}
