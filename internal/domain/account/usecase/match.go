package usecase

import (
	"context"
	"database/sql"
	"go-tonify-backend/internal/container"
	accountConverter "go-tonify-backend/internal/domain/account/converter"
	"go-tonify-backend/internal/domain/account/model"
	accountRepository "go-tonify-backend/internal/domain/account/repository"
	categoryConverter "go-tonify-backend/internal/domain/category/converter"
	categoryRepository "go-tonify-backend/internal/domain/category/repository"
	"go-tonify-backend/internal/domain/entity"
	"go-tonify-backend/internal/domain/provider/transaction"
	"go-tonify-backend/pkg/logger"
)

type Match interface {
	MatchableAccounts(ctx context.Context, accountID int64, limit int64) ([]model.Account, error)
	MatchAction(ctx context.Context, accountID int64, targetID int64, action model.MatchAction) (model.MatchResult, error)
}

type match struct {
	container           container.Container
	transactionProvider *transaction.Provider
	accountRepository   accountRepository.Account
	tagRepository       accountRepository.Tag
	categoryRepository  categoryRepository.Category
}

func NewMatch(
	container container.Container,
	transactionProvider *transaction.Provider,
	accountRepository accountRepository.Account,
	tagRepository accountRepository.Tag,
	categoryRepository categoryRepository.Category,
) Match {
	return &match{
		container:           container,
		transactionProvider: transactionProvider,
		accountRepository:   accountRepository,
		tagRepository:       tagRepository,
		categoryRepository:  categoryRepository,
	}
}

func (m *match) MatchableAccounts(ctx context.Context, accountID int64, limit int64) ([]model.Account, error) {
	log := m.container.GetLogger()
	accountEntity, err := m.accountRepository.GetByID(ctx, accountID)
	if err != nil {
		log.Error("fail to get account by id", logger.FError(err))
		switch err {
		case sql.ErrNoRows:
			return nil, model.EntityNotFoundError
		default:
			return nil, err
		}
	}
	if err := m.accountRepository.DeleteDislikes(ctx, accountID, 1); err != nil {
		log.Error("fail to clear dislikes", logger.FError(err))
		return nil, err
	}
	accountEntities, err := m.accountRepository.GetMatchableAccounts(ctx, accountEntity.ID, accountEntity.Role.Opposite(), limit)
	if err != nil {
		log.Error("fail to get matchable accounts", logger.FError(err))
		return nil, err
	}
	accounts := make([]model.Account, 0, len(accountEntities))
	for _, accountEntity := range accountEntities {
		account := accountConverter.ConvertEntity2AccountModel(&accountEntity)
		tags, err := m.tagRepository.GetTagsByAccountID(ctx, accountEntity.ID)
		if err != nil {
			log.Error(
				"fail to get tags by account id",
				logger.FError(err),
				logger.F("account_id", accountEntity.ID),
			)
		}
		tagModels := accountConverter.ConvertEntities2TagModels(tags)
		account.Tags = &tagModels
		categories, err := m.categoryRepository.GetCategoriesByAccountID(ctx, accountEntity.ID)
		categoryModels := categoryConverter.ConvertEntities2CategoriesModel(categories)
		account.Categories = &categoryModels
		accounts = append(accounts, *account)
	}
	return accounts, nil
}

func (m *match) MatchAction(ctx context.Context, accountID int64, targetID int64, action model.MatchAction) (model.MatchResult, error) {
	log := m.container.GetLogger()
	dislikeAccount := entity.DislikeAccount{
		DislikerID: accountID,
		DislikedID: targetID,
	}
	likeAccount := entity.LikeAccount{
		LikerID: accountID,
		LikedID: targetID,
	}
	var matchResult model.MatchResult
	err := m.transactionProvider.Transact(func(composed transaction.ComposedRepository) error {
		switch action {
		case model.LikeMatchAction:
			err := composed.Account.DeleteDislikeAccount(ctx, dislikeAccount)
			if err != nil {
				log.Error("fail to delete dislike account", logger.FError(err))
				matchResult = model.ErrorMatchResult
				return err
			}
			exists, err := composed.Account.ExistsLike(ctx, likeAccount)
			if err != nil {
				log.Error("fail to perform exists like", logger.FError(err))
				matchResult = model.ErrorMatchResult
				return err
			}
			if exists {
				log.Error("like is already exists")
				matchResult = model.ErrorMatchResult
				return model.DuplicateMatchActionError
			}
			if err := composed.Account.LikeAccount(ctx, likeAccount); err != nil {
				log.Error("fail to perform like account", logger.FError(err))
				matchResult = model.ErrorMatchResult
				return err
			}
			likeAccount = entity.LikeAccount{
				LikerID: targetID,
				LikedID: accountID,
			}
			exists, err = composed.Account.ExistsLike(ctx, likeAccount)
			if err != nil {
				log.Error("fail to perform exists like", logger.FError(err))
				matchResult = model.ErrorMatchResult
				return err
			}
			if exists {
				matchResult = model.MatchAccountMatchResult
			} else {
				matchResult = model.LikeMatchResult
			}
		case model.DislikeMatchAction:
			err := composed.Account.DeleteLikeAccount(ctx, likeAccount)
			if err != nil {
				log.Error("fail to delete like account", logger.FError(err))
				matchResult = model.ErrorMatchResult
				return err
			}
			exists, err := composed.Account.ExistsDislike(ctx, dislikeAccount)
			if err != nil {
				log.Error("fail to perform exists dislike", logger.FError(err))
				matchResult = model.ErrorMatchResult
				return err
			}
			if exists {
				log.Error("dislike is already exists")
				matchResult = model.ErrorMatchResult
				return model.DuplicateMatchActionError
			}
			if err := composed.Account.DislikeAccount(ctx, dislikeAccount); err != nil {
				log.Error("fail to perform dislike account", logger.FError(err))
				matchResult = model.ErrorMatchResult
				return err
			}
			matchResult = model.DislikeMatchResult
		default:
			log.Error("unknown or unhandled match action")
			matchResult = model.ErrorMatchResult
			return model.UnhandledMatchActionError
		}
		return nil
	})
	if err != nil {
		log.Error("fail to execute db transaction for match action", logger.FError(err))
		return model.ErrorMatchResult, err
	}
	return matchResult, nil
}
