package model

// ParamInfo 表示需要填充的参数信息
type ParamInfo struct {
	Param string `json:"param"` // 参数字段名称
	Type  string `json:"type"`  // 参数类型：string、url、number等
	Value string `json:"value"` // 参数值，用户输入后回填
}

// HiddenParams 用于存储隐藏参数
type HiddenParams struct {
	NeedFillParams []ParamInfo `json:"needFillParams"` // 需要填充的参数列表
	Raw            string      // 保存原始JSON字符串
	ShellCode      string      `json:"shellCode"` // 保存代码片段
}

// HasParameters 返回是否有需要填充的参数
func (h *HiddenParams) HasParameters() bool {
	return len(h.NeedFillParams) > 0
}
