// created by yutooou
package unix

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"unsafe"
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

// 进程信息
type ProcessInfo struct {
	Pid    uintptr            
	Status syscall.WaitStatus 
	Rusage syscall.Rusage     
}

// fork调用
func ForkProc() (pid uintptr, err error) {
	r1, r2, errMsg := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	darwin := runtime.GOOS == "darwin"
	if errMsg != 0 {
		return 0, fmt.Errorf("system call: fork(); error: %s", errMsg)
	}
	if darwin {
		if r2 == 1 {
			pid = 0
		} else {
			pid = r1
		}
	} else {
		if r1 == 0 && r2 == 0 {
			pid = 0
		} else {
			pid = r1
		}
	}
	return pid, nil
}

// 重映射文件描述符
func RedirectFileDescriptor(to int, path string, flag int, perm uint32) (fd int, err error) {
	fd, errMsg := getFileDescriptor(path, flag, perm)
	if errMsg == nil {
		errMsg = syscall.Dup2(fd, to)
		if errMsg != nil {
			syscall.Close(fd)
			return -1, errMsg
		}
		return fd, nil
	} else {
		return -1, errMsg
	}
}
// 打开并获取文件的描述符
func getFileDescriptor(path string, flag int, perm uint32) (fd int, err error) {
	var filed = 0
	_, errMsg := os.Stat(path)
	if errMsg != nil {
		if os.IsNotExist(err) {
			return 0, errMsg
		}
	}
	filed, errMsg = syscall.Open(path, flag, perm)
	return filed, nil
}

type RLimit struct {
	Which int
	RLim  syscall.Rlimit
}
type ITimerVal struct {
	ItInterval TimeVal
	ItValue    TimeVal
}
type TimeVal struct {
	TvSec  uint64
	TvUsec uint64
}

func GetRLimitEntity(cur, max uint64) syscall.Rlimit {
	return syscall.Rlimit{Cur: cur, Max: max}
}


// 硬件计时器
func SetHardTimer(realTimeLimit int) error {
	var prealt ITimerVal
	prealt.ItInterval.TvSec = uint64(math.Floor(float64(realTimeLimit) / 1000.0))
	prealt.ItInterval.TvUsec = uint64(realTimeLimit % 1000 * 1000)
	prealt.ItValue.TvSec = prealt.ItInterval.TvSec
	prealt.ItValue.TvUsec = prealt.ItInterval.TvUsec
	_, _, err := syscall.RawSyscall(syscall.SYS_SETITIMER, 0, uintptr(unsafe.Pointer(&prealt)), 0)
	if err != 0 {
		return fmt.Errorf("system call setitimer() error: %s", err)
	}
	return nil
}