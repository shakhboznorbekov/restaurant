package menu_category

import "mime/multipart"

type Filter struct {
	Limit  *int
	Offset *int
}

type AdminGetListResponse struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
}

type AdminGetDetailResponse struct {
	ID        int64   `json:"id"`
	Name      *string `json:"name"`
	Logo      *string `json:"logo"`
	CreatedAt *string `json:"created_at"`
}

type AdminCreateRequest struct {
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"logo" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
}

type AdminCreateResponse struct {
	ID        int64   `json:"id"`
	Name      *string `json:"name"`
	Logo      *string `json:"-"`
	CreatedAt *string `json:"created_at"`
}

type AdminUpdateRequest struct {
	ID       int64                 `json:"id" form:"id"`
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"logo" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
}

type BranchGetListResponse struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
}

type BranchGetDetailResponse struct {
	ID        int64   `json:"id"`
	Name      *string `json:"name"`
	Logo      *string `json:"logo"`
	CreatedAt *string `json:"created_at"`
}

type BranchCreateRequest struct {
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"logo" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
}

type BranchCreateResponse struct {
	ID        int64   `json:"id"`
	Name      *string `json:"name"`
	Logo      *string `json:"-"`
	CreatedAt *string `json:"created_at"`
}

type BranchUpdateRequest struct {
	ID       int64                 `json:"id" form:"id"`
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"logo" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
}

type CashierGetListResponse struct {
	ID   int64   `json:"id"`
	Name *string `json:"name"`
	Logo *string `json:"logo"`
}

type CashierGetDetailResponse struct {
	ID        int64   `json:"id"`
	Name      *string `json:"name"`
	Logo      *string `json:"logo"`
	CreatedAt *string `json:"created_at"`
}

type CashierCreateRequest struct {
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"logo" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
}

type CashierCreateResponse struct {
	ID        int64   `json:"id"`
	Name      *string `json:"name"`
	Logo      *string `json:"-"`
	CreatedAt *string `json:"created_at"`
}

type CashierUpdateRequest struct {
	ID       int64                 `json:"id" form:"id"`
	Name     *string               `json:"name" form:"name"`
	Logo     *multipart.FileHeader `json:"logo" form:"logo"`
	LogoLink *string               `json:"-" form:"-"`
}
