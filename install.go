package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	confilePath = "github.com/arsiac/go-tools-installer/resource.json"
)

//Depandence Tool 安装需要的依赖包
type Depandence struct {
	Address string `json:"address"`
	Path    string `json:"path"`
}

//Tool tool的详细数据
type Tool struct {
	Address     string       `json:"address"`
	Path        string       `json:"path"`
	Depandences []Depandence `json:"depandences"`
	Install     []string     `json:"install"`
}

func getConfig(path string) (tools []Tool, err error) {
	configInfo, err := os.Stat(path)

	if err != nil {
		return tools, err
	}

	if configInfo.Size() == 0 {
		return tools, errors.New("config file is empty")
	}

	confileContent := make([]byte, configInfo.Size())
	confile, err := os.Open(path)
	if err != nil {
		return tools, err
	}

	count, err := confile.Read(confileContent)
	if int64(count) != configInfo.Size() {
		return tools, errors.New(fmt.Sprint("count[", count, "] < file size[", configInfo.Size(), "]"))
	}

	if err != nil {
		return tools, err
	}

	err = json.Unmarshal(confileContent, &tools)

	if err != nil {
		return tools, err
	}

	return tools, nil
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("not exist:", path)
		} else {
			fmt.Println(err.Error())
		}
		return false
	}

	// fmt.Println("exist:", path)
	return true
}

func genAbPath(gopath, path string) string {
	return fmt.Sprint(gopath, "/src/", path)
}

func download(address, path string) {
	fmt.Println("download:", address, "to", path)
	cmd := exec.Command("git", "clone", address, path)
	output, err := cmd.CombinedOutput()

	if err != nil {
		if len(output) > 0 {
			fmt.Println(string(output))
		}
		log.Fatal(err)
	}
}

func install(tool Tool, gopath string) {
	toolPath := genAbPath(gopath, tool.Path)

	if !isExist(toolPath) {
		download(tool.Address, genAbPath(gopath, tool.Path))
	}

	for _, v := range tool.Depandences {
		if !isExist(genAbPath(gopath, v.Path)) {
			download(v.Address, genAbPath(gopath, v.Path))
		}
	}

	for _, v := range tool.Install {
		cmd := exec.Command("go", "get", "-v", v)
		fmt.Println("install:", v)
		output, err := cmd.CombinedOutput()

		if err != nil {
			if len(output) > 0 {
				fmt.Println(string(output))
			}
			log.Fatal(err)
		}
	}
}

func main() {
	gopath := os.Getenv("GOPATH")
	tools, err := getConfig(genAbPath(gopath, confilePath))

	if err != nil {
		log.Fatal(err)
	}

	for _, v := range tools {
		install(v, gopath)
	}
	fmt.Println("SUCCESS")
}
