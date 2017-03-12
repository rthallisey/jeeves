/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"fmt"

	"github.com/jeeves/pkg/build"
	"github.com/jeeves/pkg/render"
	"github.com/jeeves/pkg/yamlData"
)

func Run() *int {
	input := new(yamlData.YamlData)
	yamlData.ReadYaml("/home/rhallisey/src/github.com/jeeves/examples/heal.yaml", input)

	r := &render.RenderData{Data: input}
	r.RenderTemplates()

	client, err := build.DockerClient()
	if err != nil {
		fmt.Println(err)
	}

	for {
		fmt.Println("building image")
		build.DockerBuild(client)
		i := new(int)
		return i
	}
}
