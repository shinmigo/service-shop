package rpc

import (
	"context"
	"crypto/md5"
	"fmt"
	"goshop/service-shop/pkg/utils"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/shinmigo/pb/shoppb"
)

const ImagePath = "./static/images"

type Image struct {
}

func NewImage() *Image {
	return &Image{}
}

func (p *Image) Upload(ctx context.Context, req *shoppb.UploadReq) (*shoppb.UploadRes, error) {
	if req.Name == "" || req.Content == nil {
		return nil, fmt.Errorf("上传的文件内容不合法")
	}

	list := strings.Split(req.Name, ".")
	extName := list[len(list)-1]
	randStr := md5.Sum([]byte(time.Now().String() + req.Name))
	name := fmt.Sprintf("%x", randStr) + "." + extName

	if !utils.DirIsExists(ImagePath) {
		if err := os.MkdirAll(ImagePath, 0755); err != nil {
			return nil, fmt.Errorf("上传文件时，创建目录失败， err: %v", err)
		}
	}

	f, err := os.Create(ImagePath + "/" + name)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败, err: %v", err)
	}
	_, _ = f.Write(req.Content)
	_ = f.Close()

	return &shoppb.UploadRes{
		ImageId: name,
	}, nil
}

func (p *Image) GetImage(ctx context.Context, req *shoppb.GetImageReq) (*shoppb.ImageContent, error) {
	if req.ImageId == "" {
		return nil, fmt.Errorf("文件Id不合法")
	}

	f, err := os.Open(ImagePath + "/" + req.ImageId)
	if err != nil {
		return nil, fmt.Errorf("获取文件失败, err: %v", err)
	}
	defer func() {
		_ = f.Close()
	}()

	by, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("读取文件内容失败， err: %v", err)
	}

	return &shoppb.ImageContent{
		Content: by,
		ImageId: req.ImageId,
	}, nil
}
