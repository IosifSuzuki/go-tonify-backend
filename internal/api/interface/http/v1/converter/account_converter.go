package converter

import (
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/domain/account/model"
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
		CreatedAt:  accountModel.CreatedAt,
		UpdatedAt:  accountModel.UpdatedAt,
	}
	if accountModel.Company != nil {
		company := dto.Company{
			ID:          accountModel.Company.ID,
			Name:        &accountModel.Company.Name,
			Description: &accountModel.Company.Description,
			CreatedAt:   accountModel.Company.CreatedAt,
			UpdatedAt:   accountModel.Company.UpdatedAt,
		}
		account.Company = &company
	}
	if accountModel.AvatarAttachment != nil {
		avatarAttachment := dto.Attachment{
			ID:        &accountModel.AvatarAttachment.ID,
			Name:      &accountModel.AvatarAttachment.Name,
			Path:      &accountModel.AvatarAttachment.Path,
			CreatedAt: accountModel.AvatarAttachment.CreatedAt,
			UpdatedAt: accountModel.AvatarAttachment.UpdatedAt,
		}
		account.AvatarAttachment = &avatarAttachment
	}
	if accountModel.DocumentAttachment != nil {
		documentAttachment := dto.Attachment{
			ID:        &accountModel.DocumentAttachment.ID,
			Name:      &accountModel.DocumentAttachment.Name,
			Path:      &accountModel.DocumentAttachment.Path,
			CreatedAt: accountModel.DocumentAttachment.CreatedAt,
			UpdatedAt: accountModel.DocumentAttachment.UpdatedAt,
		}
		account.DocumentAttachment = &documentAttachment
	}
	return &account
}
