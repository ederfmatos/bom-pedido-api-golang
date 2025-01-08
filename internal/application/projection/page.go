package projection

type (
	Page[Item any] struct {
		CurrentPage int64  `json:"currentPage"`
		PageSize    int64  `json:"pageSize"`
		TotalPages  int64  `json:"totalPages"`
		TotalItems  int64  `json:"totalItems"`
		LastPage    bool   `json:"lastPage"`
		Items       []Item `json:"items"`
	}
)
