package catalog

import (
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/service/banner"
	"github.com/restaurant/internal/service/basket"
	"github.com/restaurant/internal/service/feedback"
	"github.com/restaurant/internal/service/full_search"
	"github.com/restaurant/internal/service/service_percentage"
	"github.com/restaurant/internal/service/story"
	"net/http"
	"reflect"
)

type Controller struct {
	useCase *catalog.UseCase
}

func NewController(useCase *catalog.UseCase) *Controller {
	return &Controller{useCase}
}

// #measureUnit

// @admin

func (uc Controller) AdminGetMeasureUnitList(c *web.Context) error {
	var filter measure_unit_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetMeasureUnitList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

// @branch

func (uc Controller) BranchGetMeasureUnitList(c *web.Context) error {
	var filter measure_unit_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetMeasureUnitList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

// @super-admin

func (uc Controller) SuperAdminGetMeasureUnitList(c *web.Context) error {
	var filter measure_unit_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.SuperAdminGetMeasureUnitList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminGetMeasureUnitDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.SuperAdminGetMeasureUnitDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminCreateMeasureUnit(c *web.Context) error {
	var request measure_unit_service.SuperAdminCreateRequest

	if err := c.BindFunc(&request, "Name"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.SuperAdminCreateMeasureUnit(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateMeasureUnitAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request measure_unit_service.SuperAdminUpdateRequest

	if err := c.BindFunc(&request, "ID", "Name"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.SuperAdminUpdateMeasureUnit(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateMeasureUnitColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request measure_unit_service.SuperAdminUpdateRequest

	if err := c.BindFunc(&request, "ID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.SuperAdminUpdateMeasureUnitColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminDeleteMeasureUnit(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.SuperAdminDeleteMeasureUnit(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// #region

func (uc Controller) SuperAdminGetRegionList(c *web.Context) error {
	var filter region_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.SuperAdminGetRegionList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminGetRegionDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.SuperAdminGetRegionDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminCreateRegion(c *web.Context) error {
	var request region_service.SuperAdminCreateRequest

	if err := c.BindFunc(&request, "Name", "Code"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.SuperAdminCreateRegion(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateRegionAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request region_service.SuperAdminUpdateRequest

	if err := c.BindFunc(&request, "ID", "Name", "Code"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.SuperAdminUpdateRegion(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateRegionColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request region_service.SuperAdminUpdateRequest

	if err := c.BindFunc(&request, "ID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.SuperAdminUpdateRegionColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminDeleteRegion(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.SuperAdminDeleteRegion(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// #district

func (uc Controller) SuperAdminGetDistrictList(c *web.Context) error {
	var filter district_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if name, ok := c.GetQueryFunc(reflect.String, "name").(*string); ok {
		filter.Name = name
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.SuperAdminGetDistrictList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminGetDistrictDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.SuperAdminGetDistrictDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminCreateDistrict(c *web.Context) error {
	var request district_service.SuperAdminCreateRequest

	if err := c.BindFunc(&request, "Name", "Code", "RegionID"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.SuperAdminCreateDistrict(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateDistrictAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request district_service.SuperAdminUpdateRequest

	if err := c.BindFunc(&request, "ID", "Name", "Code", "RegionID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.SuperAdminUpdateDistrict(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateDistrictColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request district_service.SuperAdminUpdateRequest

	if err := c.BindFunc(&request, "ID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.SuperAdminUpdateDistrictColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminDeleteDistrict(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.SuperAdminDeleteDistrict(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// #story

// @admin

func (uc Controller) AdminGetStoryList(c *web.Context) error {
	var filter story.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if expired, ok := c.GetQueryFunc(reflect.Bool, "expired").(*bool); ok {
		filter.Expired = expired
	}
	if status, ok := c.GetQueryFunc(reflect.String, "status").(*string); ok {
		filter.Status = status
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetStoryList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateStory(c *web.Context) error {
	var request story.AdminCreateRequest

	if err := c.BindFunc(&request, "Name", "File", "Duration", "Type"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateStory(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateStatusStory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminUpdateStatusStory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteStory(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteStory(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @client

func (uc Controller) ClientGetStoryList(c *web.Context) error {
	var filter story.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	list, count, err := uc.useCase.ClientGetStoryList(c.Ctx, filter)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) ClientSetStoryAsViewed(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.ClientSetViewed(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "viewed!",
		"status": true,
	}, http.StatusOK)
}

// @client

func (uc Controller) SuperAdminGetStoryList(c *web.Context) error {
	var filter story.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if status, ok := c.GetQueryFunc(reflect.String, "status").(*string); ok {
		filter.Status = status
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	list, count, err := uc.useCase.SuperAdminGetStoryList(c.Ctx, filter)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateStoryStatus(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	body := struct {
		Status string `json:"status"`
	}{}

	if err := c.BindFunc(&body, "Status"); err != nil {
		return c.RespondError(err)
	}

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.SuperAdminUpdateStoryStatus(c.Ctx, int64(id), body.Status)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// #banner

// @branch

func (uc Controller) BranchGetBannerList(c *web.Context) error {
	var filter banner.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if expired, ok := c.GetQueryFunc(reflect.Bool, "expired").(*bool); ok {
		filter.Expired = expired
	}
	if status, ok := c.GetQueryFunc(reflect.String, "status").(*string); ok {
		filter.Status = status
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetBannerList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchGetBannerByID(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	c.Set("lang", c.GetHeader("Accept-Language"))

	res, err := uc.useCase.BranchGetBannerByID(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   res,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateBanner(c *web.Context) error {
	var request banner.BranchCreateRequest

	if err := c.BindFunc(&request, "Title", "Description", "Price", "MenuIds", "Photo"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateBanner(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateAllBanner(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}
	var request banner.BranchUpdateRequest

	if err := c.BindFunc(&request, "Title", "Description", "Price", "Photo"); err != nil {
		return c.RespondError(err)
	}
	request.ID = int64(id)

	err := uc.useCase.BranchUpdateAll(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateColumnBanner(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}
	var request banner.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}
	request.ID = int64(id)

	err := uc.useCase.BranchUpdateColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateStatusBanner(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	body := struct {
		ExpireAt string `json:"expire_at"`
	}{}

	if err := c.BindFunc(&body, "ExpireAt"); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchUpdateStatusBanner(c.Ctx, int64(id), body.ExpireAt)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteBanner(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteBanner(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @client

func (uc Controller) ClientGetBannerList(c *web.Context) error {
	var filter banner.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if lat, ok := c.GetQueryFunc(reflect.Float64, "lat").(*float64); ok {
		filter.Lat = lat
	}
	if lon, ok := c.GetQueryFunc(reflect.Float64, "lon").(*float64); ok {
		filter.Lon = lon
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	list, count, err := uc.useCase.ClientGetBannerList(c.Ctx, filter)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) ClientGetBannerDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	c.Set("lang", c.GetHeader("Accept-Language"))

	res, err := uc.useCase.ClientGetBannerDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   res,
		"status": true,
	}, http.StatusOK)
}

// @super-admin

func (uc Controller) SuperAdminGetBannerList(c *web.Context) error {
	var filter banner.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}
	if status, ok := c.GetQueryFunc(reflect.String, "status").(*string); ok {
		filter.Status = status
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	list, count, err := uc.useCase.SuperAdminGetBannerList(c.Ctx, filter)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminGetBannerDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	c.Set("lang", c.GetHeader("Accept-Language"))

	res, err := uc.useCase.SuperAdminGetBannerDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   res,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminUpdateBannerStatus(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)
	body := struct {
		Status string `json:"status"`
	}{}

	if err := c.BindFunc(&body, "Status"); err != nil {
		return c.RespondError(err)
	}

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.SuperAdminUpdateBannerStatus(c.Ctx, int64(id), body.Status)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// #feedback

// @admin

func (uc Controller) AdminGetFeedBackList(c *web.Context) error {
	var filter feedback.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetFeedbackList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminGetFeedBackByID(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	res, err := uc.useCase.AdminGetFeedbackByID(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   res,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateFeedBack(c *web.Context) error {
	var request feedback.AdminCreate

	if err := c.BindFunc(&request, "Name"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateFeedback(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateColumnFeedBack(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}
	var request feedback.AdminUpdate

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}
	request.ID = int64(id)

	err := uc.useCase.AdminUpdateFeedbackColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteFeedBack(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteFeedback(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @client

func (uc Controller) ClientGetFeedBackList(c *web.Context) error {
	var filter feedback.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.ClientGetFeedbackList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

// #basket

// @client

func (uc Controller) ClientGetBasket(c *web.Context) error {
	branchID := c.GetParam(reflect.Int, "branch-id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	authStr := c.Request.Header.Get("Authorization")

	detail, err := uc.useCase.ClientGetBasket(c.Ctx, int64(branchID), authStr)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": detail,
		},
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) ClientUpdateBasket(c *web.Context) error {
	var request basket.Update

	if err := c.BindFunc(&request, "Menu", "BranchID", "TableID"); err != nil {
		return c.RespondMobileError(err)
	}

	authStr := c.Request.Header.Get("Authorization")

	err := uc.useCase.ClientUpdateBasket(c.Ctx, request, authStr)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

func (uc Controller) ClientDeleteBasket(c *web.Context) error {
	branchID := c.GetParam(reflect.Int, "branch-id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondMobileError(err)
	}

	authStr := c.Request.Header.Get("Authorization")

	err := uc.useCase.ClientDeleteBasket(c.Ctx, int64(branchID), authStr)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
		"error":  nil,
	}, http.StatusOK)
}

//func (uc Controller) ClientCreateBasket(c *web.Context) error {
//	var request basket.CreateRequest
//
//	if err := c.BindFunc(&request, "Value", "BranchID"); err != nil {
//		return c.RespondMobileError(err)
//	}
//
//	var expiration time.Duration
//	if request.Expiration != nil && *request.Expiration != 0 {
//		expiration = time.Duration(rand.Int31n(int32(*request.Expiration))) * time.Second
//	} else {
//		expiration = 0
//	}
//
//	authStr := c.Request.Header.Get("Authorization")
//
//	data := basket.AdminCreate{
//		Value:      request.Value,
//		Expiration: &expiration,
//		BranchID:   request.BranchID,
//	}
//	key, err := uc.useCase.ClientCreateBasket(c.Ctx, data, authStr)
//	if err != nil {
//		return c.RespondMobileError(err)
//	}
//
//	data.Key = &key
//
//	return c.Respond(map[string]interface{}{
//		"data":   data,
//		"status": true,
//		"error":  nil,
//	}, http.StatusOK)
//}

// #full-search

// @client

func (uc Controller) ClientGetFullSearchList(c *web.Context) error {
	var filter full_search.Filter

	if search, ok := c.GetQueryFunc(reflect.String, "search").(*string); ok {
		filter.Search = search
	}
	if menu, ok := c.GetQueryFunc(reflect.String, "menu").(*string); ok {
		filter.Menu = menu
	}
	if lon, ok := c.GetQueryFunc(reflect.Float64, "lon").(*float64); ok {
		filter.Lon = lon
	}
	if lat, ok := c.GetQueryFunc(reflect.Float64, "lat").(*float64); ok {
		filter.Lat = lat
	}
	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondMobileError(err)
	}

	if filter.Search == nil {
		return c.RespondMobileError(errors.New("search is required"))
	}

	list, err := uc.useCase.ClientGetFullSearchList(c.Ctx, filter)
	if err != nil {
		return c.RespondMobileError(err)
	}

	return c.Respond(map[string]interface{}{
		"error":  nil,
		"data":   list,
		"status": true,
	}, http.StatusOK)
}

// #notification

// @admin

func (uc Controller) AdminGetNotificationList(c *web.Context) error {
	var filter notification_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if status, ok := c.GetQueryFunc(reflect.String, "status").(*string); ok {
		filter.Status = status
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetNotificationList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminGetNotificationDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetNotificationDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateNotification(c *web.Context) error {
	var request notification_service.AdminCreateRequest

	if err := c.BindFunc(&request, "Title", "Description"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateNotification(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateNotificationAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request notification_service.AdminUpdateRequest

	if err := c.BindFunc(&request, "Title", "Description"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateNotification(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateNotificationColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request notification_service.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateColumnNotification(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteNotification(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteNotification(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// @super-admin

func (uc Controller) SuperAdminUpdateNotificationStatus(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	body := struct {
		Status string `json:"status"`
	}{}

	if err := c.BindFunc(&body, "Status"); err != nil {
		return c.RespondError(err)
	}

	ID := int64(id)

	err := uc.useCase.SuperAdminUpdateNotificationStatus(c.Ctx, ID, body.Status)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminGetNotificationList(c *web.Context) error {
	var filter notification_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if status, ok := c.GetQueryFunc(reflect.String, "status").(*string); ok {
		filter.Status = status
	}
	if whose, ok := c.GetQueryFunc(reflect.String, "whose").(*string); ok {
		filter.Whose = whose
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.SuperAdminGetNotificationList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) SuperAdminSendNotification(c *web.Context) error {
	var request notification_service.SuperAdminSendRequest

	if err := c.BindFunc(&request, "Title", "Description", "Status"); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.SuperAdminSendNotification(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "sent",
		"status": true,
	}, http.StatusOK)
}

// @client

func (uc Controller) ClientGetNotificationList(c *web.Context) error {
	var filter notification_service.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if page, ok := c.GetQueryFunc(reflect.Int, "page").(*int); ok {
		filter.Page = page
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	//c.Set("lang", c.GetHeader("Accept-Language"))

	list, count, err := uc.useCase.ClientGetNotificationList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientGetCountUnseenNotification(c *web.Context) error {
	count, err := uc.useCase.ClientGetCountUnseenNotification(c.Ctx)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   count,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientSetNotificationAsViewed(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.ClientSetNotificationAsViewed(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) ClientSetAllNotificationsAsViewed(c *web.Context) error {
	err := uc.useCase.ClientSetAllNotificationsAsViewed(c.Ctx)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// #service-percentage

// @branch

func (uc Controller) BranchGetServicePercentageList(c *web.Context) error {
	var filter service_percentage.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetServicePercentageList(c.Ctx, filter)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchGetServicePercentageByID(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	res, err := uc.useCase.AdminGetServicePercentageDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   res,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateServicePercentage(c *web.Context) error {
	var request service_percentage.AdminCreateRequest

	if err := c.BindFunc(&request, "Percent"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateServicePercentage(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateAllServicePercentage(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}
	var request service_percentage.AdminUpdateRequest

	if err := c.BindFunc(&request, "Percent"); err != nil {
		return c.RespondError(err)
	}
	request.ID = int64(id)

	err := uc.useCase.AdminUpdateServicePercentageAll(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteServicePercentage(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteServicePercentage(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}
