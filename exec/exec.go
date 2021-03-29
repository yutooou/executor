// created by yutooou
package exec

import (
	"encoding/json"
	"executor/compiler"
	"fmt"
	"io"
	"os"
	"path"
)

type Runner struct {
	JudgeId			string	`json:"judge_id"`			// 本次判题唯一标识
	ProblemDir		string	`json:"problem_dir"`		// 题目配置目录
	CodeLanguage	string	`json:"code_language"`		// 源码语言
	SourceCode		string	`json:"source_code"`		// 源代码
	workDir			string								// 当前测评工作目录
	compiler		compiler.Compiler					// 编译组件
	runCommands		[]string							// 可执行文件执行命令
	judgeConfig		JudgeConfiguration					// 当前题目配置文件
}

// 创建运行对象
func NewRunner(bytes []byte) (*Runner, error) {
	var s Runner
	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Runner) Judge() Result {
	result := Result{}
	result.JudgeId = r.JudgeId
	// 创建工作目录
	err := r.createWorkDir(&result)
	if err != nil {
		return result
	}
	// 编译
	err = r.compileCode(&result)
	if err != nil {
		return result
	}
	// 运行
	// 加载题目配置
	err = r.loadProblemConfig(&result)
	if err != nil {
		return result
	}
	
	var resultNumbers []int // 判题结果切片 位置对应测试点
	// 遍历测试每一个用例
	for i := 0; i < len(r.judgeConfig.TestCases); i++ {
		oneResult := r.runCase(r.judgeConfig.TestCases[i])
		// fmt.Println(oneResult)
		isFault := r.isDisastrousFault(&result, oneResult)
		result.MemoryUsed = Max32(oneResult.MemoryUsed, result.MemoryUsed)
		result.TimeUsed = Max32(oneResult.TimeUsed, result.TimeUsed)
		resultNumbers = append(resultNumbers, oneResult.JudgeResult)

		if isFault {
			break
		}
		if oneResult.JudgeResult != RESULT_AC && oneResult.JudgeResult != RESULT_PE {
			break
		}
	}
	// 生成结果
	r.final(&result, resultNumbers)

	return result
}

// 创建工作空间目录
func (r *Runner) createWorkDir(result *Result) error {
	// 工作目录规则：系统指定工作根目录WORKDIR/id
	r.workDir = path.Join(WORKDIR, r.JudgeId)
	_, err := os.Stat(r.workDir)
	if os.IsNotExist(err) {
		// 创建工作空间文件夹
		err := os.Mkdir(r.workDir, os.ModePerm)
		// 文件夹创建出错
		if err != nil {
			result.JudgeResult = RESULT_SE
			result.SeInfo = err.Error()
			return err
		}
	}
	return nil
}

// 加载题目配置以及评测点信息
func (r *Runner) loadProblemConfig(result *Result) error {
	fpath := path.Join(r.ProblemDir, PROBLEM_CONFIG_FILENAME)
	file, err := os.Open(fpath)
	if err != nil {
		// 文件打开失败
		result.JudgeResult = RESULT_SE
		result.SeInfo = err.Error()
		return err
	}
	// 加载完毕关闭文件
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		// 文件加载失败
		result.JudgeResult = RESULT_SE
		result.SeInfo = err.Error()
		return err
	}
	// 内容解析
	err = json.Unmarshal(bytes, &r.judgeConfig)
	if err != nil {
		result.JudgeResult = RESULT_SE
		result.SeInfo = err.Error()
		return err
	}
	return nil
}

// 判定是否是灾难性结果 灾难性后果的代码片段不宜向下继续测评
func (r *Runner) isDisastrousFault(judgeResult *Result, caseResult *TestCaseResult) bool {
	if caseResult.JudgeResult == RESULT_SE {
		judgeResult.JudgeResult = RESULT_SE
		judgeResult.SeInfo = fmt.Sprintf("testcase %s caused a problem", caseResult.Id)
		return true
	}
	return false
}

// 生成当前源码结果数据
func (r *Runner) final(result *Result, resultNumbers []int) {
	acCount, peCount, waCount := 0, 0, 0
	for _, num := range resultNumbers {
		// 如果，不是AC、PE、WA
		if num != RESULT_WA && num != RESULT_PE && num != RESULT_AC {
			//直接应用结果
			result.JudgeResult = num
			return
		}
		if num == RESULT_WA {
			waCount++
		}
		if num == RESULT_PE {
			peCount++
		}
		if num == RESULT_AC {
			acCount++
		}
	}
	// 非AC、PE不会向下执行，所以判断长度是否相符
	//if len(resultNumbers) != len(r.judgeConfig.TestCases) {
	//	// 如果测试数据未全部跑完
	//	result.JudgeResult = RESULT_WA
	//} else {
	//	// 如果测试数据未全部跑了
	//
	//}
	if waCount > 0 {
		// 如果存在WA，报WA
		result.JudgeResult = RESULT_WA
	} else if peCount > 0 {
		result.JudgeResult = RESULT_PE
	} else {
		result.JudgeResult = RESULT_AC
	}
}