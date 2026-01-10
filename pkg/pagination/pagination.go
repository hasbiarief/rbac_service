package pagination

import "gin-scalable-api/internal/constants"

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

// Normalize ensures pagination parameters are within valid ranges
func (p *PaginationParams) Normalize() {
	if p.Limit <= 0 {
		p.Limit = constants.DefaultLimit
	}
	if p.Limit > constants.MaxLimit {
		p.Limit = constants.MaxLimit
	}
	if p.Offset < 0 {
		p.Offset = constants.DefaultOffset
	}
}

// GetLimitOffset returns normalized limit and offset
func (p *PaginationParams) GetLimitOffset() (int, int) {
	p.Normalize()
	return p.Limit, p.Offset
}

// PaginationResponse represents paginated response metadata
type PaginationResponse struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasMore bool  `json:"has_more"`
}

// NewPaginationResponse creates pagination response metadata
func NewPaginationResponse(total int64, limit, offset, currentCount int) *PaginationResponse {
	return &PaginationResponse{
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+currentCount) < total,
	}
}
