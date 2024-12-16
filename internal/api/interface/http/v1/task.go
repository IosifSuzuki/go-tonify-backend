package v1

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/api/interface/http/v1/converter"
	"go-tonify-backend/internal/api/interface/http/validator"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/domain/task/model"
	taskUsecase "go-tonify-backend/internal/domain/task/usecase"
	"go-tonify-backend/pkg/logger"
	"net/http"
)

type TaskHandler struct {
	container   container.Container
	validation  validator.HttpValidator
	taskUsecase taskUsecase.Task
}

func NewTaskHandler(
	container container.Container,
	validation validator.HttpValidator,
	taskUsecase taskUsecase.Task,
) *TaskHandler {
	return &TaskHandler{
		container:   container,
		validation:  validation,
		taskUsecase: taskUsecase,
	}
}

// CreateTask godoc
//
//	@Summary		Create a task
//	@Description	The account must have a client role. Each account has a limit on task creation.
//	@Description	If everything goes well, the server will return the created task as a response
//	@Tags			task
//	@Param			Authorization	header		string					true	"account's access token"
//	@Param			request			body		dto.CreateTask	true	"create task parameters"
//	@Produce		json
//	@Success		201	{object}	dto.Response{response=dto.Task}			"created task"
//	@Failure		400	{object}	dto.Response{response=dto.Empty}		"detailed error message"
//	@Failure		401	{object}	dto.Response{response=dto.Empty}		"the authorization token is invalid/expired/missing"
//	@Failure		403	{object}	dto.Response{response=dto.Empty}		"account has reached the task limit or has an incorrect role"
//	@Failure		500	{object}	dto.Response{response=dto.Empty}		"detailed error message"
//	@Router			/v1/task/create [post]
//	@Security		ApiKeyAuth
func (t *TaskHandler) CreateTask(ctx *gin.Context) {
	log := t.container.GetLogger()
	accountID, err := getAccountID(ctx)
	if err != nil {
		log.Error("fail to get account id", logger.FError(err))
		failResponse(ctx, http.StatusUnauthorized, err, nil)
		return
	}
	var createTask dto.CreateTask
	if err := ctx.ShouldBindJSON(&createTask); err != nil {
		log.Error("fail to bind get match accounts", logger.FError(err))
		badRequestResponse(ctx, t.validation, dto.BadRequestError, err)
		return
	}
	var createTaskModel = model.Task{
		OwnerID:     *accountID,
		Title:       createTask.Title,
		Description: createTask.Description,
	}
	createdTask, err := t.taskUsecase.CreateTask(ctx, &createTaskModel)
	if err != nil {
		log.Error("fail to execute a create task use case", logger.FError(err))
		switch err {
		case model.CreateTaskLimitError:
			failResponse(ctx, http.StatusForbidden, err, dto.CreateTaskLimitError)
		default:
			failResponse(ctx, http.StatusInternalServerError, err, dto.FailProcessRequestError)
		}
		return
	}
	task := converter.ConvertModel2TaskResponse(createdTask)
	successResponse(ctx, http.StatusCreated, task)
}

// GetListTask godoc
//
//	@Summary		List of tasks by account id
//	@Description    Get list of tasks by account ID with pagination parameters
//	@Tags			task
//	@Param			Authorization	header		string					true	"account's access token"
//	@Param			account_id		query		int						true	"account id"
//	@Param			offset			query		int						false	"pagination offset"
//	@Param			limit			query		int						true	"pagination limit"
//	@Produce		json
//	@Success		200	{object}	dto.Response{response=[]dto.Task}		"list of tasks"
//	@Failure		400	{object}	dto.Response{response=dto.Empty}		"detailed error message"
//	@Failure		401	{object}	dto.Response{response=dto.Empty}		"the authorization token is invalid/expired/missing"
//	@Failure		500	{object}	dto.Response{response=dto.Empty}		"detailed error message"
//	@Router			/v1/task/list [get]
//	@Security		ApiKeyAuth
func (t *TaskHandler) GetListTask(ctx *gin.Context) {
	log := t.container.GetLogger()
	var getListTask dto.GetListTask
	if err := ctx.ShouldBindQuery(&getListTask); err != nil {
		log.Error("fail to bind get list task", logger.FError(err))
		badRequestResponse(ctx, t.validation, dto.BadRequestError, err)
		return
	}
	taskModels, err := t.taskUsecase.GetList(ctx, getListTask.AccountID, getListTask.Offset, getListTask.Limit)
	if err != nil {
		log.Error("fail to execute get list usecase", logger.FError(err))
		failResponse(ctx, http.StatusInternalServerError, err, dto.FailProcessRequestError)
		return
	}
	tasks := make([]dto.Task, 0, len(taskModels))
	for _, taskModel := range taskModels {
		task := converter.ConvertModel2TaskResponse(&taskModel)
		tasks = append(tasks, *task)
	}
	successResponse(ctx, http.StatusOK, tasks)
}
