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
	PROBLEM_CONFIG_FILENAME = "config.json"
)

const (
	RESULT_AC  = iota	// 0 Accepted
	RESULT_PE 			// 1 Presentation Error
	RESULT_WA 			// 2 Wrong Answer
	RESULT_MLE			// 3 Memory Limit Exceeded
	RESULT_OLE			// 4 Output Limit Exceeded
	RESULT_TLE			// 5 Time Limit Exceeded
	RESULT_CE 			// 6 Compile Error
	RESULT_RE 			// 7 Runtime Error
	RESULT_SE 			// 8 System Error
)