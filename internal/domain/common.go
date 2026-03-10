package domain

type Pagination struct {
	Limit  *uint
	Offset *uint
}

type Paginate struct {
	Total int64 `json:"total"`
	Count int64 `json:"count"`
}
