package service

import (
	"context"
	"fmt"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
)

type DockerExecutor struct {
	cli *client.Client
}

func NewDockerExecutor(cli *client.Client) *DockerExecutor {
	return &DockerExecutor{cli: cli}
}

func (de DockerExecutor) Execute(ctx context.Context, step domain.Step, sourceCodePath string) error {
	const workingDir = "/ci"

	config := &containertypes.Config{
		Image:      step.Image,
		Env:        step.Environment,
		Cmd:        step.Command,
		Tty:        true,
		WorkingDir: workingDir,
	}

	logs, err := de.cli.ImagePull(ctx, config.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer logs.Close()

	_, err = io.Copy(os.Stdout, logs)
	if err != nil && err != io.EOF {
		return err
	}

	container, err := de.cli.ContainerCreate(ctx, config, &containertypes.HostConfig{
		Binds: []string{fmt.Sprintf("%s:%s", sourceCodePath, workingDir)},
	}, nil, nil, "")
	if err != nil {
		return err
	}

	err = de.cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	statusCh, errCh := de.cli.ContainerWait(ctx, container.ID, containertypes.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	logs, err = de.cli.ContainerLogs(ctx, container.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
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
