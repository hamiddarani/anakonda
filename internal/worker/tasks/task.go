package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/hamiddarani/anakonda/pkg/db"
	"github.com/hamiddarani/anakonda/pkg/db/model"
	"github.com/hibiken/asynq"
)

type TaskPayload struct {
	Id      string
	Image   string
	Runtime string
	Script  string
	Status  string
}

func HandleExecuteTask(ctx context.Context, t *asynq.Task) error {
	db := db.GetDb()

	var p TaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	if p.Runtime != "docker" {
		db.Model(&model.Task{}).Where("id = ?", p.Id).Update("status", "failed")
		return fmt.Errorf("container runtime is not valid: %s", p.Runtime)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("could not to connect to %s", p.Runtime)
	}
	defer cli.Close()

	out, err := cli.ImagePull(ctx, p.Image, types.ImagePullOptions{})
	if err != nil {
		db.Model(&model.Task{}).Where("id = ?", p.Id).Update("status", "failed")
		return fmt.Errorf("could not to pull image %s", p.Image)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{Image: p.Image}, nil, nil, nil, "")
	if err != nil {
		db.Model(&model.Task{}).Where("id = ?", p.Id).Update("status", "failed")
		return fmt.Errorf("could not to create container %s", p.Image)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		db.Model(&model.Task{}).Where("id = ?", p.Id).Update("status", "failed")
		return fmt.Errorf("could not to start container %s", p.Image)
	}

	stats, _ := cli.ContainerStats(ctx, resp.ID, false)
	var containerStats map[string]interface{}
	json.NewDecoder(stats.Body).Decode(&containerStats)

	if out, err = cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true}); err != nil {
		db.Model(&model.Task{}).Where("id = ?", p.Id).Update("status", "failed")
		return fmt.Errorf("could not to read container log  %s", resp.ID)
	}
	log, _ := io.ReadAll(out)
	taskUpdate := map[string]interface{}{
		"status": "success",
		"result": log,
	}
	db.Model(&model.Task{}).Where("id = ?", p.Id).Updates(taskUpdate)

	cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})

	return nil
}
