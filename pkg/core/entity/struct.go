package entity

type Pagination struct {
	TotalItems   int  `json:"total_items"`
	TotalPages   int  `json:"total_pages"`
	CurrentPage  int  `json:"current_page"`
	NextPage     *int `json:"next_page,omitempty"`
	PreviousPage *int `json:"prev_page,omitempty"`
	FirstPage    int  `json:"first_page"`
	LastPage     int  `json:"last_page"`
}
