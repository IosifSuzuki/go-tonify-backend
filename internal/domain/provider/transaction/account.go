package transaction

import (
	"database/sql"
	accountRepository "go-tonify-backend/internal/domain/account/repository"
	categoryRepository "go-tonify-backend/internal/domain/category/repository"
	"go-tonify-backend/pkg/psql"
)

type Provider struct {
	db *sql.DB
}

type ComposedRepository struct {
	Attachment accountRepository.Attachment
	Account    accountRepository.Account
	Company    accountRepository.Company
	Tag        accountRepository.Tag
	Category   categoryRepository.Category
}

func NewProvider(db *sql.DB) *Provider {
	return &Provider{
		db: db,
	}
}

func (p *Provider) Transact(txFunc func(composed ComposedRepository) error) error {
	return psql.RunInTx(p.db, func(tx *sql.Tx) error {
		composed := ComposedRepository{
			Attachment: accountRepository.NewAttachment(tx),
			Account:    accountRepository.NewAccount(tx),
			Company:    accountRepository.NewCompany(tx),
			Tag:        accountRepository.NewTag(tx),
			Category:   categoryRepository.NewCategory(tx),
		}
		return txFunc(composed)
	})
}
