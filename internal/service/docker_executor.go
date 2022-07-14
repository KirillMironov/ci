package service

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
)

// DockerExecutor is a service that can execute a step in a docker container.
type DockerExecutor struct {
	cli        *client.Client
	workingDir string
}

// NewDockerExecutor creates a new DockerExecutor with a provided docker client.
func NewDockerExecutor(cli *client.Client, workingDir string) *DockerExecutor {
	return &DockerExecutor{
		cli:        cli,
		workingDir: workingDir,
	}
}

// ExecuteStep executes a step in a container.
func (de DockerExecutor) ExecuteStep(ctx context.Context, step domain.Step, sourceCodeArchive io.Reader) (
	logs io.ReadCloser, err error) {
	config := &containertypes.Config{
		Image:      step.Image,
		Env:        step.Environment,
		Cmd:        step.Command,
		Tty:        true,
		WorkingDir: de.workingDir,
	}

	_, err = de.cli.ImagePull(ctx, config.Image, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}

	container, err := de.cli.ContainerCreate(ctx, config, nil, nil, nil, "")
	if err != nil {
		return nil, err
	}

	err = de.cli.CopyToContainer(ctx, container.ID, de.workingDir, sourceCodeArchive, types.CopyToContainerOptions{})
	if err != nil {
		return nil, err
	}

	err = de.cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	statusCh, errCh := de.cli.ContainerWait(ctx, container.ID, containertypes.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		if err != nil {
			return nil, err
		}
	case <-statusCh:
	}

	return de.cli.ContainerLogs(ctx, container.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
}
