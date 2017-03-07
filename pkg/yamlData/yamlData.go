package yamlData

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

// Handle user input to a file location
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
