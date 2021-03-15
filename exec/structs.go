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
	Id				int	`json:"id"`
	Input           string 	`json:"input"`           // 测试用例输入文件
	Output          string 	`json:"output"`          // 测试用例输出文件
}

// 评测结果
type Result struct {
	JudgeId			string					`json:"judge_id"`		// 本次判题唯一标识
	JudgeResult 	int                   	`json:"judge_result"` 	// 运行结果
	TimeUsed    	int                   	`json:"time_used"`    	// 使用时间
	MemoryUsed  	int                   	`json:"memory_used"`  	// 最大内存占用
	ReInfo      	string                	`json:"re_info"`      	// Runtime Error 提示信息
	SeInfo      	string                	`json:"se_info"`      	// System Error 提示信息
	CeInfo      	string                	`json:"ce_info"`      	// Compile Error 提示信息
}