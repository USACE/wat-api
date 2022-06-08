package utils

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/usace/wat-api/model"
	"golang.org/x/net/context"
)

func StartContainer(plugin model.Plugin, payloadPath string, environmentVariables []string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}
	cli.NegotiateAPIVersion(ctx)
	reader, err := cli.ImagePull(ctx, plugin.ImageAndTag, types.ImagePullOptions{})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	io.Copy(os.Stdout, reader)
	var chc *container.HostConfig
	var nnc *network.NetworkingConfig
	var vp *v1.Platform
	args := append(plugin.CommandLineArgs, payloadPath)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        plugin.ImageAndTag,
		Cmd:          args,
		Tty:          true,
		AttachStdout: true,
		Env:          environmentVariables,
	}, chc, nnc, vp, "")
	if err != nil {
		return "", err
	}
	//retrieve container messages and parrot to lambda standard out.
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, Follow: true})
	if err != nil {
		return "", err
	}
	//defer out.Close()
	io.Copy(os.Stdout, out)
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	statuschn, errchn := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errchn:
		if err != nil {
			log.Fatal(err)
		}
	case status := <-statuschn:
		log.Printf("status.StatusCode: %#+v\n", status.StatusCode)
		//cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	}

	return resp.ID, err
}
