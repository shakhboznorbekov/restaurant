package warehouse

import "github.com/restaurant/foundation/web"

type Controller struct {
	useCase *warehouse.UseCase
}

func NewController(useCase *warehouse.UseCase) *Controller {
	return &Controller{useCase}
}

// @warehouse

// #admin

func (uc Controller) AdminGetWarehouseList(c *web.Context) error {
	var filter warehouse2.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if branchId, ok := c.GetQueryFunc(reflect.Int64, "branch_id").(*int64); ok {
		filter.BranchID = branchId
	}
	if Type, ok := c.GetQueryFunc(reflect.String, "type").(*string); ok {
		filter.Type = Type
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetWarehouseList(c.Ctx, filter)
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

func (uc Controller) AdminGetWarehouseDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetWarehouseDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminCreateWarehouse(c *web.Context) error {
	var request warehouse2.AdminCreateRequest

	if err := c.BindFunc(&request, "Name", "Type", "BranchID"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateWarehouse(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateWarehouseAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request warehouse2.AdminUpdateRequest

	if err := c.BindFunc(&request, "Name", "Location", "Type", "BranchID"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateWarehouseAll(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateWarehouseColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request warehouse2.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.AdminUpdateWarehouseColumns(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteWarehouse(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteWarehouse(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// #branch

func (uc Controller) BranchGetWarehouseList(c *web.Context) error {
	var filter warehouse2.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}
	if Type, ok := c.GetQueryFunc(reflect.String, "type").(*string); ok {
		filter.Type = Type
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetWarehouseList(c.Ctx, filter)
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

func (uc Controller) BranchGetWarehouseDetail(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetWarehouseDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchCreateWarehouse(c *web.Context) error {
	var request warehouse2.BranchCreateRequest

	if err := c.BindFunc(&request, "Name", "Type"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateWarehouse(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateWarehouseAll(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request warehouse2.BranchUpdateRequest

	if err := c.BindFunc(&request, "Name", "Location", "Type"); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateWarehouseAll(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateWarehouseColumns(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request warehouse2.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	request.ID = int64(id)

	err := uc.useCase.BranchUpdateWarehouseColumns(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteWarehouse(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteWarehouse(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

// warehouse-transaction

//@admin

func (uc Controller) AdminCreateWarehouseTransaction(c *web.Context) error {
	var request warehouse_transaction.AdminCreateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateWarehouseTransaction(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateWarehouseTransactionColumn(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request warehouse_transaction.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	Id := int64(id)
	request.ID = &Id

	err := uc.useCase.AdminUpdateWarehouseTransactionColumns(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteWarehouseTransaction(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteWarehouseTransaction(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminGetWarehouseTransactionList(c *web.Context) error {
	var filter warehouse_transaction.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetWarehouseTransactionList(c.Ctx, filter)
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

func (uc Controller) AdminGetWarehouseTransactionByID(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetWarehouseTransactionDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

//@branch

func (uc Controller) BranchCreateWarehouseTransaction(c *web.Context) error {
	var request warehouse_transaction.BranchCreateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateWarehouseTransaction(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateWarehouseTransactionColumn(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request warehouse_transaction.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	Id := int64(id)
	request.ID = &Id

	err := uc.useCase.BranchUpdateWarehouseTransactionColumns(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteWarehouseTransaction(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteWarehouseTransaction(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchGetWarehouseTransactionList(c *web.Context) error {
	var filter warehouse_transaction.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetWarehouseTransactionList(c.Ctx, filter)
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

func (uc Controller) BranchGetWarehouseTransactionByID(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetWarehouseTransactionDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

// warehouse-state

//@admin

func (uc Controller) AdminGetWarehouseStateByWarehouseIdList(c *web.Context) error {
	id := c.GetParam(reflect.Int, "warehouse-id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var filter warehouse_state.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetWarehouseStateByWarehouseIDList(c.Ctx, int64(id), filter)
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

//@branch

func (uc Controller) BranchGetWarehouseStateByWarehouseIdList(c *web.Context) error {
	id := c.GetParam(reflect.Int, "warehouse-id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var filter warehouse_state.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetWarehouseStateByWarehouseIDList(c.Ctx, int64(id), filter)
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

// warehouse-transaction-product

//@admin

func (uc Controller) AdminCreateWarehouseTransactionProduct(c *web.Context) error {
	var request warehouse_transaction_product.AdminCreateRequest

	if err := c.BindFunc(&request, "ProductID", "TransactionId", "Amount"); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminCreateWarehouseTransactionProduct(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminUpdateWarehouseTransactionProductColumn(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request warehouse_transaction_product.AdminUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	Id := int64(id)
	request.ID = &Id

	err := uc.useCase.AdminUpdateWarehouseTransactionProductColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminDeleteWarehouseTransactionProduct(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.AdminDeleteWarehouseTransactionProduct(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) AdminGetWarehouseTransactionProductList(c *web.Context) error {
	transactionId := c.GetParam(reflect.Int, "transaction-id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}
	var filter warehouse_transaction_product.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetWarehouseTransactionProductList(c.Ctx, filter, int64(transactionId))
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

func (uc Controller) AdminGetWarehouseTransactionProductByID(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.AdminGetWarehouseTransactionDetail(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

//@branch

func (uc Controller) BranchCreateWarehouseTransactionProduct(c *web.Context) error {
	var request warehouse_transaction_product.BranchCreateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchCreateWarehouseTransactionProduct(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchUpdateWarehouseTransactionProductColumn(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var request warehouse_transaction_product.BranchUpdateRequest

	if err := c.BindFunc(&request); err != nil {
		return c.RespondError(err)
	}

	Id := int64(id)
	request.ID = &Id

	err := uc.useCase.BranchUpdateWarehouseTransactionProductColumn(c.Ctx, request)
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchDeleteWarehouseTransactionProduct(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	err := uc.useCase.BranchDeleteWarehouseTransactionProduct(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   "ok!",
		"status": true,
	}, http.StatusOK)
}

func (uc Controller) BranchGetWarehouseTransactionProductList(c *web.Context) error {
	transactionId := c.GetParam(reflect.Int, "transaction-id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}
	var filter warehouse_transaction_product.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetWarehouseTransactionProductList(c.Ctx, filter, int64(transactionId))
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

func (uc Controller) BranchGetWarehouseTransactionProductByID(c *web.Context) error {
	id := c.GetParam(reflect.Int, "id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	response, err := uc.useCase.BranchGetWarehouseTransactionProductByID(c.Ctx, int64(id))
	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"data":   response,
		"status": true,
	}, http.StatusOK)
}

// warehouse-state

//@admin

func (uc Controller) AdminGetWarehouseStateByWarehouseProductIdList(c *web.Context) error {
	id := c.GetParam(reflect.Int, "warehouse-id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var filter warehouse_state.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.AdminGetWarehouseStateByWarehouseIDList(c.Ctx, int64(id), filter)
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

//@branch

func (uc Controller) BranchGetWarehouseStateByWarehouseProductIdList(c *web.Context) error {
	id := c.GetParam(reflect.Int, "warehouse-id").(int)

	if err := c.ValidParam(); err != nil {
		return c.RespondError(err)
	}

	var filter warehouse_state.Filter

	if limit, ok := c.GetQueryFunc(reflect.Int, "limit").(*int); ok {
		filter.Limit = limit
	}
	if offset, ok := c.GetQueryFunc(reflect.Int, "offset").(*int); ok {
		filter.Offset = offset
	}

	if err := c.ValidQuery(); err != nil {
		return c.RespondError(err)
	}

	list, count, err := uc.useCase.BranchGetWarehouseStateByWarehouseIDList(c.Ctx, int64(id), filter)
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
