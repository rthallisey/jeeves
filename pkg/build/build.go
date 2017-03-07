package build

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/jeeves/pkg/yamlData"
)

type Build struct {
	data *yamlData.YamlData
}

func NewBuild() *Build {
	d := yamlData.NewYamlData()
	return &Build{data: d}
}

func (build Build) BuildContainer() {
	build.data.ReadYaml()
	build.SetupBuildir()
	build.RenderDockerfile()
	build.DockerBuild()
}

func (build Build) SetupBuildir() {
	// Need script basename
	dest_dir := "buildir/" + build.data.Name
	err := os.MkdirAll(dest_dir, 0777)
	if err != nil {
		fmt.Println(err)
	}
	dest := dest_dir + "/" + build.data.Script
	build.CopyFile("examples/"+build.data.Script, dest)

}

func (build Build) CopyFile(src string, dst string) {
	// Need some serious error checking
	sf, err := os.OpenFile(src, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(sf)
	}
	df, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		fmt.Println(df)
	}
	if _, err = io.Copy(df, sf); err != nil {
		fmt.Println(err)
		fmt.Println("Copy error")
	}
}

func (build Build) RenderDockerfile() {
	input, err := ioutil.ReadFile("examples/dockerfiles/" + build.data.Container)
	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "{{SCRIPT}}") {
			lines[i] = "COPY " + build.data.Script + " ."
		}
		if strings.Contains(line, "{{COMMAND}}") {
			lines[i] = "CMD ['./" + build.data.Script + "']"
		}
	}
	output := strings.Join(lines, "\n")
	dst, err := os.Create("buildir/" + build.data.Name + "/Dockerfile")
	fmt.Println(dst)

	err = ioutil.WriteFile("buildir/"+build.data.Name+"/Dockerfile", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (build Build) DockerBuild() {
	// Docker build
	cmdName := "docker"
	cmdArgs := []string{"build", "-t", build.data.Container + "-" + build.data.Name, "buildir" + build.data.Name}

	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		fmt.Println(os.Stderr, "Error", err)
	}
	fmt.Println(string(cmdOut))
}
