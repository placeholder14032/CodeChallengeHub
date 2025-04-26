package judge

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func compile(ctx context.Context, sourcePath, outputPath string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	// Resolve absolute paths
	sourcePath, _ = filepath.Abs(sourcePath)
	outputPath, _ = filepath.Abs(outputPath)

	// Create the container
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "golang:latest",
		Cmd:   []string{"go", "build", "-o", "/app/submission.out", "/app/main.go"},
		Tty:   false,
		WorkingDir: "/app",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: sourcePath,
				Target: "/app/main.go",
			},
			{
				Type:   mount.TypeBind,
				Source: outputPath,
				Target: "/app/submission.out",
			},
		},
	}, nil, nil, "")

	if err != nil {
		return err
	}
	defer func() {
		cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
	}()

	// Pull image if not available
	fmt.Println("Starting container to compile...")
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	// Wait for compilation to finish
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
		fmt.Println("Compilation complete.")
	}

	// check compilation logs
	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return err
	}
	fmt.Println("Compiler output:")
	io.Copy(os.Stdout, out)
	// check exit code
	inspection, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil { return err }
	if inspection.State.ExitCode != 0 { return errors.New("compilation error") }

	return nil
}

