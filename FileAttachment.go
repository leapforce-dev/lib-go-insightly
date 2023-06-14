package insightly

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	i_types "github.com/leapforce-libraries/go_insightly/types"
	"io"
	"net/http"
)

// FileAttachment stores FileAttachment from Service
type FileAttachment struct {
	FileId         int                     `json:"FILE_ID"`
	FileName       string                  `json:"FILE_NAME"`
	ContentType    string                  `json:"CONTENT_TYPE"`
	FileSize       int                     `json:"FILE_SIZE"`
	FileCategoryId int                     `json:"FILE_CATEGORY_ID"`
	OwnerUserId    int                     `json:"OWNER_USER_ID"`
	DateCreatedUtc *i_types.DateTimeString `json:"DATE_CREATED_UTC"`
	DateUpdatedUtc *i_types.DateTimeString `json:"DATE_UPDATED_UTC"`
	Url            string                  `json:"URL"`
}

// GetFileAttachment returns a specific file attachments as a slice of bytes
func (service *Service) GetFileAttachment(fileId int64) ([]byte, *errortools.Error) {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodGet,
		Url:    service.url(fmt.Sprintf("fileattachments/%v", fileId)),
	}

	_, response, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	return b, nil
}
