package exec

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"syscall"
)

// 进行文本比较
// Compare the text
func (r *Runner) diffText(result *TestCaseResult) error {
	answerInfo, err := os.Stat(path.Join(r.ProblemDir, result.Output))
	if err != nil {
		result.JudgeResult = RESULT_SE
		result.TextDiffLog = fmt.Sprintf("Get answer file info failed: %s", err.Error())
		return err
	}
	useroutInfo, err := os.Stat(path.Join(r.workDir, result.ProgramOut))
	if err != nil {
		result.JudgeResult = RESULT_SE
		result.TextDiffLog = fmt.Sprintf("Get userout file info failed: %s", err.Error())
		return err
	}

	useroutLen := useroutInfo.Size()
	answerLen := answerInfo.Size()

	sizeText := fmt.Sprintf("tcLen=%d, ansLen=%d", answerLen, useroutLen)

	var useroutBuffer, answerBuffer []byte
	errText := ""

	answerBuffer, errText, err = readFileWithTry(path.Join(r.ProblemDir, result.Output), "answer", 3)
	if err != nil {
		result.JudgeResult = RESULT_SE
		result.TextDiffLog = errText
		return err
	}

	useroutBuffer, errText, err = readFileWithTry(path.Join(r.workDir, result.ProgramOut), "userout", 3)
	if err != nil {
		result.JudgeResult = RESULT_SE
		result.TextDiffLog = errText
		return err
	}

	if useroutLen == 0 && answerLen == 0 {
		// Empty File AC
		result.JudgeResult = RESULT_AC
		result.TextDiffLog = sizeText + "; Accepted with zero size."
		return nil
	} else if useroutLen > 0 && answerLen > 0 {
		if (useroutLen > int64(r.judgeConfig.FileSizeLimit)) || (useroutLen >= answerLen * 2) {
			// OLE
			result.JudgeResult = RESULT_OLE
			if useroutLen > int64(r.judgeConfig.FileSizeLimit) {
				result.TextDiffLog = sizeText + "; WA: larger then limitation."
				return nil
			} else {
				result.TextDiffLog = sizeText + "; WA: larger then 2 times."
				return nil
			}
		}
	} else {
		result.JudgeResult = RESULT_WA
		result.TextDiffLog = sizeText + "; WA: less then zero size."
		return nil
	}

	rel, logText := charDiffIoUtil(useroutBuffer, answerBuffer, useroutLen, answerLen)
	result.JudgeResult = rel

	if rel != RESULT_WA {
		// PE or AC or SE
		if rel == RESULT_AC {
			// AC 时执行强制检查，可以排除空白字符的顺序不一致也是AC的情况
			sret := strictDiff(useroutBuffer, answerBuffer, useroutLen, answerLen)
			if !sret {
				result.JudgeResult = RESULT_PE
				logText = "Strict check: Presentation Error."
			} else {
				logText = "Accepted."
			}
		}
	} else {
		// WA
		sameLines, totalLines := lineDiff(r, result)
		result.SameLines = sameLines
		result.TotalLines = totalLines
	}
	result.TextDiffLog = sizeText + "; " + logText
	return nil
}


// 文件读写(有重试次数，checker专用)
func readFileWithTry(filePath string, name string, tryOnFailed int) ([]byte, string, error) {
	errCnt, errText := 0, ""
	var err error
	for errCnt < tryOnFailed {
		fp, err := OpenFile(filePath, os.O_RDONLY|syscall.O_NONBLOCK, 0)
		if err != nil {
			errText = err.Error()
			errCnt++
			continue
		}
		data, err := io.ReadAll(fp)
		if err != nil {
			_ = fp.Close()
			errText = fmt.Sprintf("Read file(%s) i/o error: %s", name, err.Error())
			errCnt++
			continue
		}
		_ = fp.Close()
		return data, errText, nil
	}
	return nil, errText, err
}

// 打开文件并获取描述符 (强制文件检查)
func OpenFile(filePath string, flag int, perm os.FileMode) (*os.File, error) {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file (%s) not exists", filePath)
		} else {
			return nil, fmt.Errorf("open file (%s) error: %s", filePath, err.Error())
		}
	} else {
		if fp, err := os.OpenFile(filePath, flag, perm); err != nil {
			return nil, fmt.Errorf("open file (%s) error: %s", filePath, err.Error())
		} else {
			return fp, nil
		}
	}
}

// 比较每一个字符，但是忽略空白
// Compare each char in buffer, but ignore the 'SpaceChar'
func charDiffIoUtil(useroutBuffer, answerBuffer []byte, useroutLen, answerLen int64) (rel int, logtext string) {
	var (
		leftPos, rightPos   int64 = 0, 0
		maxLength                 = Max(useroutLen, answerLen)
		leftByte, rightByte byte
	)
	for (leftPos < maxLength) && (rightPos < maxLength) && (leftPos < useroutLen) && (rightPos < answerLen) {
		if leftPos < useroutLen {
			leftByte = useroutBuffer[leftPos]
		}
		if rightPos < answerLen {
			rightByte = answerBuffer[rightPos]
		}

		for leftPos < useroutLen && isSpaceChar(leftByte) {
			leftPos++
			if leftPos < useroutLen {
				leftByte = useroutBuffer[leftPos]
			} else {
				leftByte = 0
			}
		}
		for rightPos < answerLen && isSpaceChar(rightByte) {
			rightPos++
			if rightPos < answerLen {
				rightByte = answerBuffer[rightPos]
			} else {
				rightByte = 0
			}
		}

		if leftByte != rightByte {
			return RESULT_WA, fmt.Sprintf(
				"WA: at leftPos=%d, rightPos=%d, leftByte=%d, rightByte=%d",
				leftPos,
				rightPos,
				leftByte,
				rightByte,
			)
		}
		leftPos++
		rightPos++
	}

	// 如果左游标没跑完
	for leftPos < useroutLen {
		leftByte = useroutBuffer[leftPos]
		if !isSpaceChar(leftByte) {
			return RESULT_WA, fmt.Sprintf(
				"WA: leftPos=%d, rightPos=%d, leftLen=%d, rightLen=%d",
				leftPos,
				rightPos,
				useroutLen,
				answerLen,
			)
		}
		leftPos++
	}
	// 如果右游标没跑完
	for rightPos < answerLen {
		rightByte = answerBuffer[rightPos]
		if !isSpaceChar(rightByte) {
			return RESULT_WA, fmt.Sprintf(
				"WA: leftPos=%d, rightPos=%d, leftLen=%d, rightLen=%d",
				leftPos,
				rightPos,
				useroutLen,
				answerLen,
			)
		}
		rightPos++
	}
	// 左右匹配，说明AC
	// if left cursor's position equals right cursor's, means Accepted.
	if leftPos == rightPos {
		return RESULT_AC, "AC!"
	} else {
		return RESULT_PE, fmt.Sprintf(
			"PE: leftPos=%d, rightPos=%d, leftLen=%d, rightLen=%d",
			leftPos,
			rightPos,
			useroutLen,
			answerLen,
		)
	}
}

func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
func Max32(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// 判断是否是空白字符
func isSpaceChar(ch byte) bool {
	return ch == '\n' || ch == '\r' || ch == ' ' || ch == '\t'
}
// 严格比较每一个字符
// Compare each char in buffer strictly
func strictDiff(useroutBuffer, answerBuffer []byte, useroutLen, answerLen int64) bool {
	if useroutLen != answerLen {
		return false
	}
	pos := int64(0)
	for ; pos < useroutLen; pos++ {
		leftByte, rightByte := useroutBuffer[pos], answerBuffer[pos]
		if leftByte != rightByte {
			return false
		}
	}
	return true
}

// 逐行比较，获取错误行数
// Compare each line, to find out the number of wrong line
func lineDiff(r *Runner, rst *TestCaseResult) (sameLines int, totalLines int) {
	answer, err := os.OpenFile(path.Join(r.ProblemDir, rst.Output), os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return 0, 0
	}
	defer answer.Close()
	userout, err := os.Open(path.Join(r.workDir, rst.ProgramOut))
	if err != nil {
		return 0, 0
	}
	defer userout.Close()

	useroutBuffer := bufio.NewReader(userout)
	answerBuffer := bufio.NewReader(answer)

	var (
		leftStr, rightStr       = "", ""
		leftErr, rightErr error = nil, nil
		leftCnt, rightCnt       = 0, 0
	)

	for leftErr == nil {
		leftStr, leftErr = readLineAndRemoveSpaceChar(answerBuffer)
		if leftStr == "" {
			continue
		}

		leftCnt++

		for rightStr == "" && rightErr == nil {
			rightStr, rightErr = readLineAndRemoveSpaceChar(useroutBuffer)
		}

		if rightStr == leftStr {
			rightCnt++
		}
		rightStr = ""
	}

	return rightCnt, leftCnt
}

// 使用IOReader读写每一行， 并移除空白字符
func readLineAndRemoveSpaceChar(buf *bufio.Reader) (string, error) {
	line, isContinue, err := buf.ReadLine()
	for isContinue && err == nil {
		var next []byte
		next, isContinue, err = buf.ReadLine()
		line = append(line, next...)
	}
	return removeSpaceChar(string(line)), err
}

// 移除空白字符
// Remove the special blank words in a string
func removeSpaceChar(source string) string {
	source = strings.Replace(source, "\t", "", -1)
	source = strings.Replace(source, "\r", "", -1)
	source = strings.Replace(source, "\n", "", -1)
	source = strings.Replace(source, " ", "", -1)
	return source
}