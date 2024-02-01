package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hamiddarani/anakonda/internal/config"
	"github.com/hamiddarani/anakonda/pkg/cache"
	"github.com/hamiddarani/anakonda/pkg/logging"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

type Controller struct {
	cfg    *config.Config
	redis  *redis.Client
	logger logging.Logger
}

var ctx = context.Background()

type TaskPayload struct {
	Id      string
	Runtime string
}

func New(cfg *config.Config, lg logging.Logger) *Controller {
	return &Controller{
		cfg:    cfg,
		redis:  cache.GetRedis(),
		logger: lg,
	}
}

func NewHandleTask(msg string) (*asynq.Task, error) {
	return asynq.NewTask("task:deliverer", []byte(msg)), nil
}

func (c *Controller) watch_dog() {
	for {
		c.redis.ExpireXX(ctx, c.cfg.Controller.ControllerLeaderRedisKey, 1*time.Second)
	}
}

func (c *Controller) InitController() {
	controller_name := fmt.Sprintf("%s%d", c.cfg.Controller.ControllerPrefixKey, time.Now().Unix())
	c.redis.Set(ctx, controller_name, time.Now().Unix(), 0)

	for {
		if c.redis.SetNX(ctx, c.cfg.Controller.ControllerLeaderRedisKey, time.Now().Unix(), time.Second*1).Val() {
			fmt.Println("I'm Leader")
			break
		} else {
			current_leader := c.redis.Get(ctx, c.cfg.Controller.ControllerLeaderRedisKey)
			fmt.Printf("current leader is: %s\n", current_leader)
			time.Sleep(3 * time.Second)
		}
	}

	go c.watch_dog()

	subscriber := c.redis.Subscribe(ctx, c.cfg.Queue.Channels["anakonda_new_task"])
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%d", c.cfg.Redis.Host, c.cfg.Redis.Port), Password: c.cfg.Redis.Password})
	defer client.Close()

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			c.logger.Fatal(logging.Controller, logging.Subscribe, err.Error(), nil)
		}

		var t TaskPayload

		if err := json.Unmarshal([]byte(msg.Payload), &t); err != nil {
			c.logger.Fatal(logging.Controller, logging.UnmarshalTask, err.Error(), nil)
		}

		fmt.Println("Received message from " + msg.Channel + " channel.")
		fmt.Printf("%v", t)

		task, err := NewHandleTask(msg.Payload)
		if err != nil {
			c.logger.Warn(logging.Controller, logging.QueueTask, err.Error(), nil)
		}
		info, err := client.Enqueue(task)
		if err != nil {
			c.logger.Fatalf("could not enqueue task: %v", err)
		}
		c.logger.Infof("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	}

}
