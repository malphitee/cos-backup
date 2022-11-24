package cos_tool

import (
	config2 "backFolderToCos/config"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var client *cos.Client

type CosTool struct {
	Prefix    string
	Delimiter string
	Config    *config2.Config
}

func (t *CosTool) GetCosClient() *cos.Client {
	if client == nil {
		once := &sync.Once{}
		once.Do(func() {
			var err error
			config := t.Config
			u, err := url.Parse(config.Url)
			if err != nil {
				panic(err)
			}
			b := &cos.BaseURL{BucketURL: u}
			client = cos.NewClient(b, &http.Client{
				Transport: &cos.AuthorizationTransport{
					SecretID:  config.SecretId,
					SecretKey: config.SecretKey,
				},
			})
		})
	}
	return client
}

func (t *CosTool) GetBucketFileList() *cos.BucketGetResult {
	var err error
	opt := &cos.BucketGetOptions{
		Prefix:    t.Prefix,
		Delimiter: t.Delimiter,
		Marker:    "",
		MaxKeys:   1000,
	}
	client := t.GetCosClient()
	v, _, err := client.Bucket.Get(context.Background(), opt)
	if err != nil {
		panic(fmt.Sprintf("获取Bucket内容出错,err = %v", err))
	}
	return v
}

func (t *CosTool) UploadToCos(path string, fileName string) error {
	client := t.GetCosClient()
	response, _, err := client.Object.Upload(context.Background(), fileName, path, nil)
	if err != nil {
		fmt.Println("上传文件至COS时出错,文件名=", fileName, " 文件路径=", path, " err = ", err)
	} else {
		fmt.Println("上传文件至COS成功,文件名=", fileName, " md5=", response.ETag, " url=", response.Location)
	}
	return err
}
