package usecase

import (
	"context"
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain/entity"
	"go-tonify-backend/internal/domain/task/converter"
	"go-tonify-backend/internal/domain/task/model"
	"go-tonify-backend/internal/domain/task/repository"
	"go-tonify-backend/pkg/logger"
)

const (
	MaxTasksByAccount int64 = 3
)

type Task interface {
	CreateTask(ctx context.Context, task *model.Task) (*model.Task, error)
	GetList(ctx context.Context, ownerID int64, offset int64, limit int64) ([]model.Task, error)
}

type task struct {
	container      container.Container
	taskRepository repository.Task
}

func NewTask(container container.Container, taskRepository repository.Task) Task {
	return &task{
		container:      container,
		taskRepository: taskRepository,
	}
}

func (t *task) CreateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	log := t.container.GetLogger()
	createdTasks, err := t.taskRepository.CountByID(ctx, task.OwnerID)
	if err != nil {
		log.Error("fail to get count tasks", logger.FError(err))
		return nil, err
	}
	if createdTasks == nil {
		log.Error("createdTasks contains nil value")
		return nil, dto.NilError
	}
	if *createdTasks >= MaxTasksByAccount {
		log.Error("account has exceeded the limit of created tasks")
		return nil, model.CreateTaskLimitError
	}
	taskEntity := entity.Task{
		OwnerID:     task.OwnerID,
		Title:       task.Title,
		Description: task.Description,
	}
	createdTaskID, err := t.taskRepository.Create(ctx, &taskEntity)
	if err != nil {
		log.Error("fail to record task to db", logger.FError(err))
		return nil, err
	}
	if createdTaskID == nil {
		log.Error("createdTaskID contains nil value")
		return nil, dto.NilError
	}
	createdTaskEntity, err := t.taskRepository.GetByID(ctx, *createdTaskID)
	if err != nil {
		log.Error("fail to get created task from db", logger.FError(err))
		return nil, err
	}
	createdTask := converter.ConvertEntity2TaskModel(createdTaskEntity)
	return createdTask, nil
}

func (t *task) GetList(ctx context.Context, ownerID int64, offset int64, limit int64) ([]model.Task, error) {
	log := t.container.GetLogger()
	taskEntities, err := t.taskRepository.GetList(ctx, ownerID, offset, limit)
	if err != nil {
		log.Error("fail to get list task from db", logger.FError(err))
		return nil, err
	}
	tasks := make([]model.Task, 0, len(taskEntities))
	for _, taskEntity := range taskEntities {
		task := converter.ConvertEntity2TaskModel(&taskEntity)
		tasks = append(tasks, *task)
	}
	return tasks, nil
}
