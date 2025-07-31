// internal/application/dto/common.go
package dto

// ListRequest represents common pagination and filtering parameters
type ListRequest struct {
	Page     int    `json:"page" validate:"min=1" default:"1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100" default:"10"`
	SortBy   string `json:"sort_by" validate:"omitempty,oneof=created_at updated_at name balance"`
	SortDir  string `json:"sort_dir" validate:"omitempty,oneof=asc desc" default:"desc"`
	Search   string `json:"search" validate:"omitempty,max=100"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// SuccessResponse represents success response structure
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
