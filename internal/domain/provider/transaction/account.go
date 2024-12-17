package transaction

import (
	"database/sql"
	"go-tonify-backend/internal/domain/account/repository"
	"go-tonify-backend/pkg/psql"
)

type Provider struct {
	db *sql.DB
}

type ComposedRepository struct {
	Attachment repository.Attachment
	Account    repository.Account
	Company    repository.Company
	Tag        repository.Tag
}

func NewProvider(db *sql.DB) *Provider {
	return &Provider{
		db: db,
	}
}

func (p *Provider) Transact(txFunc func(composed ComposedRepository) error) error {
	return psql.RunInTx(p.db, func(tx *sql.Tx) error {
		composed := ComposedRepository{
			Attachment: repository.NewAttachment(tx),
			Account:    repository.NewAccount(tx),
			Company:    repository.NewCompany(tx),
			Tag:        repository.NewTag(tx),
		}
		return txFunc(composed)
	})
}
