package converter

import (
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/domain/account/model"
	"go-tonify-backend/pkg/datetime"
)

func ConvertModel2AccountResponse(accountModel *model.Account) *dto.Account {
	account := dto.Account{
		ID:         accountModel.ID,
		TelegramID: accountModel.TelegramID,
		FirstName:  accountModel.FirstName,
		MiddleName: accountModel.MiddleName,
		LastName:   accountModel.LastName,
		Role:       accountModel.Role,
		Nickname:   accountModel.Nickname,
		AboutMe:    accountModel.AboutMe,
		Gender:     accountModel.Gender,
		Country:    accountModel.Country,
		Location:   accountModel.Location,
	}
	if createdAt := accountModel.CreatedAt; createdAt != nil {
		dt := datetime.Datetime(*createdAt)
		account.CreatedAt = &dt
	}
	if updatedAt := accountModel.CreatedAt; updatedAt != nil {
		dt := datetime.Datetime(*updatedAt)
		account.CreatedAt = &dt
	}
	if accountModel.Company != nil {
		company := ConvertModel2CompanyResponse(accountModel.Company)
		account.Company = company
	}
	if accountModel.AvatarAttachment != nil {
		avatarAttachment := ConvertModel2AttachmentResponse(accountModel.AvatarAttachment)
		account.AvatarAttachment = avatarAttachment
	}
	if accountModel.DocumentAttachment != nil {
		documentAttachment := ConvertModel2AttachmentResponse(accountModel.DocumentAttachment)
		account.DocumentAttachment = documentAttachment
	}
	if accountModel.Tags != nil {
		tags := ConvertModels2TagsResponse(*accountModel.Tags)
		account.Tags = &tags
	}
	if accountModel.Categories != nil {
		categories := ConvertModels2CategoriesResponse(*accountModel.Categories)
		account.Categories = &categories
	}
	return &account
}
