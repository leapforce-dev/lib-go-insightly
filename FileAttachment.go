package insightly

import (
	i_types "github.com/leapforce-libraries/go_insightly/types"
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
