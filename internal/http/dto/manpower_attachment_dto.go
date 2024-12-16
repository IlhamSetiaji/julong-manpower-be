package dto

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/spf13/viper"
)

func ConvertManpowerAttachmentsToResponse(manpowerAttachments *[]entity.ManpowerAttachment, viper *viper.Viper) []*response.ManpowerAttachmentResponse {
	var res []*response.ManpowerAttachmentResponse
	for _, manpowerAttachment := range *manpowerAttachments {
		fullURL := viper.GetString("app.url") + manpowerAttachment.FilePath
		res = append(res, &response.ManpowerAttachmentResponse{
			ID:       manpowerAttachment.ID.String(),
			FileName: manpowerAttachment.FileName,
			FilePath: fullURL,
			FileType: manpowerAttachment.FileType,
		})
	}
	return res
}
