package uploader

import (
	"time"

	"github.com/lampnick/doctron/conf"
)

type MockUploader struct {
	DoctronUploader
}

func (ins *MockUploader) Upload() (*UploadeResult, error) {
	start := time.Now()
	defer func() {
		ins.uploadElapsed = time.Since(start)
	}()
	res := &UploadeResult{
		URL: "http://" + conf.LoadedConfig.Oss.PrivateServerDomain + "/" + ins.UploadConfig.Key,
	}
	return res, nil
}

func (ins *MockUploader) GetUploadElapsed() time.Duration {
	return ins.uploadElapsed
}
