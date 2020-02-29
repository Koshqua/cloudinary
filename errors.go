package cloudinary

import "errors"

var (
	errNotCloudinary = errors.New("URL scheme is not cloudinary")
	errNoAPISecret   = errors.New("There is no api secret provided in URL")
)
