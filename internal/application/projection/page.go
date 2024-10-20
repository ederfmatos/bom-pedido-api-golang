package projection

type (
	Page[Item any] struct {
		CurrentPage int32  `json:"currentPage"`
		PageSize    int32  `json:"pageSize"`
		TotalPages  int32  `json:"totalPages"`
		TotalItems  int32  `json:"totalItems"`
		LastPage    bool   `json:"lastPage"`
		Items       []Item `json:"items"`
	}
)
