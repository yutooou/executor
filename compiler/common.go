package compiler

import (
	"context"
	"errors"
	"executor/unix"
	"os"
	"strings"
)

const (
	filePERM = 0644
)

// 编译命令集
var compileCommands = struct {
	C    string
	CPP  string
	Java string
}{
	C:   "gcc %s -o %s -std=c11",
	CPP: "g++ %s -o %s -std=c++11",
}

// 编译器对外暴露的所有接口且必须实现
type Compiler interface {
	Init(code, workDir string) error // 初始化
	Compile() error                  // 编译
	RunArgs() (args []string)        // 可执行文件执行命令
}

// 检查工作目录
func checkWorkDir(workDir string) error {
	// 获取文件属性
	_, err := os.Stat(workDir)
	if err != nil {
		// 文件不存在
		if os.IsNotExist(err) {
			err = errors.New("workDir not exists")
		}
	}
	return err
}

func shell(command string, ctx context.Context) error {
	cmdArgs := strings.Split(command, " ")
	if len(cmdArgs) <= 1 {
		return errors.New("error command")
	}
	res, err := unix.UnixShell(&unix.ShellOptions{
		Context: ctx,
		Name:    cmdArgs[0],
		Args:    cmdArgs[1:],
		Error:   nil,
		Output:  nil,
		Input:   nil,
	})
	if err != nil {
		return err
	}
	if !res.Success {
		return errors.New(res.Stderr)
	}
	return nil
}

// 创建编译程序 对外暴露
func NewCompiler(name string) (Compiler, error) {
	switch name {
	case "c", "C":
		return &c{}, nil
	case "cpp", "CPP", "C++", "c++":
		return &cpp{}, nil
	default:
		return nil, errors.New("Language not supported")
	}
}
