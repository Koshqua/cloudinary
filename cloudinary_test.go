package cloudinary

import (
	"net/url"
	"testing"
)

// func Example() {
// 	//Dial takes a usual cloudinary auth link and returns *Service and error.
// 	service, err := initService("cloudinary://api_key:api_secret@cloud_name")
// 	if err != nil {
// 		//Do something with error
// 	}
// 	//After initialisation we receive a service which is ready to use.
// 	fmt.Println(service)
// }

func TestInitService(t *testing.T) {
	t.Run("Sample test case", func(t *testing.T) {
		service, err := initService("cloudinary://api_key:api_secret@cloud_name")
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		upURL, err := url.Parse("https://api.cloudinary.com/v1_1/cloud_name/image/upload/")
		if err != nil {
			t.Errorf("Unexpected error:%v", err)
		}
		admURL, err := url.Parse("https://api_key:api_secret@api.cloudinary.com/v1_1/cloud_name")
		if err != nil {
			t.Errorf("Unexpected error:%v", err)
		}
		testService := &Service{"cloud_name", "api_key", "api_secret", upURL, admURL, false, false, 0}
		testAdmLink := testService.adminURL.String()
		serviceAdmLink := service.adminURL.String()
		if testAdmLink != serviceAdmLink {
			t.Errorf("Not equal adminURL %v, %v", testAdmLink, serviceAdmLink)
		}
	})
	t.Run("No host provided", func(t *testing.T) {
		_, err := initService("https:://api_key:api_secret@cloud_name")
		if err != errNotCloudinary {
			t.Errorf("Expected error %v, got %v", errNotCloudinary, err)
		}
	})
}
