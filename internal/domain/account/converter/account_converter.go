package converter

import (
	"go-tonify-backend/internal/domain/account/model"
	"go-tonify-backend/internal/domain/entity"
)

func ConvertEntity2AccountModel(accountEntity *entity.Account) *model.Account {
	account := model.Account{
		ID:         accountEntity.ID,
		TelegramID: accountEntity.TelegramID,
		FirstName:  accountEntity.FirstName,
		MiddleName: accountEntity.MiddleName,
		LastName:   accountEntity.LastName,
		Role:       accountEntity.Role.String(),
		Nickname:   accountEntity.Nickname,
		AboutMe:    accountEntity.AboutMe,
		Gender:     accountEntity.Gender.String(),
		Country:    accountEntity.Country,
		Location:   accountEntity.Location,
		CreatedAt:  accountEntity.CreatedAt,
		UpdatedAt:  accountEntity.UpdatedAt,
	}
	if accountEntity.Company != nil {
		company := model.Company{
			ID:          accountEntity.Company.ID,
			Name:        accountEntity.Company.Name,
			Description: accountEntity.Company.Description,
			CreatedAt:   accountEntity.Company.CreatedAt,
			UpdatedAt:   accountEntity.Company.UpdatedAt,
		}
		account.Company = &company
	}
	if accountEntity.AvatarAttachment != nil {
		avatarAttachment := model.Attachment{
			ID:        accountEntity.AvatarAttachment.ID,
			Name:      accountEntity.AvatarAttachment.FileName,
			CreatedAt: accountEntity.AvatarAttachment.CreatedAt,
			UpdatedAt: accountEntity.AvatarAttachment.UpdatedAt,
		}
		if accountEntity.AvatarAttachment.Path != nil {
			avatarAttachment.Path = *accountEntity.AvatarAttachment.Path
		}
		account.AvatarAttachment = &avatarAttachment
	}
	if accountEntity.DocumentAttachment != nil {
		documentAttachment := model.Attachment{
			ID:        accountEntity.DocumentAttachment.ID,
			Name:      accountEntity.DocumentAttachment.FileName,
			CreatedAt: accountEntity.DocumentAttachment.CreatedAt,
			UpdatedAt: accountEntity.DocumentAttachment.UpdatedAt,
		}
		if accountEntity.DocumentAttachment.Path != nil {
			documentAttachment.Path = *accountEntity.DocumentAttachment.Path
		}
		account.DocumentAttachment = &documentAttachment
	}
	return &account
}
