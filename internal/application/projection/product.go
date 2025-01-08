package projection

type (
	ProductListFilter struct {
		CurrentPage int64
		PageSize    int64
		TenantId    string
	}

	ProductListItem struct {
		Id          string  `json:"id"`
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
		Price       string  `json:"price"`
		Status      string  `json:"status"`
		ImageURL    string  `json:"imageURL,omitempty"`
		Category    struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"category"`
	}
)
