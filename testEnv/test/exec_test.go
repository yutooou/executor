package test

import (
	"executor/exec"
	"fmt"
	"io"
	"os"
	"testing"
)

/*
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
*/
func TestRunnerAC(t *testing.T) {
	result := testRunner("/Users/yutooou/go/src/executor/testEnv/jsondata/ac.json")
	switch result.JudgeResult {
	case 1:
		fmt.Println("\x1b[0;42m!!!AC!!!\x1b[0m")
	default:
		fmt.Printf("\x1b[0;41m!!!NOT AC :%d!!!\x1b[0m\n", result.JudgeResult)
	}
}

func TestRunnerPE(t *testing.T) {
	result := testRunner("/Users/yutooou/go/src/executor/testEnv/jsondata/pe1.json")
	switch result.JudgeResult {
	case 2:
		fmt.Println("\x1b[0;42m!!!PE!!!\x1b[0m")
	default:
		fmt.Printf("\x1b[0;41m!!!NOT PE1 :%d!!!\x1b[0m\n", result.JudgeResult)
	}

	result = testRunner("/Users/yutooou/go/src/executor/testEnv/jsondata/pe2.json")
	switch result.JudgeResult {
	case 2:
		fmt.Println("\x1b[0;42m!!!PE!!!\x1b[0m")
	default:
		fmt.Printf("\x1b[0;41m!!!NOT PE2 :%d!!!\x1b[0m\n", result.JudgeResult)
	}

	result = testRunner("/Users/yutooou/go/src/executor/testEnv/jsondata/pe3.json")
	switch result.JudgeResult {
	case 2:
		fmt.Println("\x1b[0;42m!!!PE!!!\x1b[0m")
	default:
		fmt.Printf("\x1b[0;41m!!!NOT PE3 :%d!!!\x1b[0m\n", result.JudgeResult)
	}
}

func TestRunnerWA(t *testing.T) {
	result := testRunner("/Users/yutooou/go/src/executor/testEnv/jsondata/wa.json")
	switch result.JudgeResult {
	case 3:
		fmt.Println("\x1b[0;42m!!!WA!!!\x1b[0m")
	default:
		fmt.Printf("\x1b[0;41m!!!NOT WA :%d!!!\x1b[0m\n", result.JudgeResult)
	}

	result = testRunner("/Users/yutooou/go/src/executor/testEnv/jsondata/wa2.json")
	switch result.JudgeResult {
	case 3:
		fmt.Println("\x1b[0;42m!!!WA!!!\x1b[0m")
	default:
		fmt.Printf("\x1b[0;41m!!!NOT WA2 :%d!!!\x1b[0m\n", result.JudgeResult)
	}
}


func TestRunnerMLE(t *testing.T) {
	result := testRunner("/Users/yutooou/go/src/executor/testEnv/jsondata/mle.json")
	switch result.JudgeResult {
	case 4:
		fmt.Println("\x1b[0;42m!!!MLE!!!\x1b[0m")
	default:
		fmt.Printf("\x1b[0;41m!!!NOT MLE :%d!!!\x1b[0m\n", result.JudgeResult)
	}
}


func TestRunnerOLE(t *testing.T) {
	result := testRunner("/Users/yutooou/go/src/executor/testEnv/jsondata/ole1.json")
	switch result.JudgeResult {
	case 5:
		fmt.Println("\x1b[0;42m!!!OLE!!!\x1b[0m")
	default:
		fmt.Printf("\x1b[0;41m!!!NOT OLE1 :%d!!!\x1b[0m\n", result.JudgeResult)
	}

	result = testRunner("/Users/yutooou/go/src/executor/testEnv/jsondata/ole2.json")
	switch result.JudgeResult {
	case 5:
		fmt.Println("\x1b[0;42m!!!OLE!!!\x1b[0m")
	default:
		fmt.Printf("\x1b[0;41m!!!NOT OLE2 :%d!!!\x1b[0m\n", result.JudgeResult)
	}
}

func TestRunnerTLE(t *testing.T) {
	result := testRunner("/Users/yutooou/go/src/executor/testEnv/jsondata/tle1.json")
	fmt.Println(result.TimeUsed)
	switch result.JudgeResult {
	case 6:
		fmt.Println("\x1b[0;42m!!!TLE!!!\x1b[0m")
	default:
		fmt.Printf("\x1b[0;41m!!!NOT TLE1 :%d!!!\x1b[0m\n", result.JudgeResult)
	}

	// 未引入sleep
	result = testRunner("/Users/yutooou/go/src/executor/testEnv/jsondata/tle2.json")
	switch result.JudgeResult {
	case 6:
		fmt.Println("\x1b[0;42m!!!TLE!!!\x1b[0m")
	default:
		fmt.Printf("\x1b[0;41m!!!NOT TLE2 :%d!!!\x1b[0m\n", result.JudgeResult)
	}
}


func TestRunnerCE(t *testing.T) {
	result := testRunner("/Users/yutooou/go/src/executor/testEnv/jsondata/ce.json")
	switch result.JudgeResult {
	case 7:
		fmt.Println("\x1b[0;42m!!!CE!!!\x1b[0m")
	default:
		fmt.Printf("\x1b[0;41m!!!NOT CE :%d!!!\x1b[0m\n", result.JudgeResult)
	}
}

func TestRunnerRE(t *testing.T) {
	result := testRunner("/Users/yutooou/go/src/executor/testEnv/jsondata/re1.json")
	switch result.JudgeResult {
	case 8:
		fmt.Println("\x1b[0;42m!!!RE!!!\x1b[0m")
	default:
		fmt.Printf("\x1b[0;41m!!!NOT RE1 :%d!!!\x1b[0m\n", result.JudgeResult)
	}

	// 未引入sleep
	result = testRunner("/Users/yutooou/go/src/executor/testEnv/jsondata/re2.json")
	switch result.JudgeResult {
	case 8:
		fmt.Println("\x1b[0;42m!!!RE!!!\x1b[0m")
	default:
		fmt.Printf("\x1b[0;41m!!!NOT RE2 :%d!!!\x1b[0m\n", result.JudgeResult)
	}
}

func testRunner(filePath string) exec.Result {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	bytes, err := io.ReadAll(file)
	runner, err := exec.NewRunner(bytes)
	if err != nil {
		panic(err)
	}
	result := runner.Judge()
	return result
}