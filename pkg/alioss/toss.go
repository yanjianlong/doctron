package alioss

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"net/http"
	"net/url"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gopkg.in/go-playground/validator.v9"
)

type TOssHelper struct {
	client *cos.Client
	config OssConfig
}

func NewTOssHelper(c OssConfig, options ...oss.ClientOption) (*TOssHelper, error) {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		return nil, errors.New("toss uploader config not set")
	}
	u, _ := url.Parse(c.Endpoint)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: c.AccessKeyId, // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: c.AccessKeySecret, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
		},
	})

	return &TOssHelper{
		client: client,
		config: c,
	}, nil
}

func (h TOssHelper) Upload(objectKey string, b []byte, options ...oss.Option) (string, string, error) {
	// h.UploadImageMogr2(objectKey, b)
	parts := strings.Split(objectKey, "/")
	thumbPart := parts[len(parts)-1]

	// text := base64.StdEncoding.EncodeToString([]byte("西瓜新房"))
	pic := &cos.PicOperations{
		IsPicInfo: 1,
		Rules: []cos.PicOperationsRules{
			// {
			// 	FileId: "water_" + thumbPart,
			// 	Rule:   fmt.Sprintf("watermark/2/text/%s/gravity/center", text),
			// },
			{
				FileId: "thum30_" + thumbPart,
				Rule:   "imageMogr2/thumbnail/!30p",
			},
		},
	}
	opt := &cos.ObjectPutOptions{}
	opt.ObjectPutHeaderOptions = &cos.ObjectPutHeaderOptions{
		XOptionHeader: &http.Header{},
	}
	opt.XOptionHeader.Add("Pic-Operations", cos.EncodePicOperations(pic))
	// 获取存储空间。
	reader := bytes.NewReader(b)
	imgresult, resp, err := h.client.CI.Put(context.Background(), objectKey, reader, opt)
	fmt.Println("Upload", imgresult, resp, err)
	if err != nil {
		return "", "", err
	}
	return objectKey, imgresult.ProcessResults[0].Key, nil
}
