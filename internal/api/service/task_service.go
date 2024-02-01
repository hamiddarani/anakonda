package service

import (
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/hamiddarani/anakonda/internal/api/dto"
	"github.com/hamiddarani/anakonda/internal/config"
	"github.com/hamiddarani/anakonda/pkg/cache"
	"github.com/hamiddarani/anakonda/pkg/db"
	"github.com/hamiddarani/anakonda/pkg/db/model"
	"github.com/hamiddarani/anakonda/pkg/logging"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TaskService struct {
	cfg      *config.Config
	database *gorm.DB
	logger   logging.Logger
	redis    *redis.Client
}

func NewTaskService(cfg *config.Config) *TaskService {
	return &TaskService{
		cfg:      cfg,
		database: db.GetDb(),
		logger:   logging.NewLogger(cfg.Logger),
		redis:    cache.GetRedis(),
	}
}

func (s *TaskService) Create(ctx echo.Context, req *dto.CreateTaskRequest) (*dto.CreateTaskResponse, error) {
	task := &model.Task{
		Name:      req.Name,
		Image:     req.Image,
		Namespace: req.Namespace,
		Runtime:   req.Runtime,
		Status:    "new",
		Script:    req.Script,
	}

	tx := s.database.Begin()
	if err := tx.Create(task).Error; err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return nil, err
	}
	tx.Commit()

	payload, err := json.Marshal(task)
	if err != nil {
		s.logger.Error(logging.OS, logging.MarshalTask, err.Error(), nil)
		return nil, err
	}

	cmd := s.redis.Publish(ctx.Request().Context(), s.cfg.Queue.Channels["anakonda_new_task"], payload)
	fmt.Println("------------------------------------")
	fmt.Printf("%v", cmd)
	fmt.Println("------------------------------------")

	return s.GetById(ctx, task.ID)
}

func (s *TaskService) GetById(ctx echo.Context, id uuid.UUID) (*dto.CreateTaskResponse, error) {
	task := new(model.Task)
	if err := s.database.Where("id = ?", id).First(task).Error; err != nil {
		return nil, err
	}
	return &dto.CreateTaskResponse{
		Id:        task.ID,
		Name:      task.Name,
		Image:     task.Image,
		Namespace: task.Namespace,
		Runtime:   task.Runtime,
		Status:    task.Status,
		Script:    task.Script,
		Result:    task.Result,
	}, nil
}
