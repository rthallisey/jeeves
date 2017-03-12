package build

import (
	"archive/tar"
	"bytes"
	"fmt"
	"time"

	"github.com/fsouza/go-dockerclient"
)

func DockerClient() (*docker.Client, error) {
	endpoint := "unix:///var/run/docker.sock"
	return docker.NewClient(endpoint)
}

func DockerBuild(client *docker.Client) {
	t := time.Now()
	inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputbuf)
	tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: 10, ModTime: t, AccessTime: t, ChangeTime: t})
	tr.Write([]byte("FROM base\n"))
	tr.Close()
	opts := docker.BuildImageOptions{
		Name: "test-jeeves",
		//		InputStream:  inputbuf,
		OutputStream: outputbuf,
		ContextDir:   "buildir/jeeves",
	}

	err := client.BuildImage(opts)
	if err != nil {
		fmt.Println(err)
	}

	hostConfig := docker.HostConfig{}
	createOpts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: "test-jeeves",
		},
		HostConfig: &hostConfig,
	}
	DockerRun(client, hostConfig, createOpts)

}

func DockerRun(client *docker.Client, hostConfig docker.HostConfig, opts docker.CreateContainerOptions) {
	container, err := client.CreateContainer(opts)
	if err != nil {
		fmt.Println(err)
	}
	err = client.StartContainer(container.ID, &hostConfig)
	if err != nil {
		fmt.Println(err)
	}
}
