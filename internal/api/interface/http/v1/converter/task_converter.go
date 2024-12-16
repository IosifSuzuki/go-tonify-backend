package converter

import (
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/domain/task/model"
	"go-tonify-backend/pkg/datetime"
)

func ConvertModel2TaskResponse(taskModel *model.Task) *dto.Task {
	var task = dto.Task{
		OwnerID:     taskModel.OwnerID,
		Title:       taskModel.Title,
		Description: taskModel.Description,
	}
	if createdAt := taskModel.CreatedAt; createdAt != nil {
		dt := datetime.Datetime(*createdAt)
		task.CreatedAt = &dt
	}
	if updatedAt := taskModel.UpdatedAt; updatedAt != nil {
		dt := datetime.Datetime(*updatedAt)
		task.CreatedAt = &dt
	}
	return &task
}
