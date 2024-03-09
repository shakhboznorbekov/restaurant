package router

import (
	"github.com/redis/go-redis/v9"
	"github.com/restaurant/foundation/web"
	"github.com/restaurant/internal/auth"
	restaurant_controller "github.com/restaurant/internal/controller/http/v1/restaurant"
	user4 "github.com/restaurant/internal/controller/http/v1/user"
	"github.com/restaurant/internal/middleware"
	"github.com/restaurant/internal/pkg/repository/postgresql"
	"github.com/restaurant/internal/repository/postgres/device"
	"github.com/restaurant/internal/repository/postgres/restaurant"
	"github.com/restaurant/internal/repository/postgres/user"
	"github.com/restaurant/internal/repository/redis/hashing"
	device2 "github.com/restaurant/internal/service/device"
	"github.com/restaurant/internal/service/fcm"
	hashing2 "github.com/restaurant/internal/service/hashing"
	"github.com/restaurant/internal/service/sms"
	"github.com/restaurant/internal/socket"

	auth3 "github.com/restaurant/internal/controller/http/v1/auth"
	restaurant2 "github.com/restaurant/internal/service/restaurant"
	restaurantCategory2 "github.com/restaurant/internal/service/restaurant_category"
	user2 "github.com/restaurant/internal/service/user"
	auth2 "github.com/restaurant/internal/usecase/auth"
	restaurant3 "github.com/restaurant/internal/usecase/restaurant"
	user3 "github.com/restaurant/internal/usecase/user"

	restaurantCategory "github.com/restaurant/internal/repository/postgres/restaurant_category"
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
	//repositories:
	//-postgresql
	userPostgres := user.NewRepository(r.postgresDB)
	devicePostgres := device.NewRepository(r.postgresDB)
	restaurantPostgres := restaurant.NewRepository(r.postgresDB)
	restaurantCategoryPostgres := restaurantCategory.NewRepository(r.postgresDB)

	// - redis
	hashingRedis := hashing.NewRepository(r.redisDB)
	//branchRedis := redisBranch.NewRepository(r.redisDB)
	//basketRedis := basket.NewRepository(r.redisDB)
	hashing2.Hashing = hashingRedis

	//-service
	userService := user2.NewService(userPostgres)
	smsService := sms.NewService(r.smsConfig, r.redisDB, r.postgresDB)
	deviceService := device2.NewService(devicePostgres)
	fcmService := fcm.NewFCMService(r.fcm)
	restaurantService := restaurant2.NewService(restaurantPostgres)
	restaurantCategoryService := restaurantCategory2.NewService(restaurantCategoryPostgres)

	// cron service for scheduled tasks
	schedule := cron.New()

	//-use_case
	userUseCase := user3.NewUseCase(userService, smsService, r.auth, schedule)
	authUseCase := auth2.NewUseCase(userService, smsService, deviceService, r.auth)
	restaurantUseCase := restaurant3.NewUseCase(
		restaurantService,
		userService,
		restaurantCategoryService,
	)
	//-controller
	userController := user4.NewController(userUseCase)
	authController := auth3.NewController(authUseCase)
	restaurantController := restaurant_controller.NewController(restaurantUseCase)
	fileController := file.NewController(r.App)

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
		r.Post("/api/v1/sign-in", authController.SignIn)
		r.Post("/api/v1/client/fill-up", authController.ClientFillUp, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Put("/api/v1/client/update/phone", authController.ClientUpdatePhone, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Delete("/api/v1/sign-out/:device-id", authController.ClientLogOut, middleware.Authenticate(r.auth, auth.RoleClient))
		r.Post("/api/v1/change/lang", authController.ChangeDeviceLang, middleware.Authenticate(r.auth, auth.RoleClient))

		// admin - super_admin - cashier
		r.Post("/api/v1/admin/sign-in", authController.AdminSignIn)

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
		r.Delete("/api/v1/super-admin/restaurant/:id", restaurantController.SuperAdminDeleteRestaurant, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))

		// #restaurant_category
		r.Get("/api/v1/super-admin/restaurant-category/list", restaurantController.SuperAdminGetRestaurantCategoryList, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Get("/api/v1/super-admin/restaurant-category/:id", restaurantController.SuperAdminGetRestaurantCategoryDetail, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Post("/api/v1/super-admin/restaurant-category/create", restaurantController.SuperAdminCreateRestaurantCategory, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Put("/api/v1/super-admin/restaurant-category/:id", restaurantController.SuperAdminUpdateRestaurantCategoryAll, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Patch("/api/v1/super-admin/restaurant-category/:id", restaurantController.SuperAdminUpdateRestaurantCategoryColumns, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))
		r.Delete("/api/v1/super-admin/restaurant-category/:id", restaurantController.SuperAdminDeleteRestaurantCategory, middleware.Authenticate(r.auth, auth.RoleSuperAdmin))

	}
	// @admin
	{
		// #restaurant_category
		//r.Get("/api/v1/admin/restaurant-category/list", restaurantController.AdminGetRestaurantCategoryList, middleware.Authenticate(r.auth, auth.RoleAdmin))
	}

	// @site [landing page]
	{
		r.Get("/api/v1/site/restaurant/list", restaurantController.SiteGetRestaurantList)
		r.Get("/api/v1/site/restaurant/category/list", restaurantController.SiteGetRestaurantCategoryList)
		r.Get("/api/v1/getme", userController.GetMe)
	}

	return r.Run(r.port)
}
