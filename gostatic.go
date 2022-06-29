package gostatic

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/go-git/go-git/v5"
)

func Build(repo string, directory string, image string, cmd []string, ctx context.Context) error {
	if repo == "" {
		return errors.New("No repository URL given")
	}
	if directory == "" {
		return errors.New("No target directory given")
	}

	r, err := git.PlainOpen(directory)

	if err == git.ErrRepositoryNotExists {
		_, err := git.PlainClone(directory, false, &git.CloneOptions{
			URL:      repo,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
	} else {
		w, err := r.Worktree()
		if err != nil {
			return err
		}
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return err
		}
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      image,
		Cmd:        cmd,
		Tty:        false,
		WorkingDir: "/repo",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/home/eric/src/gostatic/cmd/gostatic/example",
				Target: "/repo",
			},
		},
	}, nil, nil, "")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return err
	}

	if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, out); err != nil {
		return err
	}

	statusCh, _ := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	exit := <-statusCh
	if exit.StatusCode > 0 {
		return errors.New(fmt.Sprintf("exited with code %v", exit.StatusCode))
	}

	return nil
}
