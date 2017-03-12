package render

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jeeves/pkg/yamlData"
)

type RenderData struct {
	Data *yamlData.YamlData
}

func (r *RenderData) RenderTemplates() {
	r.SetupBuildir()
	r.RenderDockerfile()
}

func (r *RenderData) SetupBuildir() {
	// Need script basename
	dest_dir := "buildir/" + r.Data.Name
	err := os.MkdirAll(dest_dir, 0777)
	if err != nil {
		fmt.Println(err)
	}
	dest := dest_dir + "/" + r.Data.Script
	r.CopyFile("examples/"+r.Data.Script, dest)
}

func (r *RenderData) CopyFile(src string, dst string) {
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

func (r *RenderData) RenderDockerfile() {
	input, err := ioutil.ReadFile("examples/dockerfiles/" + r.Data.Container)
	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "{{SCRIPT}}") {
			lines[i] = "COPY " + r.Data.Script + " ."
		}
		if strings.Contains(line, "{{COMMAND}}") {
			lines[i] = "CMD ['./" + r.Data.Script + "']"
		}
	}
	output := strings.Join(lines, "\n")
	_, err = os.Create("buildir/" + r.Data.Name + "/Dockerfile")

	err = ioutil.WriteFile("buildir/"+r.Data.Name+"/Dockerfile", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
