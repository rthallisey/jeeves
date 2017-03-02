package build

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type YamlData struct {
	Name      string `yaml:"name"`
	App       string `yaml:"app"`
	Template  string `yaml:"template"`
	Container string `yaml:"container"`
	Script    string `yaml:"script"`
	When      string `yaml:"when"`
}

func NewYamlData() *YamlData {
	return &YamlData{}
}

func (data *YamlData) BuildContainer() {
	data.ReadYaml()
	data.SetupBuildir()
	data.RenderDockerfile()
	data.Build()
}

func (data *YamlData) ReadYaml() {
	yamlFile, err := ioutil.ReadFile("/home/rhallisey/src/github.com/jeeves/examples/heal.yaml")
	if err != nil {
		fmt.Println("yamlFile.Get error: #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, data)
	if err != nil {
		fmt.Println("Unmarshal: %v", err)
	}
}

func (data *YamlData) SetupBuildir() {
	// Need scipt basename
	dest_dir := "buildir/" + data.Name
	err := os.MkdirAll(dest_dir, 0777)
	if err != nil {
		fmt.Println(err)
	}
	dest := dest_dir + "/" + data.Script
	CopyFile("examples/"+data.Script, dest)

}

func CopyFile(src string, dst string) {
	// Need some serious error checking
	sf, err := os.OpenFile(src, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(sf)
	}
	fmt.Println(dst)
	df, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		fmt.Println(df)
	}
	fmt.Println("copy")
	fmt.Println(sf)
	fmt.Println(df)
	var b int64
	if b, err = io.Copy(df, sf); err != nil {
		fmt.Println(err)
		fmt.Println("Copy error")
	}
	fmt.Println(b)
}

func (data *YamlData) RenderDockerfile() {
	input, err := ioutil.ReadFile("examples/dockerfiles/" + data.Container)
	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "{{SCRIPT}}") {
			lines[i] = "COPY " + data.Script + " ."
		}
		if strings.Contains(line, "{{COMMAND}}") {
			lines[i] = "CMD ['./" + data.Script + "']"
		}
	}
	output := strings.Join(lines, "\n")
	dst, err := os.Create("buildir/" + data.Name + "/Dockerfile")
	fmt.Println(dst)

	err = ioutil.WriteFile("buildir/"+data.Name+"/Dockerfile", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (data *YamlData) Build() {
	// Docker build
	cmdName := "docker"
	cmdArgs := []string{"build", "-t", data.Container + "-" + data.Name, "buildir" + data.Name}

	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		fmt.Println(os.Stderr, "Error", err)
	}
	fmt.Println(string(cmdOut))
}
