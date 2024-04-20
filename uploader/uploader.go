package uploader

import (
	"time"
)

type UploadFunc func(uc UploadConfig) (url string, err error)

//go:generate mockgen -source=./uploader.go -destination=./mock/uploader_mock.go -package=mock_uploader

type UploadeResult struct {
	URL      string `json:"url"`
	ScaleURL string `json:"scaleURL"`
}

type Uploader interface {
	Upload() (*UploadeResult, error)
	GetUploadElapsed() time.Duration
}
