package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
)

// DockerExecutor used to execute a step in a container.
type DockerExecutor struct {
	cli *client.Client
	// Container working directory.
	workingDir string
}

// NewDockerExecutor creates a new DockerExecutor with a provided docker client.
func NewDockerExecutor(cli *client.Client, workingDir string) *DockerExecutor {
	return &DockerExecutor{
		cli:        cli,
		workingDir: workingDir,
	}
}

// ExecuteStep copies the source code to the container, executes the step and returns container logs.
func (de DockerExecutor) ExecuteStep(ctx context.Context, step domain.Step, sourceCodeArchive io.Reader) (
	logs io.ReadCloser, err error) {
	config := &containertypes.Config{
		Image:      step.Image,
		Env:        step.Environment,
		Entrypoint: step.Command,
		Cmd:        step.Args,
		Tty:        true,
		WorkingDir: de.workingDir,
	}

	pullLogs, err := de.cli.ImagePull(ctx, config.Image, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	defer pullLogs.Close()
	_, _ = io.Copy(io.Discard, pullLogs)

	container, err := de.cli.ContainerCreate(ctx, config, nil, nil, nil, "")
	if err != nil {
		return nil, err
	}

	err = de.cli.CopyToContainer(ctx, container.ID, de.workingDir, sourceCodeArchive, types.CopyToContainerOptions{})
	if err != nil {
		return logs, err
	}

	err = de.cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		return logs, err
	}

	logs, err = de.cli.ContainerLogs(ctx, container.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return nil, err
	}

	resultCh, errCh := de.cli.ContainerWait(ctx, container.ID, containertypes.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		return logs, err
	case result := <-resultCh:
		switch {
		case result.Error != nil:
			return logs, errors.New(result.Error.Message)
		case result.StatusCode != 0:
			return logs, fmt.Errorf("exit code: %d", result.StatusCode)
		default:
			return logs, nil
		}
	}
}
