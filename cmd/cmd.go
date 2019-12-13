package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"os"
)

type FabCmd struct {
	args     []string
	dir      string
	fileName string
}

func NewFabCmd(fileName, host, sshuser, sshpwd, sshport, sshkey string) *FabCmd {
	var args []string
	var checkAdd = func(f, val string) {
		if val != "" {
			args = append(args, f, val)
		}
	}
	if sshuser == "" {
		sshuser = GlobalConfig.SshUserName
	}
	if sshpwd == "" {
		sshpwd = GlobalConfig.SshPwd
	}
	if sshport == "" {
		sshport = GlobalConfig.SshPort
	}
	if sshport == "" {
		sshport = "22"
	}
	if sshkey == "" {
		sshkey = GlobalConfig.SshKey
	}
	checkAdd("-H", host)
	checkAdd("--port", sshport)
	checkAdd("-u", sshuser)
	checkAdd("-p", sshpwd)
	if PathExists(sshkey) {
		checkAdd("-i", sshkey)
	}
	return &FabCmd{args: append(args, "-f"), dir: ScriptPath(), fileName: fileName}
}

//判断文件或文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func NewLocalFabCmd(fileName string) *FabCmd {
	return &FabCmd{args: []string{"-f"}, dir: ScriptPath(), fileName: fileName}
}

func (c *FabCmd) SetDir(dir string) {
	c.dir = dir
}

func (c *FabCmd) SetFileName(name string) {
	c.fileName = name
}

func (c *FabCmd) Run(function string, args ...string) ([]byte, error) {
	return c.FileRun(c.fileName, function, args...)
}

//func (c *FabCmd) RunGet(function string, args ...string) (string, error) {
//	ret, err := c.FileRun(c.fileName, function, args...)
//	if err != nil {
//		return "", err
//	}
//	rets, err := getTagFields(string(ret))
//	if err != nil {
//		return "", err
//	}
//	return rets[len(rets)-1], nil
//}

func (c *FabCmd) RunShow(function string, args ...string) error {
	ret, err := c.FileRun(c.fileName, function, args...)
	if len(ret) != 0 {
		fmt.Println(string(ret))
	}
	return err
}

func (c *FabCmd) FileRun(fileName, function string, args ...string) ([]byte, error) {
	if fileName == "" {
		return nil, fmt.Errorf("fileName is empty")
	} else if function == "" {
		return nil, fmt.Errorf("function name is empty")
	}
	arg := function
	if len(args) != 0 {
		arg += ":" + strings.Join(args, ",")
	}
	_args := append(c.args, fileName, arg)
	fmt.Println("***************************************************************************")
	fmt.Println("fab:", _args)
	cmd := exec.Command("fab", _args...)
	cmd.Dir = c.dir
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if len(stderr.Bytes()) != 0 {
		if err != nil {
			fmt.Print(stderr.String())
			err = fmt.Errorf("exec python method failed")
		} else {
			fmt.Print(stderr.String())
		}
	}

	return stdout.Bytes(), err
}
