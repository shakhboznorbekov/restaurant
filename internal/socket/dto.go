package socket

type Message struct {
	UserID       int64
	BranchID     int64
	RestaurantID int64
	Action       string
	Message      []byte
}

type ResMessage struct {
	Message []byte
}

type Action struct {
	Action string `json:"action"`
}

type Response struct {
	Action   string `json:"action"`
	Property any    `json:"property"`
}

type NewFoodPrinterResponse struct {
	OrderNumber *int64        `json:"order_number"`
	Ip          *string       `json:"ip"`
	Foods       []FoodPrinter `json:"foods"`
}

type FoodPrinter struct {
	Name  *string `json:"name"`
	Count *int    `json:"count"`
}

type NewOrderResponse struct {
	OrderID     int64    `json:"id"`
	OrderNumber int64    `json:"number"`
	TableID     int64    `json:"table_id"`
	TableNumber int64    `json:"table_number"`
	Foods       any      `json:"foods"`
	Price       *float32 `json:"price"`
	CreatedAt   *string  `json:"created_at"`
}

type NewFoodsResponse struct {
	OrderID     int64    `json:"id"`
	OrderNumber int64    `json:"number"`
	TableID     int64    `json:"table_id"`
	TableNumber int64    `json:"table_number"`
	Foods       any      `json:"foods"`
	Price       *float32 `json:"price"`
}

type OrderPaymentResponse struct {
	OrderID     int64 `json:"id"`
	OrderNumber int64 `json:"number"`
}

type NewOrderWaiterResponse struct {
	OrderID     int64    `json:"id"`
	OrderNumber int64    `json:"number"`
	TableID     int64    `json:"table_id"`
	TableNumber int64    `json:"table_number"`
	Foods       any      `json:"foods"`
	ClientCount int      `json:"client_count"`
	Price       *float32 `json:"price"`
	CreatedAt   *string  `json:"created_at"`
	Time        int      `json:"time"`
}

type OrderAcceptedResponse struct {
	OrderID     *int64  `json:"id"`
	OrderNumber *int64  `json:"number"`
	BranchID    *int64  `json:"branch_id"`
	TableID     *int64  `json:"table_id"`
	TableNumber *int64  `json:"table_number"`
	ClientID    *int64  `json:"client_id"`
	Status      string  `json:"status"`
	WaiterID    *int64  `json:"waiter_id"`
	WaiterName  *string `json:"waiter_name"`
	WaiterPhoto *string `json:"waiter_photo"`
	AcceptedAt  *string `json:"accepted_at"`
}

type ClientCall struct {
	OrderID     *int64 `json:"id"`
	OrderNumber *int64 `json:"number"`
	TableID     *int64 `json:"table_id"`
	TableNumber *int64 `json:"table_number"`
}
