package dto

type PaginationResponse struct {
	Size        int `json:"size"`
	Total       int `json:"total"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
}
