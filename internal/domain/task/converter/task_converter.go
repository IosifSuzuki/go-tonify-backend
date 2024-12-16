package converter

import (
	"go-tonify-backend/internal/domain/entity"
	"go-tonify-backend/internal/domain/task/model"
)

func ConvertEntity2TaskModel(taskEntity *entity.Task) *model.Task {
	return &model.Task{
		ID:          taskEntity.ID,
		Title:       taskEntity.Title,
		OwnerID:     taskEntity.OwnerID,
		Description: taskEntity.Description,
		CreatedAt:   taskEntity.CreatedAt,
		UpdatedAt:   taskEntity.UpdatedAt,
	}
}
