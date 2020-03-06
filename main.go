// Package cloudinary provides an easy way of connection between go and cloudinary
package cloudinary

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
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

// Service represents cloudinary service
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

// UploadResponse ...
type UploadResponse struct {
	PublicID     string    `json:"public_id"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	Format       string    `json:"format"`
	ResourceType string    `json:"resource_type"`
	CreatedAt    time.Time `json:"created_at"`
	SecureURL    string    `json:"secure_url"`
	URL          string    `json:"url"`
}

const (
	baseUploadURL string = "https://api.cloudinary.com/v1_1"
	imageType     int    = 0
)

// NewService inits a *Service from cloudinaryURL
// Link should be given in format
// "cloudinary://api_key:api_secret@cloud_name"
// After initialization returns ready to use service or an error in case of incorrect URL
func NewService(cloudinaryURL string) (*Service, error) {
	u, err := url.Parse(cloudinaryURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "cloudinary" {
		return nil, errNotCloudinary
	}
	secret, exists := u.User.Password()
	if !exists {
		return nil, errNoAPISecret
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

// Upload uploads a file
func (s *Service) Upload(fileName string, fileBody io.Reader, randomPublicID bool) (*UploadResponse, error) {
	var publicID string
	file, err := ioutil.ReadAll(fileBody)
	if err != nil {
		return nil, err
	}
	if len(file) == 0 {
		return nil, fmt.Errorf("not allowed to upload empty files: %s", fileName)
	}
	fNameWithoutExt := trimExt(fileName)

	// Creating a form body for request
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	//Writing a public_id field for request
	if !randomPublicID {
		publicID = fNameWithoutExt
		pi, err := mw.CreateFormField("public_id")
		if err != nil {
			return nil, err
		}
		pi.Write([]byte(publicID))
	}
	// Writing an API key
	ak, err := mw.CreateFormField("api_key")
	if err != nil {
		return nil, err
	}
	ak.Write([]byte(s.apiKey))

	// Writing timestamp
	ts, err := mw.CreateFormField("timestamp")
	if err != nil {
		return nil, err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	ts.Write([]byte(timestamp))

	// Writing signature
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

	fi, err := mw.CreateFormFile("file", fNameWithoutExt)
	if err != nil {
		return nil, err
	}
	fi.Write(file)

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
