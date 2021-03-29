package exec

import (
	"executor/unix"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"syscall"
)

// 测试单个测试点
func (r *Runner) runCase(testCase TestCase) *TestCaseResult {
	// 为单个测试点记录结果
	tcRes := TestCaseResult{}
	tcRes.Id = testCase.Id
	tcRes.Input = testCase.Input
	tcRes.Output = testCase.Output
	// 定义输出路径 judgeID_caseID.fileType
	tcRes.ProgramOut = r.JudgeId + "_" + fmt.Sprintf("%d", testCase.Id) + ".out"
	tcRes.ProgramError = r.JudgeId + "_" + fmt.Sprintf("%d", testCase.Id) + ".err"
	// 开始起程序
	pInfo, err := r.runProgram(&tcRes)
	// 程序开启期间任何问题都是SE
	if err != nil {
		tcRes.JudgeResult = RESULT_SE
		tcRes.SeInfo = err.Error()
		return &tcRes
	}
	r.saveExitRusage(&tcRes, pInfo)
	// 分析目标程序的状态
	r.analysisExitStatus(&tcRes, pInfo)
	// 只有什么状态都未曾写入的时候才进行文本比较！
	if tcRes.JudgeResult == RESULT_AC {
		// 进行文本比较
		err = r.diffText(&tcRes)
		if err != nil {
			tcRes.JudgeResult = RESULT_SE
			tcRes.SeInfo = err.Error()
			return &tcRes
		}
	}
	return &tcRes
}


// 运行程序
func (r *Runner) runProgram(rst *TestCaseResult) (*unix.ProcessInfo, error) {
	// 创建进程信息
	pinfo := unix.ProcessInfo{}
	// 创建子进程
	pid, fds, err := runProgramProcess(r, rst)

	if err != nil {
		if pid <= 0 {
			// 如果是子进程错误了，输出到程序的error去
			panic(err)
		}
		return nil, err
	}
	pinfo.Pid = pid

	// 获得子进程状态信息以及资源使用信息
	_, err = syscall.Wait4(int(pid), &pinfo.Status, syscall.WUNTRACED, &pinfo.Rusage)
	if err != nil {
		return nil, err
	}

	for _, fd := range fds {
		if fd > 0 {
			_ = syscall.Close(fd)
		}
	}
	return &pinfo, err
}

// 运行目标程序子进程
func runProgramProcess(r *Runner, rst *TestCaseResult) (uintptr, []int, error) {
	var (
		err error
		pid uintptr
		fds []int
	)

	fds = make([]int, 3)

	// 创建子进程，返回pid
	pid, err = unix.ForkProc()
	if err != nil {
		// 进程创建失败
		return 0, fds, fmt.Errorf("fork process error: %s", err.Error())
	}
	if pid == 0 {
		// Redirect test-case input to STDIN
		fds[0], err = unix.RedirectFileDescriptor(
			syscall.Stdin,
			path.Join(r.ProblemDir, rst.Input),
			os.O_RDONLY,
			0,
		)
		if err != nil {
			return 0, fds, err
		}

		// Redirect userOut to STDOUT
		fds[1], err = unix.RedirectFileDescriptor(
			syscall.Stdout,
			path.Join(r.workDir, rst.ProgramOut),
			os.O_WRONLY|os.O_CREATE,
			0644,
		)
		if err != nil {
			return 0, fds, err
		}

		// Redirect programError to STDERR
		fds[2], err = unix.RedirectFileDescriptor(
			syscall.Stderr,
			path.Join(r.workDir, rst.ProgramError),
			os.O_WRONLY|os.O_CREATE,
			0644,
		)
		if err != nil {
			return 0, fds, err
		}

		// Set Resource Limit
		tl, ml, rtl, fsl := getLimitation(r)
		err = setLimit(tl, ml, rtl, fsl)
		if err != nil {
			return 0, fds, err
		}

		// Run Program
		commands := r.runCommands
		// 参考exec.Command，从环境变量获取编译器/VM真实的地址
		programPath := commands[0]
		if filepath.Base(programPath) == programPath {
			if programPath, err = exec.LookPath(programPath); err != nil {
				return 0, fds, err
			}
		}
		if len(commands) > 1 {
			err = syscall.Exec(programPath, commands[1:], []string{})
		} else {
			err = syscall.Exec(programPath, nil, []string{})
		}
		//it won't be run.
	} else if pid < 0 {
		return 0, fds, fmt.Errorf("fork process error: pid < 0")
	}
	// parent process
	return pid, fds, err
}

// 获取资源限制的参数
func getLimitation(r *Runner) (int, int, int, int) {
	memoryLimitExtend := 0
	return r.judgeConfig.TimeLimit,
		r.judgeConfig.MemoryLimit + memoryLimitExtend,
		r.judgeConfig.RealTimeLimit,
		r.judgeConfig.FileSizeLimit
}

func setLimit(timeLimit, memoryLimit, realTimeLimit, fileSizeLimit int) error {

	// Set stack limit
	stack := uint64(memoryLimit * 1024)
	if runtime.GOOS == "darwin" { 
		stack = uint64(65500 * 1024)
	}

	rlimits := []unix.RLimit{
		// Set time limit: RLIMIT_CPU
		{
			Which: syscall.RLIMIT_CPU,
			RLim: unix.GetRLimitEntity(
				uint64(math.Ceil(float64(timeLimit)/1000.0)),
				uint64(math.Ceil(float64(timeLimit)/1000.0)),
			),
		},
		// Set memory limit: RLIMIT_DATA
		{
			Which: syscall.RLIMIT_DATA,
			RLim: unix.GetRLimitEntity(
				uint64(memoryLimit*1024),
				uint64(memoryLimit*1024),
			),
		},
		// Set memory limit: RLIMIT_AS
		{
			Which: syscall.RLIMIT_AS,
			RLim: unix.GetRLimitEntity(
				uint64(memoryLimit*1024*2),
				uint64(memoryLimit*1024*2+1024),
			),
		},
		// Set stack limit
		{
			Which: syscall.RLIMIT_STACK,
			RLim: unix.GetRLimitEntity(
				stack,
				stack,
			),
		},
		// Set file size limit: RLIMIT_FSIZE
		{
			Which: syscall.RLIMIT_FSIZE,
			RLim: unix.GetRLimitEntity(
				uint64(fileSizeLimit),
				uint64(fileSizeLimit),
			),
		},
	}

	for _, rlimit := range rlimits {
		err := syscall.Setrlimit(rlimit.Which, &rlimit.RLim)
		if err != nil {
			return fmt.Errorf("setrlimit(%d) error: %s", rlimit.Which, err)
		}
	}

	// Set time limit: setITimer
	if realTimeLimit > 0 {
		err := unix.SetHardTimer(realTimeLimit)
		if err != nil {
			return err
		}
	}

	return nil
}

// 分析进程资源占用
func (r *Runner) saveExitRusage(rst *TestCaseResult, pinfo *unix.ProcessInfo) {
	ru := pinfo.Rusage
	status := pinfo.Status

	tu := int(ru.Utime.Sec*1000 + int64(ru.Utime.Usec)/1000 + ru.Stime.Sec*1000 + int64(ru.Stime.Usec)/1000)
	mu := int(ru.Minflt * int64(syscall.Getpagesize()/1024))

	// 特判
	rst.TimeUsed = tu
	rst.MemoryUsed = mu
	rst.ReSignum = int(status.Signal())
}


// 分析进程退出状态
func (r *Runner) analysisExitStatus(rst *TestCaseResult, pinfo *unix.ProcessInfo) {
	status := pinfo.Status

	// If process stopped with a signal
	if status.Signaled() {
		sig := status.Signal()
		if sig == syscall.SIGSEGV {
			// MLE or RE can also get SIGSEGV signal.
			if rst.MemoryUsed > r.judgeConfig.MemoryLimit {
				rst.JudgeResult = RESULT_MLE
			} else {
				rst.JudgeResult = RESULT_RE
				if r, e := SignalNumberMap[rst.ReSignum]; e {
					rst.ReInfo = fmt.Sprintf("%s: %s", r[0], r[1])
				}
			}
		} else if sig == syscall.SIGXFSZ {
			// SIGXFSZ signal means OLE
			rst.JudgeResult = RESULT_OLE
		} else if sig == syscall.SIGALRM || sig == syscall.SIGVTALRM || sig == syscall.SIGXCPU {
			// Normal TLE signal
			rst.JudgeResult = RESULT_TLE
		} else if sig == syscall.SIGKILL {
			// Sometimes MLE might get SIGKILL signal.
			// So if real time used lower than TIME_LIMIT - 100, it might be a TLE error.
			if rst.TimeUsed > (r.judgeConfig.TimeLimit - 100) {
				rst.JudgeResult = RESULT_TLE
			} else {
				rst.JudgeResult = RESULT_MLE
			}
		} else {
			// Otherwise, called runtime error.
			rst.JudgeResult = RESULT_RE
			if r, e := SignalNumberMap[rst.ReSignum]; e {
				rst.ReInfo = fmt.Sprintf("%s: %s", r[0], r[1])
			}
		}
	} else {
		// Sometimes setrlimit doesn't work accurately.
		if rst.TimeUsed > r.judgeConfig.TimeLimit {
			rst.JudgeResult = RESULT_MLE
		} else if rst.MemoryUsed > r.judgeConfig.MemoryLimit {
			rst.JudgeResult = RESULT_MLE
		} else {
			rst.JudgeResult = RESULT_AC
		}
	}
}