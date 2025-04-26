package judge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

const (
	imageName = "alpine"
)

// runs code inside a docker container with memory constraints and hopefully kills it before reaching time limit
func runCode(ctx context.Context, executablePath string, timeLimit time.Duration, memoryLimit int64, input, expectedOutput []byte) (RunResults, error) {
	executablePath, _ = filepath.Abs(executablePath)

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return JudgeError, err
	}
	defer cli.Close()

	containerConfig := &container.Config{
		Image:        imageName,
		Cmd:          []string{"/app/submission.out"},
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		OpenStdin:    true,
		WorkingDir:   "/app",
		NetworkDisabled: true,
	}
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: executablePath,
				Target: "/app/submission.out",
			},
		},
		Resources: container.Resources{
			Memory: memoryLimit,
			CPUPeriod: 100_000,
			CPUQuota: 100_000, // exactly one cpu
		},
	}

	// creating the container
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return JudgeError, err
	}
	// deleting the container after we are done with it
	defer func() {
		cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
	}()

	// attaching to the container stdin/stdout
	hijackedResp, err := cli.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stream: true, Stdin: true, Stdout: true, Stderr: true,
	})
	if err != nil {
		return JudgeError, err
	}
	defer hijackedResp.Close()

	// starting the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return JudgeError, err
	}

	// this channel is used to report if the program ends and adds in error
	done := make(chan error, 1)

	// checking memory usage while it's running and recording the maximum usage
	var maxMemory uint64
	statsDone := make(chan struct{})
	go func() {
		stats, err := cli.ContainerStats(ctx, resp.ID, true)
		if err != nil {
			fmt.Println("Error getting live container stats:", err)
			close(statsDone)
			return
		}
		defer stats.Body.Close()
		decoder := json.NewDecoder(stats.Body)
		for decoder.More() {
			var s struct {
				MemoryStats struct {
					Usage uint64 `json:"usage"`
				} `json:"memory_stats"`
			}
			if err := decoder.Decode(&s); err != nil {
				break
			}
			if s.MemoryStats.Usage > maxMemory {
				maxMemory = s.MemoryStats.Usage
			}
		}
		close(statsDone)
	}()
	//

	// Send input
	_, err = hijackedResp.Conn.Write(input)
	//defer hijackedResp.Conn.Close()

	// Read output
	var outputBuf bytes.Buffer

	go func() {
		_, err := stdcopy.StdCopy(&outputBuf, io.Discard, hijackedResp.Reader)
		done <- err
	}()

	select {
	case <-time.After(timeLimit):
		_ = cli.ContainerKill(ctx, resp.ID, "SIGKILL")
		fmt.Println("Time Limit Exceeded")
		return TimeLimit, nil
	case err := <-done:
		if err != nil {
			fmt.Println("Runtime Error:", err)
			return RuntimeError, nil
		}
	}

	actual := bytes.TrimSpace(outputBuf.Bytes())
	expected := bytes.TrimSpace(expectedOutput)

	// checking how it ran
	inspection, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return JudgeError, err
	}
	// Possible values are listed in the `ContainerState` docs; there do not
	// seem to be named constants for these values.
	if inspection.State.Status != "exited" {
		return JudgeError, errors.New("container is not exited")
	}

	if inspection.State.ExitCode != 0 {
		return RuntimeError, nil
	}
	if maxMemory > uint64(memoryLimit) || inspection.State.OOMKilled {
		return MemoryLimit, nil
	}

	if bytes.Equal(actual, expected) {
		return Accepted, nil
	} else {
		return Wrong, nil
	}
}
