package helper

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"oauth2-console-go/pkg/er"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
)

const (
	ImgTypeJpeg = ".jpg"
	ImgTypePng  = ".png"
	MB          = 1 << (10 * 2)
)

func CheckFormUploadImage(c *gin.Context, fieldName string, fileLimit int64) (multipart.File, string, string, error) {
	_ = c.Request.ParseMultipartForm(fileLimit * MB)
	file, fileHeader, err := c.Request.FormFile(fieldName)
	if err != nil {
		getFileErr := er.NewAppErr(http.StatusBadRequest, er.UploadFileErrUnknown, err.Error(), err)
		return nil, "", "", getFileErr
	}
	if file == nil {
		fileNilErr := er.NewAppErr(http.StatusBadRequest, er.UploadFileErrNotExist, "the file is empty or isn't exist.", nil)
		return nil, "", "", fileNilErr
	}
	if fileHeader.Size > (fileLimit * MB) {
		reqErr := er.NewAppErr(http.StatusBadRequest, er.UploadFileErrSizeOverLimit, fmt.Sprintf("max file size is %dMB.", fileLimit), nil)
		return nil, "", "", reqErr
	}
	if len(fileHeader.Filename) > 50 {
		reqErr := er.NewAppErr(http.StatusBadRequest, er.UploadFileErrFileName, "max file name is 50 characters.", nil)
		return nil, "", "", reqErr
	}
	mime, err := mimetype.DetectReader(file)
	if err != nil || mime == nil {
		fileNilErr := er.NewAppErr(http.StatusBadRequest, er.UploadFileErrEmpty, "the file is empty or isn't exist.", nil)
		return nil, "", "", fileNilErr
	}
	fileExtension := mime.Extension()
	switch fileExtension {
	case ImgTypeJpeg:
	case ImgTypePng:
	default:
		typeErr := er.NewAppErr(http.StatusBadRequest, er.UploadFileErrTypeNotMatch, "file type not allowed.", nil)
		_ = c.Error(typeErr)
		return nil, "", "", typeErr
	}
	defer file.Close()

	return file, fileHeader.Filename, fileExtension, nil
}

func CheckFormUploadCsv(c *gin.Context, fieldName string, fileLimit int64) (multipart.File, string, error) {
	_ = c.Request.ParseMultipartForm(fileLimit * MB)
	file, fileHeader, err := c.Request.FormFile(fieldName)
	if err != nil {
		getFileErr := er.NewAppErr(http.StatusBadRequest, er.UploadFileErrUnknown, err.Error(), err)
		return nil, "", getFileErr
	}
	if file == nil {
		fileNilErr := er.NewAppErr(http.StatusBadRequest, er.UploadFileErrNotExist, "the file is empty or isn't exist.", nil)
		return nil, "", fileNilErr
	}
	if fileHeader.Size > (fileLimit * MB) {
		reqErr := er.NewAppErr(http.StatusBadRequest, er.UploadFileErrSizeOverLimit, fmt.Sprintf("max file size is %dMB.", fileLimit), nil)
		return nil, "", reqErr
	}
	if len(fileHeader.Filename) > 50 {
		reqErr := er.NewAppErr(http.StatusBadRequest, er.UploadFileErrFileName, "max file name is 50 characters.", nil)
		return nil, "", reqErr
	}
	defer file.Close()

	return file, fileHeader.Filename, nil
}
