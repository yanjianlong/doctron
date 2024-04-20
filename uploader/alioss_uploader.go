package uploader

import (
	"errors"
	"time"

	"github.com/lampnick/doctron/conf"
	"github.com/lampnick/doctron/pkg/alioss"
)

type AliOssUploader struct {
	DoctronUploader
}

var ErrNoNeedToUpload = errors.New("no need to upload")

func (ins *AliOssUploader) Upload() (*UploadeResult, error) {
	if ins.Key == "" {
		return nil, ErrNoNeedToUpload
	}
	start := time.Now()
	defer func() {
		ins.uploadElapsed = time.Since(start)
	}()
	helper, err := alioss.NewOssHelper(conf.OssConfig)
	if err != nil {
		return nil, err
	}
	uploadUrl, err := helper.Upload(ins.Key, ins.Stream)
	if err != nil {
		return nil, err
	}
	return &UploadeResult{URL: uploadUrl}, nil
}

func (ins *AliOssUploader) GetUploadElapsed() time.Duration {
	return ins.uploadElapsed
}
