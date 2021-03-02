// created by yutooou
package compiler

import (
	"context"
	"executor/utils"
	"fmt"
	"os"
	"path"
	"time"
)

type c struct {
	code            string // 源代码
	isReady         bool   // 是否编译完成
	codeFilePath    string // 目标源代码目录
	codeFileName    string // 目标源代码文件
	programFilePath string // 目标程序目录
	programFileName string // 目标程序文件
	workDir         string // 工作目录
}

// 完成编译程序的初始化工作
func (c *c) Init(code, workDir string) error {
	c.code = code
	// 检查工作目录
	err := checkWorkDir(workDir)
	c.workDir = workDir
	if err != nil {
		return err
	}

	// 写入文件
	err = c.createFile(".c", ".do")
	return err
}

func (c *c) createFile(codeFileSuffix, programFileSuffix string) error {
	// 写入源代码文件信息
	randomName := utils.UUID(12)
	c.codeFileName = fmt.Sprintf("%s%s", randomName, codeFileSuffix)
	c.codeFilePath = path.Join(c.workDir, c.codeFileName)
	// 写入可执行文件信息
	c.programFileName = fmt.Sprintf("%s%s", randomName, programFileSuffix)
	c.programFilePath = path.Join(c.workDir, c.programFileName)

	// 保存文件
	file, err := os.OpenFile(c.codeFilePath, os.O_RDWR|os.O_CREATE, filePERM)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(c.code)
	return err
}

func (c *c) Compile() error {
	// 写入编译命令
	cmd := fmt.Sprintf(compileCommands.C, c.codeFilePath, c.programFilePath)
	// 设置环境信息 编译指令最多执行10秒就强制退出
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := shell(cmd, ctx)
	if err == nil {
		c.isReady = true
	}
	return err
}

func (c *c) RunArgs() (args []string) {
	return []string{c.programFilePath}
}
