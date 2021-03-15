// created by yutooou
package exec

import (
	"encoding/json"
	"executor/compiler"
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
		r.runCase(r.judgeConfig.TestCases[i])
	}

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