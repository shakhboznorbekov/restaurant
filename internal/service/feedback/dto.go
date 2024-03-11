package feedback

type Filter struct {
	Limit  *int
	Offset *int
}

// @admin

type AdminGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}

type AdminGetDetail struct {
	ID   int64             `json:"id"`
	Name map[string]string `json:"name"`
}

type AdminCreate struct {
	Name map[string]string `json:"name" form:"name"`
}

type AdminUpdate struct {
	ID   int64              `json:"id" form:"id"`
	Name *map[string]string `json:"name" form:"name"`
}

// @client

type ClientGetList struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
}
