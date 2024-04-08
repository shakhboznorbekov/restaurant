package router

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/restaurant/internal/controller/http/v1/file"
	attendance2 "github.com/restaurant/internal/repository/postgres/attendance"
	"github.com/restaurant/internal/repository/postgres/banner"
	"github.com/restaurant/internal/repository/postgres/branchReview"
	"github.com/restaurant/internal/repository/postgres/cashier"
	"github.com/restaurant/internal/repository/postgres/category"
	"github.com/restaurant/internal/repository/postgres/device"
	"github.com/restaurant/internal/repository/postgres/feedback"
	"github.com/restaurant/internal/repository/postgres/food_recipe"
	"github.com/restaurant/internal/repository/postgres/food_recipe_group"
	"github.com/restaurant/internal/repository/postgres/food_recipe_group_history"
	"github.com/restaurant/internal/repository/postgres/full_search"
	"github.com/restaurant/internal/repository/postgres/hall"
	"github.com/restaurant/internal/repository/postgres/menu"
	"github.com/restaurant/internal/repository/postgres/menu_category"
	"github.com/restaurant/internal/repository/postgres/notification"
	"github.com/restaurant/internal/repository/postgres/order"
	"github.com/restaurant/internal/repository/postgres/order_report"
	"github.com/restaurant/internal/repository/postgres/partner"
	"github.com/restaurant/internal/repository/postgres/printers"
	"github.com/restaurant/internal/repository/postgres/product_recipe"
	"github.com/restaurant/internal/repository/postgres/product_recipe_group"
	"github.com/restaurant/internal/repository/postgres/product_recipe_group_history"
	"github.com/restaurant/internal/repository/postgres/service_percentage"
	"github.com/restaurant/internal/repository/postgres/table"
	"github.com/restaurant/internal/repository/postgres/waiter"
	"github.com/restaurant/internal/repository/postgres/waiter_work_time"
	"github.com/restaurant/internal/repository/postgres/warehouse"
	"github.com/restaurant/internal/repository/postgres/warehouse_state"
	"github.com/restaurant/internal/repository/postgres/warehouse_state_history"
	"github.com/restaurant/internal/repository/postgres/warehouse_transaction"
	warehouse_transaction_product "github.com/restaurant/internal/repository/postgres/warehouse_transaction_products"
	"github.com/restaurant/internal/repository/redis/basket"
	"github.com/restaurant/internal/repository/redis/hashing"
	"github.com/restaurant/internal/service/attendance"
	banner2 "github.com/restaurant/internal/service/banner"
	basket2 "github.com/restaurant/internal/service/basket"
	cashier2 "github.com/restaurant/internal/service/cashier"
	category2 "github.com/restaurant/internal/service/category"
	device2 "github.com/restaurant/internal/service/device"
	"github.com/restaurant/internal/service/fcm"
	feedback2 "github.com/restaurant/internal/service/feedback"
	recipe2 "github.com/restaurant/internal/service/food_recipe"
	food_recipe_group2 "github.com/restaurant/internal/service/food_recipe_group"
	food_recipe_group_history2 "github.com/restaurant/internal/service/food_recipe_group_history"
	full_search2 "github.com/restaurant/internal/service/full_search"
	halls "github.com/restaurant/internal/service/hall"
	hashing2 "github.com/restaurant/internal/service/hashing"
	order_report2 "github.com/restaurant/internal/service/order_report"
	categoryrecipe2 "github.com/restaurant/internal/service/product_recipe"
	product_recipe_group2 "github.com/restaurant/internal/service/product_recipe_group"
	product_recipe_group_history2 "github.com/restaurant/internal/service/product_recipe_group_history"
	waiter2 "github.com/restaurant/internal/service/waiter"
	waiter_work_time2 "github.com/restaurant/internal/service/waiter_work_time"
	warehouse2 "github.com/restaurant/internal/service/warehouse"
	warehouse_state2 "github.com/restaurant/internal/service/warehouse_state"
	warehouse_state_history2 "github.com/restaurant/internal/service/warehouse_state_history"
	warehouse_transaction2 "github.com/restaurant/internal/service/warehouse_transaction"
	warehouse_transaction_product2 "github.com/restaurant/internal/service/warehouse_transaction_product"
	"github.com/restaurant/internal/socket"
	"github.com/restaurant/internal/usecase/food"
	warehouse3 "github.com/restaurant/internal/usecase/warehouse"
	"github.com/robfig/cron/v3"
	"log"

	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	catalog_controller "github.com/restaurant/internal/controller/http/v1/catalog"
	food_controller "github.com/restaurant/internal/controller/http/v1/food"
	restaurant_controller "github.com/restaurant/internal/controller/http/v1/restaurant"
	user4 "github.com/restaurant/internal/controller/http/v1/user"
	warehouse_controller "github.com/restaurant/internal/controller/http/v1/warehouse"
	ws_auth_controller "github.com/restaurant/internal/controller/ws/v1/auth"
	"github.com/restaurant/internal/middleware"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"github.com/restaurant/internal/repository/postgres/branch"
	"github.com/restaurant/internal/repository/postgres/district"
	food1 "github.com/restaurant/internal/repository/postgres/food"
	"github.com/restaurant/internal/repository/postgres/foodCategory"
	"github.com/restaurant/internal/repository/postgres/measureUnit"
	orderFood "github.com/restaurant/internal/repository/postgres/order_menu"
	orderPayment "github.com/restaurant/internal/repository/postgres/order_payment"
	"github.com/restaurant/internal/repository/postgres/product"
	"github.com/restaurant/internal/repository/postgres/region"
	"github.com/restaurant/internal/repository/postgres/restaurant"
	"github.com/restaurant/internal/repository/postgres/story"
	"github.com/restaurant/internal/repository/postgres/user"
	redisBranch "github.com/restaurant/internal/repository/redis/branch"
	"github.com/restaurant/internal/service/sms"

	restaurantCategory "github.com/restaurant/internal/repository/postgres/restaurant_category"
	catalog_usecase "github.com/restaurant/internal/usecase/catalog"

	auth3 "github.com/restaurant/internal/controller/http/v1/auth"
	branch2 "github.com/restaurant/internal/service/branch"
	branchReview2 "github.com/restaurant/internal/service/branchReview"
	district2 "github.com/restaurant/internal/service/district"
	food2 "github.com/restaurant/internal/service/food"
	foodCategory2 "github.com/restaurant/internal/service/foodCategory"
	measure_unit2 "github.com/restaurant/internal/service/measureUnit"
	menu2 "github.com/restaurant/internal/service/menu"
	menu_category2 "github.com/restaurant/internal/service/menu_category"
	notification2 "github.com/restaurant/internal/service/notification"
	order2 "github.com/restaurant/internal/service/order"
	orderFood2 "github.com/restaurant/internal/service/order_menu"
	orderPayment2 "github.com/restaurant/internal/service/order_payment"
	partner2 "github.com/restaurant/internal/service/partner"
	printers2 "github.com/restaurant/internal/service/printers"
	product2 "github.com/restaurant/internal/service/product"
	region2 "github.com/restaurant/internal/service/region"
	restaurant2 "github.com/restaurant/internal/service/restaurant"
	restaurantCategory2 "github.com/restaurant/internal/service/restaurant_category"
	service_percentage2 "github.com/restaurant/internal/service/service_percentage"
	story2 "github.com/restaurant/internal/service/story"
	table2 "github.com/restaurant/internal/service/tables"
	user2 "github.com/restaurant/internal/service/user"
	auth2 "github.com/restaurant/internal/usecase/auth"
	order3 "github.com/restaurant/internal/usecase/order"
	partner3 "github.com/restaurant/internal/usecase/partner"
	product3 "github.com/restaurant/internal/usecase/product"
	restaurant3 "github.com/restaurant/internal/usecase/restaurant"
	user3 "github.com/restaurant/internal/usecase/user"

	order4 "github.com/restaurant/internal/controller/http/v1/order"
	partner4 "github.com/restaurant/internal/controller/http/v1/partner"
	product4 "github.com/restaurant/internal/controller/http/v1/product"
)

type Router struct {
	*web.App
	postgresDB *postgresql.Database
	redisDB    *redis.Client
	fcm        *fcm.ConfigFCM
	port       string
	auth       *auth.Auth
	smsConfig  *sms.Config
	hub        *socket.Hub
}

func NewRouter(
	app *web.App,
	postgresDB *postgresql.Database,
	redisDB *redis.Client,
	fcm *fcm.ConfigFCM,
	port string,
	auth *auth.Auth,
	smsConfig *sms.Config,
	hub *socket.Hub,
) *Router {
	return &Router{
		app,
		postgresDB,
		redisDB,
		fcm,
		port,
		auth,
		smsConfig,
		hub,
	}
}

func (r Router) Init() error {
	// repositories:

	// - postgresql
	userPostgres := user.NewRepository(r.postgresDB)
	measureUnitPostgres := measureUnit.NewRepository(r.postgresDB)
	productPostgres := product.NewRepository(r.postgresDB)
	restaurantPostgres := restaurant.NewRepository(r.postgresDB)
	regionPostgres := region.NewRepository(r.postgresDB)
	districtPostgres := district.NewRepository(r.postgresDB)
	restaurantCategoryPostgres := restaurantCategory.NewRepository(r.postgresDB)
	branchPostgres := branch.NewRepository(r.postgresDB)
	foodCategoryPostgres := foodCategory.NewRepository(r.postgresDB)
	foodPostgres := food1.NewRepository(r.postgresDB)
	menuPostgres := menu.NewRepository(r.postgresDB)
	tablePostgres := table.NewRepository(r.postgresDB)
	orderPostgres := order.NewRepository(r.postgresDB)
	printersPostgres := printers.NewRepository(r.postgresDB)
	branchReviewPostgres := branchReview.NewRepository(r.postgresDB)
	partnerPostgres := partner.NewRepository(r.postgresDB)
	orderFoodPostgres := orderFood.NewRepository(r.postgresDB)
	orderPaymentPostgres := orderPayment.NewRepository(r.postgresDB)
	storyPostgres := story.NewRepository(r.postgresDB)
	bannerPostgres := banner.NewRepository(r.postgresDB)
	feedbackPostgres := feedback.NewRepository(r.postgresDB)
	fullSearchPostgres := full_search.NewRepository(r.postgresDB)
	foodRecipePostgres := food_recipe.NewRepository(r.postgresDB)
	notificationPostgres := notification.NewRepository(r.postgresDB)
	devicePostgres := device.NewRepository(r.postgresDB)
	warehousePostgres := warehouse.NewRepository(r.postgresDB)
	warehouseStatePostgres := warehouse_state.NewRepository(r.postgresDB)
	warehouseTransactionsPostgres := warehouse_transaction.NewRepository(r.postgresDB)
	warehouseStateHistoryPostgres := warehouse_state_history.NewRepository(r.postgresDB)
	waiterPostgres := waiter.NewRepository(r.postgresDB)
	categoryRecipePostgres := product_recipe.NewRepository(r.postgresDB)
	servicePercentagePostgres := service_percentage.NewRepository(r.postgresDB)
	attendancePostgres := attendance2.NewRepository(r.postgresDB)
	waiterWorkTimePostgres := waiter_work_time.NewRepository(r.postgresDB)
	warehouseTransactionProductPostgres := warehouse_transaction_product.NewRepository(r.postgresDB)
	categoryPostgres := category.NewRepository(r.postgresDB)
	menuCategoryPostgres := menu_category.NewRepository(r.postgresDB)
	cashierPostgres := cashier.NewRepository(r.postgresDB)
	foodRecipeGroupPostgres := food_recipe_group.NewRepository(r.postgresDB)
	foodRecipeGroupHistoryPostgres := food_recipe_group_history.NewRepository(r.postgresDB)
	productRecipeGroupPostgres := product_recipe_group.NewRepository(r.postgresDB)
	productRecipeGroupHistoryPostgres := product_recipe_group_history.NewRepository(r.postgresDB)
	hallPostgres := hall.NewRepository(r.postgresDB)
	orderReportPostgres := order_report.New(r.postgresDB)

	// - redis
	branchRedis := redisBranch.NewRepository(r.redisDB)
	basketRedis := basket.NewRepository(r.redisDB)
	hashingRedis := hashing.NewRepository(r.redisDB)

	hashing2.Hashing = hashingRedis

	// - service
	userService := user2.NewService(userPostgres)
	restaurantService := restaurant2.NewService(restaurantPostgres)
	smsService := sms.NewService(r.smsConfig, r.redisDB, r.postgresDB)
	measureUnitService := measure_unit2.NewService(measureUnitPostgres)
	productService := product2.NewService(productPostgres)
	regionService := region2.NewService(regionPostgres)
	districtService := district2.NewService(districtPostgres)
	restaurantCategoryService := restaurantCategory2.NewService(restaurantCategoryPostgres)
	branchService := branch2.NewService(branchPostgres, branchRedis)
	foodCategoryService := foodCategory2.NewService(foodCategoryPostgres)
	foodService := food2.NewService(foodPostgres)
	menuService := menu2.NewService(menuPostgres)
	tableService := table2.NewService(tablePostgres)
	orderService := order2.NewService(orderPostgres)
	printersService := printers2.NewService(printersPostgres)
	branchReviewService := branchReview2.NewService(branchReviewPostgres)
	partnerService := partner2.NewService(partnerPostgres)
	orderFoodService := orderFood2.NewService(orderFoodPostgres)
	orderPaymentService := orderPayment2.NewService(orderPaymentPostgres)
	storyService := story2.NewService(storyPostgres)
	bannerService := banner2.NewService(bannerPostgres)
	feedbackService := feedback2.NewService(feedbackPostgres)
	basketService := basket2.NewService(basketRedis)
	recipeService := recipe2.NewService(foodRecipePostgres)
	notificationService := notification2.NewService(notificationPostgres)
	fullSearchService := full_search2.NewService(fullSearchPostgres)
	deviceService := device2.NewService(devicePostgres)
	warehouseService := warehouse2.NewService(warehousePostgres)
	warehouseStateService := warehouse_state2.NewService(warehouseStatePostgres)
	warehouseTransactionService := warehouse_transaction2.NewService(warehouseTransactionsPostgres)
	warehouseStateHistoryService := warehouse_state_history2.NewService(warehouseStateHistoryPostgres)
	fcmService := fcm.NewFCMService(r.fcm)
	waiterService := waiter2.NewService(waiterPostgres)
	productRecipeService := categoryrecipe2.NewService(categoryRecipePostgres)
	servicePercentageService := service_percentage2.NewService(servicePercentagePostgres)
	attendanceService := attendance.NewService(attendancePostgres)
	waiterWorkTimeService := waiter_work_time2.NewService(waiterWorkTimePostgres)
	warehouseTransactionProductService := warehouse_transaction_product2.NewService(warehouseTransactionProductPostgres)
	categoryService := category2.NewService(categoryPostgres)
	menuCategoryService := menu_category2.NewService(menuCategoryPostgres)
	foodRecipeGroupService := food_recipe_group2.NewService(foodRecipeGroupPostgres)
	foodRecipeGroupHistoryService := food_recipe_group_history2.NewService(foodRecipeGroupHistoryPostgres)
	productRecipeGroupService := product_recipe_group2.NewService(productRecipeGroupPostgres)
	productRecipeGroupHistoryService := product_recipe_group_history2.NewService(productRecipeGroupHistoryPostgres)
	hallService := halls.NewService(hallPostgres)
	orderReportService := order_report2.New(orderReportPostgres)

	// cron service for scheduled tasks
	schedule := cron.New()
	cashierService := cashier2.NewService(cashierPostgres)

	// - use_case
	userUseCase := user3.NewUseCase(userService, smsService, waiterService, servicePercentageService, attendanceService, waiterWorkTimeService, r.auth, schedule, cashierService)
	if err := userUseCase.CalculateWaitersKPI(context.TODO()); err != nil {
		log.Println("cron use_case:", err)
	}

	restaurantUseCase := restaurant3.NewUseCase(
		restaurantService,
		userService,
		restaurantCategoryService,
		branchService,
		tableService,
		branchReviewService,
		printersService,
		servicePercentageService,
		hallService,
	)
	authUseCase := auth2.NewUseCase(userService, smsService, deviceService, r.auth)
	catalogUseCase := catalog_usecase.NewUseCase(
		measureUnitService, regionService,
		districtService, storyService,
		bannerService, feedbackService,
		basketService,
		fullSearchService,
		menuService,
		r.auth,
		notificationService,
		fcmService,
		deviceService,
		servicePercentageService,
	)
	productUseCase := product3.NewUseCase(productService, productRecipeService, productRecipeGroupService, productRecipeGroupHistoryService)
	foodUseCase := food.NewUseCase(
		foodCategoryService, foodService,
		menuService, recipeService,
		basketService, categoryService,
		menuCategoryService, foodRecipeGroupService,
		foodRecipeGroupHistoryService,
	)
	orderUseCase := order3.NewUseCase(orderService, orderFoodService, orderPaymentService, r.hub, orderReportService)
	partnerUseCase := partner3.NewUseCase(partnerService)
	warehouseUseCase := warehouse3.NewUseCase(
		warehouseService,
		warehouseStateService,
		warehouseTransactionService,
		warehouseStateHistoryService,
		warehouseTransactionProductService,
	)

	// - controller
	userController := user4.NewController(userUseCase)
	restaurantController := restaurant_controller.NewController(restaurantUseCase)
	authController := auth3.NewController(authUseCase)
	catalogController := catalog_controller.NewController(catalogUseCase)
	productController := product4.NewController(productUseCase)
	foodController := food_controller.NewController(foodUseCase)
	fileController := file.NewController(r.App)
	orderController := order4.NewController(orderUseCase)
	partnerController := partner4.NewController(partnerUseCase)
	warehouseController := warehouse_controller.NewController(warehouseUseCase)

	wsAuthController := ws_auth_controller.NewController(r.hub, restaurantUseCase)

	//order checking
	err := orderUseCase.OrderChecking(context.Background(), 10)
	if err != nil {
		log.Println(err.Error())
	}

	// @media
	{
		r.GET("/media/*filepath", fileController.File)
		r.HEAD("/media/*filepath", fileController.File)
		r.RouterGroup.BasePath()
	}

	// @ws
	{
		r.Get("/api/v1/ws/auth", wsAuthController.WebSocket, middleware.WsAuthenticate(r.auth))
		r.Get("/api/v1/ws/printer/auth", wsAuthController.WebSocketPrinter)
	}

	// @auth
	{
		// client
		r.Post("/api/v1/send-sms", authController.SendSmsCode)
		r.Post("/api/v1/sign-in", authController.SignInClient)
		r.Post("/api/v1/client/fill-up", authController.ClientFillUp, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Put("/api/v1/client/update/phone", authController.ClientUpdatePhone, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Delete("/api/v1/sign-out/:device-id", authController.ClientLogOut, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Post("/api/v1/change/lang", authController.ChangeDeviceLang, middleware.Authenticate(r.auth, auth.RoleClient))

		// super_admin - admin - branch - cashier
		r.Post("/api/v1/admin/sign-in", authController.SignIn)

		// waiter
		r.Post("/api/v1/waiter/sign-in", authController.WaiterSignIn)
		r.Patch("/api/v1/waiter/update/phone", authController.WaiterUpdatePhone, middleware.Authenticate(r.auth, auth.RoleWaiter))
	}

	// @super-admin
	{
		// #user
		r.Get("/api/v1/super-admin/user/list", userController.SuperAdminGetUserList, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Get("/api/v1/super-admin/user/:id", userController.SuperAdminGetUserDetail, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Post("/api/v1/super-admin/user/create", userController.SuperAdminCreateUser, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Put("/api/v1/super-admin/user/:id", userController.SuperAdminUpdateUserAll, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Patch("/api/v1/super-admin/user/:id", userController.SuperAdminUpdateUserColumns, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Delete("/api/v1/super-admin/user/:id", userController.SuperAdminDeleteUser, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))

		// #restaurant
		r.Get("/api/v1/super-admin/restaurant/list", restaurantController.SuperAdminGetRestaurantList, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Get("/api/v1/super-admin/restaurant/:id", restaurantController.SuperAdminGetRestaurantDetail, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Post("/api/v1/super-admin/restaurant/create", restaurantController.SuperAdminCreateRestaurant, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Put("/api/v1/super-admin/restaurant/:id", restaurantController.SuperAdminUpdateRestaurantAll, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Patch("/api/v1/super-admin/restaurant/:id", restaurantController.SuperAdminUpdateRestaurantColumns, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Patch("/api/v1/super-admin/restaurant/password/:id", restaurantController.SuperAdminUpdateRestaurantAdmin, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Delete("/api/v1/super-admin/restaurant/:id", restaurantController.SuperAdminDeleteRestaurant, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))

		// #measure-unit
		r.Get("/api/v1/super-admin/measure-unit/list", catalogController.SuperAdminGetMeasureUnitList, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Get("/api/v1/super-admin/measure-unit/:id", catalogController.SuperAdminGetMeasureUnitDetail, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Post("/api/v1/super-admin/measure-unit/create", catalogController.SuperAdminCreateMeasureUnit, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Put("/api/v1/super-admin/measure-unit/:id", catalogController.SuperAdminUpdateMeasureUnitAll, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Patch("/api/v1/super-admin/measure-unit/:id", catalogController.SuperAdminUpdateMeasureUnitColumns, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Delete("/api/v1/super-admin/measure-unit/:id", catalogController.SuperAdminDeleteMeasureUnit, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))

		// #region
		r.Get("/api/v1/super-admin/region/list", catalogController.SuperAdminGetRegionList, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Get("/api/v1/super-admin/region/:id", catalogController.SuperAdminGetRegionDetail, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Post("/api/v1/super-admin/region/create", catalogController.SuperAdminCreateRegion, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Put("/api/v1/super-admin/region/:id", catalogController.SuperAdminUpdateRegionAll, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Patch("/api/v1/super-admin/region/:id", catalogController.SuperAdminUpdateRegionColumns, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Delete("/api/v1/super-admin/region/:id", catalogController.SuperAdminDeleteRegion, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))

		// #district
		r.Get("/api/v1/super-admin/district/list", catalogController.SuperAdminGetDistrictList, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Get("/api/v1/super-admin/district/:id", catalogController.SuperAdminGetDistrictDetail, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Post("/api/v1/super-admin/district/create", catalogController.SuperAdminCreateDistrict, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Put("/api/v1/super-admin/district/:id", catalogController.SuperAdminUpdateDistrictAll, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Patch("/api/v1/super-admin/district/:id", catalogController.SuperAdminUpdateDistrictColumns, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Delete("/api/v1/super-admin/district/:id", catalogController.SuperAdminDeleteDistrict, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))

		// #restaurant_category
		r.Get("/api/v1/super-admin/restaurant-category/list", restaurantController.SuperAdminGetRestaurantCategoryList, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Get("/api/v1/super-admin/restaurant-category/:id", restaurantController.SuperAdminGetRestaurantCategoryDetail, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Post("/api/v1/super-admin/restaurant-category/create", restaurantController.SuperAdminCreateRestaurantCategory, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Put("/api/v1/super-admin/restaurant-category/:id", restaurantController.SuperAdminUpdateRestaurantCategoryAll, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Patch("/api/v1/super-admin/restaurant-category/:id", restaurantController.SuperAdminUpdateRestaurantCategoryColumns, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Delete("/api/v1/super-admin/restaurant-category/:id", restaurantController.SuperAdminDeleteRestaurantCategory, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))

		// #notification
		r.Post("/api/v1/super-admin/send/notification", catalogController.SuperAdminSendNotification, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Get("/api/v1/super-admin/notification/list", catalogController.SuperAdminGetNotificationList, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Patch("/api/v1/super-admin/notification/update/status/:id", catalogController.SuperAdminUpdateNotificationStatus, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))

		// #banner
		r.Get("/api/v1/super-admin/banner/list", catalogController.SuperAdminGetBannerList, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Get("/api/v1/super-admin/banner/:id", catalogController.SuperAdminGetBannerDetail, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Patch("/api/v1/super-admin/banner/update/status/:id", catalogController.SuperAdminUpdateBannerStatus, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))

		// #story
		r.Get("/api/v1/super-admin/story/list", catalogController.SuperAdminGetStoryList, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Patch("/api/v1/super-admin/story/update/status/:id", catalogController.SuperAdminUpdateStoryStatus, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))

		// #foodCategory
		r.Get("/api/v1/super-admin/category/list", foodController.SuperAdminGetCategoryList, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Get("/api/v1/super-admin/category/:id", foodController.SuperAdminGetCategoryDetail, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Post("/api/v1/super-admin/category/create", foodController.SuperAdminCreateCategory, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Put("/api/v1/super-admin/category/:id", foodController.SuperAdminUpdateCategoryAll, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Patch("/api/v1/super-admin/category/:id", foodController.SuperAdminUpdateCategoryColumns, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Delete("/api/v1/super-admin/category/:id", foodController.SuperAdminDeleteCategory, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
	}

	// @admin
	{
		// #measure_unit
		r.Get("/api/v1/admin/measure-unit/list", catalogController.AdminGetMeasureUnitList, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #product
		r.Get("/api/v1/admin/product/list", productController.AdminGetProductList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/product/:id", productController.AdminGetProductDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/product/create", productController.AdminCreateProduct, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Put("/api/v1/admin/product/:id", productController.AdminUpdateProductAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/product/:id", productController.AdminUpdateProductColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/product/:id", productController.AdminDeleteProduct, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #branch
		r.Get("/api/v1/admin/branch/list", restaurantController.AdminGetBranchList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/branch/:id", restaurantController.AdminGetBranchDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/branch/create", restaurantController.AdminCreateBranch, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Put("/api/v1/admin/branch/:id", restaurantController.AdminUpdateBranchAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/branch/:id", restaurantController.AdminUpdateBranchColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/branch/password/:id", restaurantController.AdminUpdateBranchAdmin, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/branch/:id", restaurantController.AdminDeleteBranch, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/branch/image/delete/:id", restaurantController.AdminImageDeleteBranch, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #food
		r.Get("/api/v1/admin/food/list", foodController.AdminGetFoodList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/food/:id", foodController.AdminGetFoodDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/food/create", foodController.AdminCreateFood, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Put("/api/v1/admin/food/:id", foodController.AdminUpdateFoodAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/food/:id", foodController.AdminUpdateFoodColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/food/:id", foodController.AdminDeleteFood, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/food/image/delete/:id", foodController.AdminImageDelete, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #food_recipe
		r.Get("/api/v1/admin/recipe/list", foodController.AdminGetFoodRecipeList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/recipe/:id", foodController.AdminGetFoodRecipeDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/recipe/create", foodController.AdminCreateFoodRecipe, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Put("/api/v1/admin/recipe/:id", foodController.AdminUpdateFoodRecipeAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/recipe/:id", foodController.AdminUpdateFoodRecipeColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/recipe/:id", foodController.AdminDeleteFoodRecipe, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #menu
		r.Get("/api/v1/admin/menu/list", foodController.AdminGetMenuList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/menu/:id", foodController.AdminGetMenuDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/menu/create", foodController.AdminCreateMenu, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Put("/api/v1/admin/menu/:id", foodController.AdminUpdateMenuAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/menu/:id", foodController.AdminUpdateMenuColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/menu/:id", foodController.AdminDeleteMenu, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/menu/remove/photo/:id", foodController.AdminRemoveMenuPhoto, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #table
		r.Get("/api/v1/admin/table/list", restaurantController.AdminGetTableList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/table/:id", restaurantController.AdminGetTableDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/table/create", restaurantController.AdminCreateTable, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Put("/api/v1/admin/table/:id", restaurantController.AdminUpdateTableAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/table/:id", restaurantController.AdminUpdateTableColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/table/:id", restaurantController.AdminDeleteTable, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #restaurant_category
		r.Get("/api/v1/admin/restaurant-category/list", restaurantController.AdminGetRestaurantCategoryList, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #order
		r.Get("/api/v1/admin/order/list", orderController.AdminGetListOrder, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #story
		r.Get("/api/v1/admin/story/list", catalogController.AdminGetStoryList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/story/create", catalogController.AdminCreateStory, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/story/status/:id", catalogController.AdminUpdateStatusStory, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/story/:id", catalogController.AdminDeleteStory, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #feedback
		r.Get("/api/v1/admin/feedback/list", catalogController.AdminGetFeedBackList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/feedback/:id", catalogController.AdminGetFeedBackByID, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/feedback/create", catalogController.AdminCreateFeedBack, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/feedback/:id", catalogController.AdminUpdateColumnFeedBack, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/feedback/:id", catalogController.AdminDeleteFeedBack, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #notification
		r.Get("/api/v1/admin/notification/list", catalogController.AdminGetNotificationList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/notification/:id", catalogController.AdminGetNotificationDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/notification/create", catalogController.AdminCreateNotification, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Put("/api/v1/admin/notification/:id", catalogController.AdminUpdateNotificationAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/notification/:id", catalogController.AdminUpdateNotificationColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/notification/:id", catalogController.AdminDeleteNotification, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #warehouse
		r.Get("/api/v1/admin/warehouse/list", warehouseController.AdminGetWarehouseList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/warehouse/:id", warehouseController.AdminGetWarehouseDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/warehouse/create", warehouseController.AdminCreateWarehouse, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/warehouse/:id", warehouseController.AdminUpdateWarehouseColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/warehouse/:id", warehouseController.AdminDeleteWarehouse, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #waiter
		r.Get("/api/v1/admin/waiter/list", userController.AdminGetWaiterList, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #partner
		r.Get("/api/v1/admin/partner/list", partnerController.AdminGetPartnerList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/partner/:id", partnerController.AdminGetPartnerDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/partner/create", partnerController.AdminCreatePartner, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Put("/api/v1/admin/partner/:id", partnerController.AdminUpdatePartnerAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/partner/:id", partnerController.AdminUpdatePartnerColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/partner/:id", partnerController.AdminDeletePartner, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #product
		r.Get("/api/v1/admin/product/spending/by/branch", productController.AdminGetProductSpendingByBranch, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #warehouse-transaction
		r.Get("/api/v1/admin/warehouse/transaction/list", warehouseController.AdminGetWarehouseTransactionList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/warehouse/transaction/:id", warehouseController.AdminGetWarehouseTransactionByID, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/warehouse/transaction/create", warehouseController.AdminCreateWarehouseTransaction, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/warehouse/transaction/:id", warehouseController.AdminUpdateWarehouseTransactionColumn, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/warehouse/transaction/:id", warehouseController.AdminDeleteWarehouseTransaction, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #warehouse-transaction-product
		r.Get("/api/v1/admin/warehouse/transaction/product/list/:transaction-id", warehouseController.AdminGetWarehouseTransactionProductList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/warehouse/transaction/product/:id", warehouseController.AdminGetWarehouseTransactionProductByID, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/warehouse/transaction/product/create", warehouseController.AdminCreateWarehouseTransactionProduct, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/warehouse/transaction/product/:id", warehouseController.AdminUpdateWarehouseTransactionProductColumn, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/warehouse/transaction/product/:id", warehouseController.AdminDeleteWarehouseTransactionProduct, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #warehouse-state
		r.Get("/api/v1/admin/warehouse/state/list/by/warehouse/:warehouse-id", warehouseController.AdminGetWarehouseStateByWarehouseIdList, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #product_recipe
		r.Get("/api/v1/admin/product-recipe/list", productController.AdminGetProductRecipeList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/product-recipe/:id", productController.AdminGetProductRecipeDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/product-recipe/create", productController.AdminCreateProductRecipe, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Put("/api/v1/admin/product-recipe/:id", productController.AdminUpdateProductRecipeAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/product-recipe/:id", productController.AdminUpdateProductRecipeColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/product-recipe/:id", productController.AdminDeleteProductRecipe, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #category
		r.Get("/api/v1/admin/category/list", foodController.AdminGetCategoryList, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #cashier
		r.Get("/api/v1/admin/cashier/list", userController.AdminGetCashierList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/cashier/:id", userController.AdminGetCashierDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/cashier/create", userController.AdminCreateCashier, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Put("/api/v1/admin/cashier/:id", userController.AdminUpdateCashierAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/cashier/:id", userController.AdminUpdateCashierColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/cashier/:id", userController.AdminDeleteCashier, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/cashier/status/:id", userController.AdminUpdateCashierStatus, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/cashier/password/:id", userController.AdminUpdateCashierPassword, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/cashier/phone/:id", userController.AdminUpdateCashierPhone, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #menu_category
		r.Get("/api/v1/admin/menu-category/list", foodController.AdminGetMenuCategoryList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/menu-category/:id", foodController.AdminGetMenuCategoryDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/menu-category/create", foodController.AdminCreateMenuCategory, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/menu-category/:id", foodController.AdminUpdateMenuCategory, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/menu-category/:id", foodController.AdminDeleteMenuCategory, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #food_category
		r.Get("/api/v1/admin/food-category/list", foodController.AdminGetFoodCategoryList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/food-category/:id", foodController.AdminGetFoodCategoryDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/food-category/create", foodController.AdminCreateFoodCategory, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Put("/api/v1/admin/food-category/:id", foodController.AdminUpdateFoodCategoryAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/food-category/:id", foodController.AdminUpdateFoodCategoryColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/food-category/:id", foodController.AdminDeleteFoodCategory, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #food_recipe_group
		r.Get("/api/v1/admin/food/recipe/group/list", foodController.AdminGetFoodRecipeGroupList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/food/recipe/group/:id", foodController.AdminGetFoodRecipeGroupDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/food/recipe/group/create", foodController.AdminCreateFoodRecipeGroup, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Put("/api/v1/admin/food/recipe/group/:id", foodController.AdminUpdateFoodRecipeGroupAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/food/recipe/group/:id", foodController.AdminUpdateFoodRecipeGroupColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/food/recipe/group/single/:id", foodController.AdminDeleteFoodRecipeGroupSingleRecipe, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/food/recipe/group/:id", foodController.AdminDeleteFoodRecipeGroup, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #food_recipe_group
		r.Get("/api/v1/admin/food/recipe/group/history/list", foodController.AdminGetFoodRecipeGroupHistoryList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/food/recipe/group/history/:id", foodController.AdminGetFoodRecipeGroupHistoryDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/food/recipe/group/history/create", foodController.AdminCreateFoodRecipeGroupHistory, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/food/recipe/group/history/:id", foodController.AdminDeleteFoodRecipeGroupHistory, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #hall
		r.Get("/api/v1/admin/hall/list", restaurantController.AdminGetHallList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/hall/:id", restaurantController.AdminGetHallDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/hall/create", restaurantController.AdminCreateHall, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Put("/api/v1/admin/hall/:id", restaurantController.AdminUpdateHallAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/hall/:id", restaurantController.AdminUpdateHallColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/hall/:id", restaurantController.AdminDeleteHall, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #product_recipe_group
		r.Get("/api/v1/admin/product/recipe/group/list", productController.AdminGetProductRecipeGroupList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/product/recipe/group/:id", productController.AdminGetProductRecipeGroupDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/product/recipe/group/create", productController.AdminCreateProductRecipeGroup, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Put("/api/v1/admin/product/recipe/group/:id", productController.AdminUpdateProductRecipeGroupAll, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/product/recipe/group/:id", productController.AdminUpdateProductRecipeGroupColumns, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Patch("/api/v1/admin/product/recipe/group/single/:id", productController.AdminDeleteProductRecipeGroupSingleRecipe, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/product/recipe/group/:id", productController.AdminDeleteProductRecipeGroup, middleware.Authenticate(r.auth, auth.RoleAdmin))

		// #product_recipe_group_history
		r.Get("/api/v1/admin/product/recipe/group/history/list", productController.AdminGetProductRecipeGroupHistoryList, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Get("/api/v1/admin/product/recipe/group/history/:id", productController.AdminGetProductRecipeGroupHistoryDetail, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Post("/api/v1/admin/product/recipe/group/history/create", productController.AdminCreateProductRecipeGroupHistory, middleware.Authenticate(r.auth, auth.RoleAdmin))
		r.Delete("/api/v1/admin/product/recipe/group/history/:id", productController.AdminDeleteProductRecipeGroupHistory, middleware.Authenticate(r.auth, auth.RoleAdmin))
	}

	// @branch
	{
		// #banner
		r.Get("/api/v1/branch/banner/list", catalogController.BranchGetBannerList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/banner/:id", catalogController.BranchGetBannerByID, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/banner/create", catalogController.BranchCreateBanner, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/banner/:id", catalogController.BranchUpdateAllBanner, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/banner/:id", catalogController.BranchUpdateColumnBanner, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/banner/status/:id", catalogController.BranchUpdateStatusBanner, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/banner/:id", catalogController.BranchDeleteBanner, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #menu
		r.Get("/api/v1/branch/menu/list", foodController.BranchGetMenuList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/menu/:id", foodController.BranchGetMenuDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/menu/create", foodController.BranchCreateMenu, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/menu/:id", foodController.BranchUpdateMenuAll, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/menu/:id", foodController.BranchUpdateMenuColumns, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/menu/:id", foodController.BranchDeleteMenu, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/menu/add/printer", foodController.BranchUpdateMenuPrinterID, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/menu/printer/:id", foodController.BranchDeleteMenuPrinterID, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/menu/remove/photo/:id", foodController.BranchRemoveMenuPhoto, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #table
		r.Get("/api/v1/branch/table/list", restaurantController.BranchGetTableList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/table/:id", restaurantController.BranchGetTableDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/table/create", restaurantController.BranchCreateTable, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/table/:id", restaurantController.BranchUpdateTableAll, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/table/:id", restaurantController.BranchUpdateTableColumns, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/table/:id", restaurantController.BranchDeleteTable, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/generate/qr", restaurantController.BranchGenerateQRTable, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #food
		r.Get("/api/v1/branch/food/list", foodController.BranchGetFoodList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/food/:id", foodController.BranchGetFoodDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/food/create", foodController.BranchCreateFood, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/food/:id", foodController.BranchUpdateFoodAll, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/food/:id", foodController.BranchUpdateFoodColumns, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/food/:id", foodController.BranchDeleteFood, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/food/image/delete/:id", foodController.BranchImageDelete, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #printers
		r.Get("/api/v1/branch/printer/list", restaurantController.BranchGetPrintersList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/printer/:id", restaurantController.BranchGetPrintersDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/printer/create", restaurantController.BranchCreatePrinters, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/printer/:id", restaurantController.BranchUpdatePrintersAll, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/printer/:id", restaurantController.BranchUpdatePrintersColumns, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/printer/:id", restaurantController.BranchDeletePrinters, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #warehouse
		r.Get("/api/v1/branch/warehouse/list", warehouseController.BranchGetWarehouseList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/warehouse/:id", warehouseController.BranchGetWarehouseDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/warehouse/create", warehouseController.BranchCreateWarehouse, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/warehouse/:id", warehouseController.BranchUpdateWarehouseColumns, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/warehouse/:id", warehouseController.BranchDeleteWarehouse, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #waiters
		r.Get("/api/v1/branch/waiter/list", userController.BranchGetWaiterList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/waiter/:id", userController.BranchGetWaiterDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/waiter/create", userController.BranchCreateWaiter, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/waiter/:id", userController.BranchUpdateWaiterAll, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/waiter/:id", userController.BranchUpdateWaiterColumns, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/waiter/:id", userController.BranchDeleteWaiter, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/waiter/status/:id", userController.BranchUpdateWaiterStatus, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/waiter/password/:id", userController.BranchUpdateWaiterPassword, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/waiter/phone/:id", userController.BranchUpdateWaiterPhone, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #partner
		r.Get("/api/v1/branch/partner/list", partnerController.BranchGetPartnerList, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #service-percentage
		r.Get("/api/v1/branch/service-percentage/list", catalogController.BranchGetServicePercentageList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/service-percentage/:id", catalogController.BranchGetServicePercentageByID, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/service-percentage/create", catalogController.BranchCreateServicePercentage, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/service-percentage/:id", catalogController.BranchUpdateAllServicePercentage, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/service-percentage/:id", catalogController.BranchDeleteServicePercentage, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #warehouse-transaction
		r.Get("/api/v1/branch/warehouse/transaction/list", warehouseController.BranchGetWarehouseTransactionList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/warehouse/transaction/:id", warehouseController.BranchGetWarehouseTransactionByID, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/warehouse/transaction/create", warehouseController.BranchCreateWarehouseTransaction, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/warehouse/transaction/:id", warehouseController.BranchUpdateWarehouseTransactionColumn, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/warehouse/transaction/:id", warehouseController.BranchDeleteWarehouseTransaction, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #warehouse-transaction-product
		r.Get("/api/v1/branch/warehouse/transaction/product/list/:transaction-id", warehouseController.BranchGetWarehouseTransactionProductList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/warehouse/transaction/product/:id", warehouseController.BranchGetWarehouseTransactionProductByID, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/warehouse/transaction/product/create", warehouseController.BranchCreateWarehouseTransactionProduct, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/warehouse/transaction/product/:id", warehouseController.BranchUpdateWarehouseTransactionProductColumn, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/warehouse/transaction/product/:id", warehouseController.BranchDeleteWarehouseTransactionProduct, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #warehouse-state
		r.Get("/api/v1/branch/warehouse/state/list/by/warehouse/:warehouse-id", warehouseController.BranchGetWarehouseStateByWarehouseIdList, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #branch
		r.Get("/api/v1/branch/branch/token", restaurantController.BranchGetBranchToken, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #product
		r.Get("/api/v1/branch/product/list", productController.BranchGetProductList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/product/:id", productController.BranchGetProductDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/product/create", productController.BranchCreateProduct, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/product/:id", productController.BranchUpdateProductAll, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/product/:id", productController.BranchUpdateProductColumns, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/product/:id", productController.BranchDeleteProduct, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #product_recipe
		r.Get("/api/v1/branch/product-recipe/list", productController.BranchGetProductRecipeList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/product-recipe/:id", productController.BranchGetProductRecipeDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/product-recipe/create", productController.BranchCreateProductRecipe, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/product-recipe/:id", productController.BranchUpdateProductRecipeAll, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/product-recipe/:id", productController.BranchUpdateProductRecipeColumns, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/product-recipe/:id", productController.BranchDeleteProductRecipe, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #food_recipe
		r.Get("/api/v1/branch/recipe/list", foodController.BranchGetFoodRecipeList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/recipe/:id", foodController.BranchGetFoodRecipeDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/recipe/create", foodController.BranchCreateFoodRecipe, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/recipe/:id", foodController.BranchUpdateFoodRecipeAll, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/recipe/:id", foodController.BranchUpdateFoodRecipeColumns, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/recipe/:id", foodController.BranchDeleteFoodRecipe, middleware.Authenticate(r.auth, auth.RoleBranch))

		//	#waiter-work-time
		r.Get("/api/v1/branch/waiter/work-time/list", userController.BranchGetListWaiterWorkTime, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/waiter/work-time/:waiter_id", userController.BranchGetDetailWaiterWorkTime, middleware.Authenticate(r.auth, auth.RoleBranch))

		//	 #measure_unit
		r.Get("/api/v1/branch/measure-unit/list", catalogController.BranchGetMeasureUnitList, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #category
		r.Get("/api/v1/branch/category/list", foodController.BranchGetCategoryList, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #cashier
		r.Get("/api/v1/branch/cashier/list", userController.BranchGetCashierList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/cashier/:id", userController.BranchGetCashierDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/cashier/create", userController.BranchCreateCashier, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/cashier/:id", userController.BranchUpdateCashierAll, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/cashier/:id", userController.BranchUpdateCashierColumns, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/cashier/:id", userController.BranchDeleteCashier, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/cashier/status/:id", userController.BranchUpdateCashierStatus, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/cashier/password/:id", userController.BranchUpdateCashierPassword, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/cashier/phone/:id", userController.BranchUpdateCashierPhone, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #menu_category
		r.Get("/api/v1/branch/menu-category/list", foodController.BranchGetMenuCategoryList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/menu-category/:id", foodController.BranchGetMenuCategoryDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/menu-category/create", foodController.BranchCreateMenuCategory, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/menu-category/:id", foodController.BranchUpdateMenuCategory, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/menu-category/:id", foodController.BranchDeleteMenuCategory, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #food_category
		r.Get("/api/v1/branch/food-category/list", foodController.BranchGetFoodCategoryList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/food-category/:id", foodController.BranchGetFoodCategoryDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/food-category/create", foodController.BranchCreateFoodCategory, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/food-category/:id", foodController.BranchUpdateFoodCategoryAll, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/food-category/:id", foodController.BranchUpdateFoodCategoryColumns, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/food-category/:id", foodController.BranchDeleteFoodCategory, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #food_recipe_group
		r.Get("/api/v1/branch/recipe/group/list", foodController.BranchGetFoodRecipeGroupList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/recipe/group/:id", foodController.BranchGetFoodRecipeGroupDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/recipe/group/create", foodController.BranchCreateFoodRecipeGroup, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/recipe/group/:id", foodController.BranchUpdateFoodRecipeGroupAll, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/recipe/group/:id", foodController.BranchUpdateFoodRecipeGroupColumns, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/recipe/group/single/:id", foodController.BranchDeleteFoodRecipeGroupSingleRecipe, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/recipe/group/:id", foodController.BranchDeleteFoodRecipeGroup, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #food_recipe_group_history
		r.Get("/api/v1/branch/recipe/group/history/list", foodController.BranchGetFoodRecipeGroupHistoryList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/recipe/group/history/:id", foodController.BranchGetFoodRecipeGroupHistoryDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/recipe/group/history/create", foodController.BranchCreateFoodRecipeGroupHistory, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/recipe/group/history/:id", foodController.BranchDeleteFoodRecipeGroupHistory, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #hall
		r.Get("/api/v1/branch/hall/list", restaurantController.BranchGetHallList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/hall/:id", restaurantController.BranchGetHallDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/hall/create", restaurantController.BranchCreateHall, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/hall/:id", restaurantController.BranchUpdateHallAll, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/hall/:id", restaurantController.BranchUpdateHallColumns, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/hall/:id", restaurantController.BranchDeleteHall, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #product_recipe_group
		r.Get("/api/v1/branch/product/recipe/group/list", productController.BranchGetProductRecipeGroupList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/product/recipe/group/:id", productController.BranchGetProductRecipeGroupDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/product/recipe/group/create", productController.BranchCreateProductRecipeGroup, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Put("/api/v1/branch/product/recipe/group/:id", productController.BranchUpdateProductRecipeGroupAll, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/product/recipe/group/:id", productController.BranchUpdateProductRecipeGroupColumns, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Patch("/api/v1/branch/product/recipe/group/single/:id", productController.BranchDeleteProductRecipeGroupSingleRecipe, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/product/recipe/group/:id", productController.BranchDeleteProductRecipeGroup, middleware.Authenticate(r.auth, auth.RoleBranch))

		// #product_recipe_group_history
		r.Get("/api/v1/branch/product/recipe/group/history/list", productController.BranchGetProductRecipeGroupHistoryList, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Get("/api/v1/branch/product/recipe/group/history/:id", productController.BranchGetProductRecipeGroupHistoryDetail, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Post("/api/v1/branch/product/recipe/group/history/create", productController.BranchCreateProductRecipeGroupHistory, middleware.Authenticate(r.auth, auth.RoleBranch))
		r.Delete("/api/v1/branch/product/recipe/group/history/:id", productController.BranchDeleteProductRecipeGroupHistory, middleware.Authenticate(r.auth, auth.RoleBranch))

	}

	// @cashier
	{
		// order
		r.Put("/api/v1/cashier/order/payment/:id", orderController.CashierPaymentOrder, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Get("/api/v1/cashier/order/list", orderController.CashierGetOrderList, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Get("/api/v1/cashier/order/:id", orderController.CashierGetOrderDetail, middleware.Authenticate(r.auth, auth.RoleCashier))

		//	#waiter
		r.Get("/api/v1/cashier/waiter/work-time/list", userController.CashierGetListWaiterWorkTime, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Get("/api/v1/cashier/waiter/work-time/:waiter_id", userController.CashierGetDetailWaiterWorkTime, middleware.Authenticate(r.auth, auth.RoleCashier))
		//r.Get("/api/v1/cashier/waiter/list", userController.CashierGetWaiterList, middleware.Authenticate(r.auth, auth.RoleCashier))

		// #menu
		r.Get("/api/v1/cashier/menu/list", foodController.CashierGetMenuList, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Get("/api/v1/cashier/menu/:id", foodController.CashierGetMenuDetail, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Post("/api/v1/cashier/menu/create", foodController.CashierCreateMenu, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Put("/api/v1/cashier/menu/:id", foodController.CashierUpdateMenuAll, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Patch("/api/v1/cashier/menu/:id", foodController.CashierUpdateMenuColumn, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Delete("/api/v1/cashier/menu/:id", foodController.CashierDeleteMenu, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Patch("/api/v1/cashier/menu/add/printer", foodController.CashierUpdateMenuPrinterID, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Delete("/api/v1/cashier/menu/printer/:id", foodController.CashierDeleteMenuPrinterID, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Patch("/api/v1/cashier/menu/remove/photo/:id", foodController.CashierRemoveMenuPhoto, middleware.Authenticate(r.auth, auth.RoleCashier))

		// #food
		r.Get("/api/v1/cashier/food/list", foodController.CashierGetFoodList, middleware.Authenticate(r.auth, auth.RoleCashier))

		r.Get("/api/v1/cashier/me", userController.CashierGetMe, middleware.Authenticate(r.auth, auth.RoleCashier))

		//	#order_menu
		r.Patch("/api/v1/cashier/order/menu/status", orderController.CashierUpdateOrderMenuStatus, middleware.Authenticate(r.auth, auth.RoleCashier))

		// #waiters
		r.Get("/api/v1/cashier/waiter/list", userController.CashierGetWaiterLists, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Get("/api/v1/cashier/waiter/:id", userController.CashierGetWaiterDetails, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Post("/api/v1/cashier/waiter/create", userController.CashierCreateWaiter, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Put("/api/v1/cashier/waiter/:id", userController.CashierUpdateWaiterAll, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Patch("/api/v1/cashier/waiter/:id", userController.CashierUpdateWaiterColumns, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Delete("/api/v1/cashier/waiter/:id", userController.CashierDeleteWaiter, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Patch("/api/v1/cashier/waiter/status/:id", userController.CashierUpdateWaiterStatus, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Patch("/api/v1/cashier/waiter/password/:id", userController.CashierUpdateWaiterPassword, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Patch("/api/v1/cashier/waiter/phone/:id", userController.CashierUpdateWaiterPhone, middleware.Authenticate(r.auth, auth.RoleCashier))

		// #menu_category
		r.Get("/api/v1/cashier/menu-category/list", foodController.CashierGetMenuCategoryList, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Get("/api/v1/cashier/menu-category/:id", foodController.CashierGetMenuCategoryDetail, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Post("/api/v1/cashier/menu-category/create", foodController.CashierCreateMenuCategory, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Patch("/api/v1/cashier/menu-category/:id", foodController.CashierUpdateMenuCategory, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Delete("/api/v1/cashier/menu-category/:id", foodController.CashierDeleteMenuCategory, middleware.Authenticate(r.auth, auth.RoleCashier))

		// #food_category
		r.Get("/api/v1/cashier/food-category/list", foodController.CashierGetFoodCategoryList, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Get("/api/v1/cashier/food-category/:id", foodController.CashierGetFoodCategoryDetail, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Post("/api/v1/cashier/food-category/create", foodController.CashierCreateFoodCategory, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Put("/api/v1/cashier/food-category/:id", foodController.CashierUpdateFoodCategoryAll, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Patch("/api/v1/cashier/food-category/:id", foodController.CashierUpdateFoodCategoryColumns, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Delete("/api/v1/cashier/food-category/:id", foodController.CashierDeleteFoodCategory, middleware.Authenticate(r.auth, auth.RoleCashier))

		// #category
		r.Get("/api/v1/cashier/category/list", foodController.CashierGetCategoryList, middleware.Authenticate(r.auth, auth.RoleCashier))

		// #table
		r.Get("/api/v1/cashier/table/list", restaurantController.CashierGetTableList, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Get("/api/v1/cashier/table/:id", restaurantController.CashierGetTableDetail, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Post("/api/v1/cashier/table/create", restaurantController.CashierCreateTable, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Put("/api/v1/cashier/table/:id", restaurantController.CashierUpdateTableAll, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Patch("/api/v1/cashier/table/:id", restaurantController.CashierUpdateTableColumns, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Delete("/api/v1/cashier/table/:id", restaurantController.CashierDeleteTable, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Post("/api/v1/cashier/generate/qr", restaurantController.CashierGenerateQRTable, middleware.Authenticate(r.auth, auth.RoleCashier))

		// #product
		r.Get("/api/v1/cashier/product/list", productController.CashierGetProductList, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Get("/api/v1/cashier/product/:id", productController.CashierGetProductDetail, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Post("/api/v1/cashier/product/create", productController.CashierCreateProduct, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Put("/api/v1/cashier/product/:id", productController.CashierUpdateProductAll, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Patch("/api/v1/cashier/product/:id", productController.CashierUpdateProductColumns, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Delete("/api/v1/cashier/product/:id", productController.CashierDeleteProduct, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Get("/api/v1/cashier/product/spending", productController.CashierGetProductSpending, middleware.Authenticate(r.auth, auth.RoleCashier))

		// #hall
		r.Get("/api/v1/cashier/hall/list", restaurantController.CashierGetHallList, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Get("/api/v1/cashier/hall/:id", restaurantController.CashierGetHallDetail, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Post("/api/v1/cashier/hall/create", restaurantController.CashierCreateHall, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Put("/api/v1/cashier/hall/:id", restaurantController.CashierUpdateHallAll, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Patch("/api/v1/cashier/hall/:id", restaurantController.CashierUpdateHallColumns, middleware.Authenticate(r.auth, auth.RoleCashier))
		r.Delete("/api/v1/cashier/hall/:id", restaurantController.CashierDeleteHall, middleware.Authenticate(r.auth, auth.RoleCashier))

		// #order_report
		r.Post("/api/v1/cashier/order-report/report", orderController.CashierReportOrder, middleware.Authenticate(r.auth, auth.RoleCashier))
	}

	// @client
	{
		// #order-food
		//r.Get("/api/v1/client/order-food/list", orderController.ClientGetOrderMenuList, middleware.Authenticate(r.auth, auth.RoleClient))
		//r.Get("/api/v1/client/order-food/:id", orderController.ClientGetOrderMenuDetail, middleware.Authenticate(r.auth, auth.RoleClient))
		//r.Post("/api/v1/client/order-food/create", orderController.ClientCreateOrderFood, middleware.Authenticate(r.auth, auth.RoleClient))
		//r.Put("/api/v1/client/order-food/:id", orderController.ClientUpdateOrderFoodAll, middleware.Authenticate(r.auth, auth.RoleClient))
		//r.Patch("/api/v1/client/order-food/:id", orderController.ClientUpdateOrderFoodColumns, middleware.Authenticate(r.auth, auth.RoleClient))
		//r.AdminDelete("/api/v1/client/order-food/:id", orderController.ClientDeleteOrderFood, middleware.Authenticate(r.auth, auth.RoleClient))

		// #order-payment
		r.Get("/api/v1/client/order-payment/list", orderController.ClientGetOrderPaymentList, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Get("/api/v1/client/order-payment/:id", orderController.ClientGetOrderPaymentDetail, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Post("/api/v1/client/order-payment/create", orderController.ClientCreateOrderPayment, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Put("/api/v1/client/order-payment/:id", orderController.ClientUpdateOrderPaymentAll, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Patch("/api/v1/client/order-payment/:id", orderController.ClientUpdateOrderPaymentColumns, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Delete("/api/v1/client/order-payment/:id", orderController.ClientDeleteOrderPayment, middleware.Authenticate(r.auth, auth.RoleClient))

		// #client
		r.Get("/api/v1/client/me", userController.ClientGetUserMe, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Patch("/api/v1/client/me", userController.ClientUpdateUserColumns, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Delete("/api/v1/client/me", userController.ClientDeleteUserMe, middleware.Authenticate(r.auth, auth.RoleClient))

		// #feedback
		r.Get("/api/v1/client/feedback/list", catalogController.ClientGetFeedBackList, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Post("/api/v1/client/post-feedback", userController.SendFeedback)

		// #basket
		r.Get("/api/v1/client/basket/:branch-id", catalogController.ClientGetBasket, middleware.Authenticate(r.auth, auth.RoleClient))
		//r.Post("/api/v1/client/basket/create", catalogController.ClientCreateBasket, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Put("/api/v1/client/basket", catalogController.ClientUpdateBasket, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Delete("/api/v1/client/basket/:branch-id", catalogController.ClientDeleteBasket, middleware.Authenticate(r.auth, auth.RoleClient))

		// #branch
		r.Get("/api/v1/mobile/client/branch/list", restaurantController.ClientGetBranchList, middleware.Authenticate(r.auth, auth.RoleClient))                          // changed
		r.Get("/api/v1/mobile/client/branch/map/list", restaurantController.ClientGetMapBranchList, middleware.Authenticate(r.auth, auth.RoleClient))                   // changed
		r.Get("/api/v1/mobile/client/branch/:id", restaurantController.ClientGetBranchDetail, middleware.Authenticate(r.auth, auth.RoleClient))                         // changed
		r.Get("/api/v1/mobile/client/nearly/branch/list", restaurantController.ClientGetNearlyBranchList, middleware.Authenticate(r.auth, auth.RoleClient))             // changed
		r.Patch("/api/v1/mobile/client/branch/:id", restaurantController.ClientUpdateBranchColumns, middleware.Authenticate(r.auth, auth.RoleClient))                   // changed
		r.Patch("/api/v1/mobile/client/branch/add/search-count/:id", restaurantController.ClientAddBranchSearchCount, middleware.Authenticate(r.auth, auth.RoleClient)) // changed
		r.Get("/api/v1/mobile/client/branch/list/order/search-count", restaurantController.ClientGetBranchListOrderSearchCount, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Get("/api/v1/mobile/client/branch/list/by/category/:category_id", restaurantController.ClientGetBranchListByCategoryID, middleware.Authenticate(r.auth, auth.RoleClient)) // changed

		// #food
		r.Get("/api/v1/mobile/client/food-category/list", foodController.ClientGetFoodCategoryList, middleware.Authenticate(r.auth, auth.RoleClient))                     // changed
		r.Get("/api/v1/mobile/client/menu/list", foodController.ClientGetMenuList, middleware.Authenticate(r.auth, auth.RoleClient))                                      // changed
		r.Get("/api/v1/mobile/client/menu/list/by/category/:category_id", foodController.ClientGetMenuListByCategoryID, middleware.Authenticate(r.auth, auth.RoleClient)) // changed

		// #order
		r.Get("/api/v1/mobile/client/order/list", orderController.MobileGetOrderList, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Get("/api/v1/mobile/client/order/:id", orderController.MobileGetOrderDetail, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Post("/api/v1/mobile/client/order/create", orderController.MobileCreateOrder, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Patch("/api/v1/mobile/client/order/update/:id", orderController.MobileUpdateOrder, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Get("/api/v1/mobile/client/order/menu/often/by/branch/:branch_id", orderController.MobileClientGetOrderMenuOftenByBranchID, middleware.Authenticate(r.auth, auth.RoleClient)) // changed
		r.Post("/api/v1/mobile/client/order/waiter/call/:id", orderController.MobileWaiterCall, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Post("/api/v1/mobile/client/order/review", orderController.ClientReviewOrder, middleware.Authenticate(r.auth, auth.RoleClient))

		// #branch-review
		r.Get("/api/v1/client/branch-review/list", restaurantController.ClientGetBranchReviewList, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Get("/api/v1/client/branch-review/:id", restaurantController.ClientGetBranchReviewDetail, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Post("/api/v1/client/branch-review/create", restaurantController.ClientCreateBranchReview, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Put("/api/v1/client/branch-review/:id", restaurantController.ClientUpdateBranchReviewAll, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Patch("/api/v1/client/branch-review/:id", restaurantController.ClientUpdateBranchReviewColumns, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Delete("/api/v1/client/branch-review/:id", restaurantController.ClientDeleteBranchReview, middleware.Authenticate(r.auth, auth.RoleClient))

		// #story
		r.Get("/api/v1/client/story/list", catalogController.ClientGetStoryList, middleware.Authenticate(r.auth, auth.RoleClient))            // changed
		r.Post("/api/v1/client/story/viewed/:id", catalogController.ClientSetStoryAsViewed, middleware.Authenticate(r.auth, auth.RoleClient)) // changed

		// #banner
		r.Get("/api/v1/client/banner/list", catalogController.ClientGetBannerList, middleware.Authenticate(r.auth, auth.RoleClient))  // changed
		r.Get("/api/v1/client/banner/:id", catalogController.ClientGetBannerDetail, middleware.Authenticate(r.auth, auth.RoleClient)) // changed

		// #full_search
		r.Get("/api/v1/client/full-search/list", catalogController.ClientGetFullSearchList, middleware.Authenticate(r.auth, auth.RoleClient)) // changed

		// #notification
		r.Get("/api/v1/client/notification/list", catalogController.ClientGetNotificationList, middleware.Authenticate(r.auth, auth.RoleClient))                // changed
		r.Get("/api/v1/client/notification/unseen/count", catalogController.ClientGetCountUnseenNotification, middleware.Authenticate(r.auth, auth.RoleClient)) // changed
		r.Post("/api/v1/client/notification/viewed/:id", catalogController.ClientSetNotificationAsViewed, middleware.Authenticate(r.auth, auth.RoleClient))     // changed
		r.Post("/api/v1/client/notification/viewed/all", catalogController.ClientSetAllNotificationsAsViewed, middleware.Authenticate(r.auth, auth.RoleClient))

		// #menu_category
		r.Get("/api/v1/client/category/list", foodController.ClientGetCategoryList, middleware.Authenticate(r.auth, auth.RoleClient))

	}

	// @waiter
	{
		// #auth
		r.Get("/api/v1/waiter/me", userController.WaiterGetMe, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Get("/api/v1/waiter/personal/info", userController.WaiterGetPersonalInfo, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Patch("/api/v1/waiter/photo/update", userController.WaiterUpdatePhoto, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Patch("/api/v1/waiter/password/update", authController.WaiterUpdatePassword, middleware.Authenticate(r.auth, auth.RoleWaiter))

		// #order
		r.Get("/api/v1/waiter/order/list", orderController.WaiterGetOrderList, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Post("/api/v1/waiter/order/create", orderController.WaiterCreateOrder, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Get("/api/v1/waiter/order/:id", orderController.WaiterGetOrderDetail, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Patch("/api/v1/waiter/order/:id", orderController.WaiterUpdateOrder, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Patch("/api/v1/waiter/order/status/:id", orderController.WaiterUpdateOrderStatus, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Patch("/api/v1/waiter/order/accept/:id", orderController.WaiterAcceptOrder, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Patch("/api/v1/waiter/order/menu/status", orderController.WaiterUpdateOrderMenuStatus, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Get("/api/v1/waiter/order/history/list", orderController.WaiterHistoryActivityList, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Get("/api/v1/waiter/order/history/:id", orderController.WaiterGetMyOrderDetail, middleware.Authenticate(r.auth, auth.RoleWaiter))

		// #table
		r.Get("/api/v1/waiter/table/list", restaurantController.WaiterGetTableList, middleware.Authenticate(r.auth, auth.RoleWaiter))

		// #menu
		r.Get("/api/v1/waiter/menu/list", foodController.WaiterGetMenuList, middleware.Authenticate(r.auth, auth.RoleWaiter))

		// #food_category
		r.Get("/api/v1/waiter/food-category/list", foodController.WaiterGetFoodCategoryList, middleware.Authenticate(r.auth, auth.RoleWaiter))

		// #attendance
		r.Post("/api/v1/waiter/attendance/create/come", userController.WaiterCreateComeAttendance, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Post("/api/v1/waiter/attendance/create/gone", userController.WaiterCreateGoneAttendance, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Get("/api/v1/waiter/my/work/history", userController.WaiterGetListWorkTime, middleware.Authenticate(r.auth, auth.RoleWaiter))

		// #menu_category
		r.Get("/api/v1/waiter/category/list", foodController.WaiterGetCategoryList, middleware.Authenticate(r.auth, auth.RoleWaiter))

		//	statistics
		r.Get("/api/v1/waiter/rating/activity", userController.WaiterGetActivityStatistics, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Get("/api/v1/waiter/earned/statistics", userController.WaiterGetWeeklyActivityStatistics, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Get("/api/v1/waiter/accepted/statistics", userController.WaiterGetWeeklyAcceptedOrdersStatistics, middleware.Authenticate(r.auth, auth.RoleWaiter))
		r.Get("/api/v1/waiter/rating/activity/weekly", userController.WaiterGetWeeklyRatingStatistics, middleware.Authenticate(r.auth, auth.RoleWaiter))

		// #hall
		r.Get("/api/v1/waiter/hall/list", restaurantController.WaiterGetHallList, middleware.Authenticate(r.auth, auth.RoleWaiter))
	}

	// @site [landing page]
	{
		r.Get("/api/v1/site/restaurant/list", restaurantController.SiteGetRestaurantList)
		r.Get("/api/v1/site/restaurant/category/list", restaurantController.SiteGetRestaurantCategoryList)
		r.Get("/api/v1/getme", userController.GetMe)
	}

	return r.Run(r.port)
}
