// created by yutooou
package exec

import "executor/compiler"

type Runner struct {
	JudgeId			string	`json:"judge_id"`
	ProblemId		string	`json:"problem_id"`
	ProblemDir		string	`json:"problem_dir"`
	CodeLanguage	string	`json:"code_language"`
	SourceCode		string	`json:"source_code"`
	workDir			string
	runCommands		[]string
	compiler		compiler.Compiler
	judgeConfig		JudgeConfiguration
}

