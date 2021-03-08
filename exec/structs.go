// created by yutooou
package exec

// 测评配置
type JudgeConfiguration struct {
	TestCases     []TestCase                    `json:"test_cases"`      // 测试用例
	TimeLimit     int                           `json:"time_limit"`      // 实现限制 单位ms
	MemoryLimit   int                           `json:"memory_limit"`    // 内存限制 单位kb
	RealTimeLimit int                           `json:"real_time_limit"` // Real Time Limit (ms) (optional)
	FileSizeLimit int                           `json:"file_size_limit"` // File Size Limit (bytes) (optional)
}

// 测试数据
type TestCase struct {
	Input           string `json:"input"`             // Testcase input file path
	Output          string `json:"output"`            // Testcase output file path
}