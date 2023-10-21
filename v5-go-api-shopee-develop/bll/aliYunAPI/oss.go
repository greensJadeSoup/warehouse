package aliYunAPI

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"strings"
	"warehouse/v5-go-api-shopee/conf"
	"warehouse/v5-go-api-shopee/constant"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_util"
)

type OssBLL struct{
	Client  *oss.Client
}

var Oss OssBLL

func (this *OssBLL) init() error {
	var err error

	if this.Client != nil {
		return nil
	}

	this.Client, err = oss.New(
		conf.GetAppConfig().Oss.EndPointUrl,
		conf.GetAppConfig().Oss.AccessKeyID,
		conf.GetAppConfig().Oss.AccessKeySecret)
	if err != nil {
		return cp_error.NewSysError("oss初始化失败:" + err.Error())
	}

	return nil
}

//IsOSSFileExist 上传文件
func (this *OssBLL) UploadPdf(tmpPath string) (string, error) {
	err := this.init()
	if err != nil {
		return "", err
	}

	pdfNewName := cp_util.NewGuid() + ".pdf"
	objectPath := constant.OSS_PATH_SHOPEE_DOCUMENT + "/" + pdfNewName
	url := fmt.Sprintf("https://publice-pdf.%s.aliyuncs.com/%s/%s",
		constant.OSS_REGION_SZ, constant.OSS_PATH_SHOPEE_DOCUMENT, pdfNewName)

	bucket, err := this.Client.Bucket(constant.BUCKET_NAME_PUBLICE_PDF)
	if err != nil {
		return "", cp_error.NewSysError("bucket name非法或不存在:" + err.Error())
	}

	err = bucket.UploadFile(objectPath, tmpPath, 1024 * 1024, []oss.Option{}...)
	if err != nil {
		return "", cp_error.NewSysError("oss文件上传失败:" + err.Error())
	}

	return url, nil
}

func (this *OssBLL) DeleteOSSImage(url string) error {
	if url == "" {
		return nil
	}

	idx := strings.LastIndex(url, constant.OSS_PATH_SHOPEE_DOCUMENT)
	if idx == -1 {
		return cp_error.NewSysError("图片格式错误:" + url)
	}

	return this.DeleteOSSFile(constant.BUCKET_NAME_PUBLICE_PDF, url[idx:])
}

//IsOSSFileExist 判断文件是否存在
func (this *OssBLL) IsOSSFileExist(bucketName, object string) (bool, error) {
	err := this.init()
	if err != nil {
		return false, err
	}

	bucket, err := this.Client.Bucket(bucketName)
	if err != nil {
		return false, cp_error.NewSysError("bucket name非法或不存在:" + err.Error())
	}

	return bucket.IsObjectExist(object)
}

//DeleteOSSFile 删除OSS上面的文件
func (this *OssBLL) DeleteOSSFile(bucketName, object string) error {
	err := this.init()
	if err != nil {
		return err
	}

	bucket, err := this.Client.Bucket(bucketName)
	if err != nil {
		return cp_error.NewSysError("bucket name非法或不存在:" + err.Error())
	}

	err = bucket.DeleteObject(object)
	if err != nil {
		return cp_error.NewSysError("oss文件删除失败:" + err.Error())
	}

	return nil
}

