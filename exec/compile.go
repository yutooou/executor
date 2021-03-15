package exec

import (
	"errors"
	"executor/compiler"
)

func (r *Runner) compileCode(result *Result) error {
	// 获取对应语言编译程序
	err := r.getCompiler()
	if err != nil {
		// 编译器获取失败 可能为语言不支持
		result.JudgeResult = RESULT_SE
		result.SeInfo = err.Error()
		return err
	}

	// 编译
	err = r.compiler.Compile()
	if err != nil {
		// 编译失败 直接将err内容记录
		result.JudgeResult = RESULT_CE
		result.CeInfo = err.Error()
		err = errors.New("compile error")
		return err
	}

	// 编译后获取可执行文件执行指令
	r.runCommands = r.compiler.RunArgs()
	return nil
}

// 获取对应语言编译程序
func (r *Runner) getCompiler() (err error) {
	// var c compiler.Compiler
	c, err := compiler.NewCompiler(r.CodeLanguage)
	if err != nil {
		return err
	}
	r.compiler = c
	return nil
}
