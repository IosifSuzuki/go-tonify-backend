package usecase

import (
	"context"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain/category/converter"
	"go-tonify-backend/internal/domain/category/model"
	"go-tonify-backend/internal/domain/category/repository"
	"go-tonify-backend/pkg/logger"
)

type Category interface {
	GetAllCategories(ctx context.Context, offset int64, limit int64) (*model.Pagination[model.Category], error)
}

type category struct {
	container          container.Container
	categoryRepository repository.Category
}

func NewCategory(
	container container.Container,
	categoryRepository repository.Category,
) Category {
	return &category{
		container:          container,
		categoryRepository: categoryRepository,
	}
}

func (c *category) GetAllCategories(ctx context.Context, offset int64, limit int64) (*model.Pagination[model.Category], error) {
	log := c.container.GetLogger()
	categories, err := c.categoryRepository.GetAll(ctx, offset, limit)
	if err != nil {
		log.Error("fail to get all categories", logger.FError(err))
		return nil, err
	}
	numberOfCategories, err := c.categoryRepository.GetAllNumberRows(ctx)
	if err != nil {
		log.Error("fail to get number of all categories", logger.FError(err))
		return nil, err
	}
	if numberOfCategories == nil {
		log.Error("number_of_categories has nil value", logger.FError(err))
		return nil, model.NilError
	}
	categoryModels := converter.ConvertEntities2CategoriesModel(categories)
	pagination := model.Pagination[model.Category]{
		Offset: offset,
		Limit:  limit,
		Total:  *numberOfCategories,
		Data:   categoryModels,
	}
	return &pagination, nil
}
