package handler

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/hamiddarani/anakonda/internal/api/dto"
	"github.com/hamiddarani/anakonda/internal/api/helper"
	"github.com/hamiddarani/anakonda/internal/api/service"
	"github.com/hamiddarani/anakonda/internal/config"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(cfg *config.Config) *TaskHandler {
	return &TaskHandler{
		service: service.NewTaskService(cfg),
	}
}

// CreateTask godoc
// @Summary Create a Task
// @Description Create a Task
// @Tags Task
// @Accept json
// @produces json
// @Param Request body dto.CreateTaskRequest true "Create a Task"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.CreateTaskResponse} "Task response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/tasks [post]
func (h *TaskHandler) CreateTask(ctx echo.Context) error {
	req := new(dto.CreateTaskRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, helper.BadRequestError, err))
	}

	if errors := helper.ValidateStruct(req); errors != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationErrorCode, errors))
	}

	res, err := h.service.Create(ctx, req)

	if err != nil {
		return ctx.JSON(
			helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err),
		)
	}

	return ctx.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, helper.Success))

}


// GetTask godoc
// @Summary Get Task
// @Description Get Task
// @Tags Task
// @Accept json
// @produces json
// @Param id path string true "Id"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.CreateTaskResponse} "Task response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/tasks/{id} [get]
func (h *TaskHandler)GetById(ctx echo.Context) error {
	id, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, helper.BadRequestError, err))
	}

	res, err := h.service.GetById(ctx, id)

	if err != nil {
		return ctx.JSON(
			helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err),
		)
	}

	return ctx.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, helper.Success))
}
