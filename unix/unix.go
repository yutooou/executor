// created by yutooou
package unix

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"syscall"
)

type ShellOptions struct {
	Context context.Context // 执行环境
	Name    string          // 执行程序名
	Args    []string        // 执行参数
	Input   io.Reader       // 输入流
	Output  io.Writer       // 输出流
	Error   io.Writer       // 错误流
}

type ShellResult struct {
	Success    bool
	ExitCode   int
	Stdout     string
	Stderr     string
	Signal     int
	ErrMessage string
}

func UnixShell(options *ShellOptions) (result *ShellResult, err error) {
	// 在环境变量PATH指定的目录中搜索可执行文件
	realPath, err := exec.LookPath(options.Name)
	if err != nil {
		return nil, err
	}

	// 渲染带上下文环境的Cmd
	cmd := exec.CommandContext(options.Context, realPath, options.Args...)

	// 输出缓存 若未设置error或者stdout则先写入缓存等待处理
	var stderr, stdout bytes.Buffer

	// 输出重定向
	if options.Output != nil {
		cmd.Stdout = options.Output
	} else {
		cmd.Stdout = &stdout
	}

	// 错误重定向
	if options.Error != nil {
		cmd.Stderr = options.Error
	} else {
		cmd.Stderr = &stderr
	}

	// 输入重定向
	if options.Input != nil {
		cmd.Stdin = options.Input
	} else {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return nil, err
		}
		_ = stdin.Close()
	}
	// 开始执行 并等待执行完成
	err = cmd.Run()

	// 处理结果
	result = &ShellResult{}
	if options.Output == nil {
		// 如果未设置Output则将缓冲中的输出放入result
		result.Stdout = stdout.String()
	}
	if options.Error == nil {
		// 如果未设置Error则将缓冲中的输出放入result
		result.Stderr = stderr.String()
	}

	result.ExitCode = cmd.ProcessState.ExitCode()

	result.Signal = int(cmd.ProcessState.Sys().(syscall.WaitStatus).Signal())

	// cmd执行错误
	if err != nil {
		result.Success = false
		result.ErrMessage = err.Error()
		if stdErr := result.Stderr; stdErr == "" {
			result.Stderr += err.Error()
		}
		// command已经执行，其err不应作为函数err返回
		return result, nil
	}
	result.Success = true
	return result, nil
}
