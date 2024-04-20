package uploader

import (
	"fmt"
	"time"

	"github.com/lampnick/doctron/conf"
	"github.com/lampnick/doctron/pkg/alioss"
)

type TencentOssUploader struct {
	DoctronUploader
}

func (ins *TencentOssUploader) Upload() (*UploadeResult, error) {
	if ins.Key == "" {
		return nil, ErrNoNeedToUpload
	}
	start := time.Now()
	defer func() {
		ins.uploadElapsed = time.Since(start)
	}()
	helper, err := alioss.NewTOssHelper(conf.OssConfig)
	fmt.Println("TencentOssUploader", helper, err, conf.OssConfig)
	if err != nil {
		return nil, err
	}
	u1, u2, err := helper.Upload(ins.Key, ins.Stream)
	if err != nil {
		return nil, err
	}
	return &UploadeResult{URL: u1, ScaleURL: u2}, nil
}

func (ins *TencentOssUploader) GetUploadElapsed() time.Duration {
	return ins.uploadElapsed
}
