package build

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
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

func (data *YamlData) Build() {
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
	df, err := os.Create(dst)
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
