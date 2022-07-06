package service

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
)

type Executor struct {
	cli *client.Client
}

func NewExecutor(cli *client.Client) *Executor {
	return &Executor{cli: cli}
}

func (e *Executor) Execute(ctx context.Context, step domain.Step) error {
	config := &containertypes.Config{
		Image: step.Image,
		Env:   step.Environment,
		Cmd:   step.Command,
		Tty:   true,
	}

	logs, err := e.cli.ImagePull(ctx, config.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer logs.Close()

	_, err = io.Copy(os.Stdout, logs)
	if err != nil && err != io.EOF {
		return err
	}

	container, err := e.cli.ContainerCreate(ctx, config, nil, nil, nil, "")
	if err != nil {
		return err
	}

	err = e.cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	statusCh, errCh := e.cli.ContainerWait(ctx, container.ID, containertypes.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	logs, err = e.cli.ContainerLogs(ctx, container.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return err
	}
	defer logs.Close()

	_, err = io.Copy(os.Stdout, logs)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}
