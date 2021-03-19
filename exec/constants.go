package exec

import (
	"fmt"
	"os"
)

func init() {
	_, err := os.Stat(WORKDIR)
	if os.IsNotExist(err) {
		// 创建工作空间文件夹
		err := os.Mkdir(WORKDIR, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("sys_workdir create failed!, err = %v", err.Error()))
		}
	}
}

const (
	WORKDIR = "/tmp/ytoj"	// 工作空间
	PROBLEM_CONFIG_FILENAME = "config.json"	// 测评题目配置文件
)

// 评测结果
const (
	RESULT_AC  = iota + 1	// 1 Accepted
	RESULT_PE 				// 2 Presentation Error
	RESULT_WA 				// 3 Wrong Answer
	RESULT_MLE				// 4 Memory Limit Exceeded
	RESULT_OLE				// 5 Output Limit Exceeded
	RESULT_TLE				// 6 Time Limit Exceeded
	RESULT_CE 				// 7 Compile Error
	RESULT_RE 				// 8 Runtime Error
	RESULT_SE 				// 9 System Error
)


var SignalNumberMap = map[int][]string{
	1: {"SIGHUP", "Hangup (POSIX)."},
	2: {"SIGINT", "Interrupt (ANSI)."},
	3: {"SIGQUIT", "Quit (POSIX)."},
	4: {"SIGILL", "Illegal instruction (ANSI)."},
	5: {"SIGTRAP", "Trace trap (POSIX)."},
	6: {"SIGABRT", "Abort (ANSI)."},
	7: {"SIGBUS", "BUS error (4.2 BSD)."},
	8: {"SIGFPE", "Floating-point exception (ANSI)."},
	9: {"SIGKILL", "Kill, unblockable (POSIX)."},
	10: {"SIGUSR1", "User-defined signal 1 (POSIX)."},
	11: {"SIGSEGV", "Segmentation violation (ANSI)."},
	12: {"SIGUSR2", "User-defined signal 2 (POSIX)."},
	13: {"SIGPIPE", "Broken pipe (POSIX)."},
	14: {"SIGALRM", "Alarm clock (POSIX)."},
	15: {"SIGTERM", "Termination (ANSI)."},
	16: {"SIGSTKFLT", "Stack fault."},
	17: {"SIGCHLD", "Child status has changed (POSIX)."},
	18: {"SIGCONT", "Continue (POSIX)."},
	19: {"SIGSTOP", "Stop, unblockable (POSIX)."},
	20: {"SIGTSTP", "Keyboard stop (POSIX)."},
	21: {"SIGTTIN", "Background read from tty (POSIX)."},
	22: {"SIGTTOU", "Background write to tty (POSIX)."},
	23: {"SIGURG", "Urgent condition on socket (4.2 BSD)."},
	24: {"SIGXCPU", "CPU limit exceeded (4.2 BSD)."},
	25: {"SIGXFSZ", "File size limit exceeded (4.2 BSD)."},
	26: {"SIGVTALRM", "Virtual alarm clock (4.2 BSD)."},
	27: {"SIGPROF", "Profiling alarm clock (4.2 BSD)."},
	28: {"SIGWINCH", "Window size change (4.3 BSD, Sun)."},
	29: {"SIGIO", "I/O now possible (4.2 BSD)."},
	30: {"SIGPWR", "Power failure restart (System V)."},
	31: {"SIGSYS", "Bad system call."},
}