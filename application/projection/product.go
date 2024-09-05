package projection

type (
	ProductListFilter struct {
		CurrentPage int32
		PageSize    int32
		TenantId    string
	}

	ProductListItem struct {
		Id          string  `json:"id"`
		Name        string  `json:"name"`
		Description *string `json:"description,omitempty"`
		Price       string  `json:"price"`
		Status      string  `json:"status"`
		ImageURL    string  `json:"imageURL,omitempty"`
	}
)
