package domain

// ApiResponse 统一API响应格式
type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) ApiResponse {
	return ApiResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string) ApiResponse {
	return ApiResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

// NewSuccessResponseWithMessage 创建带自定义消息的成功响应
func NewSuccessResponseWithMessage(data interface{}, message string) ApiResponse {
	return ApiResponse{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// PaginatedData 分页数据通用结构
type PaginatedData struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Records  interface{} `json:"records,omitempty"` // 用于训练记录
	Plans    interface{} `json:"plans,omitempty"`   // 用于计划
	Templates interface{} `json:"templates,omitempty"` // 用于模板
}
