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
	confilePath = "resource.json"
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

func download(tool Tool) {
	var cmd *exec.Cmd
	cmd = exec.Command("git", "clone", tool.Address, tool.Path)
	result, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(result))

	for _, v := range tool.Depandences {
		cmd = exec.Command("git", "clone", v.Address, v.Path)
		result, err = cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(result))
	}
}

func install(tool Tool) {
	_, err := os.Stat(fmt.Sprint(os.Getenv("GOPATH"), "/src/", tool.Path))

	if err != nil {
		fmt.Println(err.Error())
		download(tool)
	}

	for _, v := range tool.Install {
		cmd := exec.Command("go", "get", "-v", v)
		result, err := cmd.Output()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(result))
	}
}

func main() {
	tools, err := getConfig(confilePath)

	if err != nil {
		log.Fatal(err)
	}

	for _, v := range tools {
		install(v)
	}
}
