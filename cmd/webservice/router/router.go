package router

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	ecMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-account/cmd/webservice/handler"
	"github.com/stellar-payment/sp-account/internal/config"
	"github.com/stellar-payment/sp-account/internal/middleware"
	"github.com/stellar-payment/sp-account/internal/service"
)

type InitRouterParams struct {
	Logger  zerolog.Logger
	Service service.Service
	Ec      *echo.Echo
	Conf    *config.Config
}

func Init(params *InitRouterParams) {
	params.Ec.Use(
		ecMiddleware.CORS(), ecMiddleware.RequestIDWithConfig(ecMiddleware.RequestIDConfig{Generator: uuid.NewString}),
		middleware.ServiceVersioner,
		middleware.RequestBodyLogger(&params.Logger),
		middleware.RequestLogger(&params.Logger),
		middleware.HandlerLogger(&params.Logger),
	)

	plainRouter := params.Ec.Group("")
	secureRouter := params.Ec.Group("", middleware.AuthorizationMiddleware(params.Service))

	// ----- Maintenance
	plainRouter.GET(PingPath, handler.HandlePing(params.Service.Ping))

	// ----- Register
	plainRouter.POST(registerCustomerPath, handler.HandleRegisterCustomer(params.Service.RegisterCustomer))
	plainRouter.OPTIONS(registerCustomerPath, handler.HandleRegisterCustomer(params.Service.RegisterCustomer))
	plainRouter.POST(registerMerchantPath, handler.HandleRegisterMerchant(params.Service.RegisterMerchant))
	plainRouter.OPTIONS(registerMerchantPath, handler.HandleRegisterMerchant(params.Service.RegisterMerchant))

	// ----- Auth
	plainRouter.POST(authSignupPath, handler.HandleSignup(params.Service.AuthSignup))
	plainRouter.OPTIONS(authSignupPath, handler.HandleSignup(params.Service.AuthSignup))
	plainRouter.POST(authLoginPath, handler.HandleAuthLogin(params.Service.AuthLogin))
	plainRouter.OPTIONS(authLoginPath, handler.HandleAuthLogin(params.Service.AuthLogin))

	// ----- Users
	secureRouter.GET(userBasepath, handler.HandleGetUsers(params.Service.GetAllUser))
	secureRouter.OPTIONS(userBasepath, handler.HandleGetUsers(params.Service.GetAllUser))
	secureRouter.GET(userMePath, handler.HandleGetUserMe(params.Service.GetUserMe))
	secureRouter.OPTIONS(userMePath, handler.HandleGetUserMe(params.Service.GetUserMe))
	secureRouter.GET(userIDPath, handler.HandleGetUserByID(params.Service.GetUser))
	secureRouter.OPTIONS(userIDPath, handler.HandleGetUserByID(params.Service.GetUser))
	secureRouter.POST(userBasepath, handler.HandleCreateUsers(params.Service.CreateUser))
	secureRouter.OPTIONS(userBasepath, handler.HandleCreateUsers(params.Service.CreateUser))
	secureRouter.PUT(userIDPath, handler.HandleUpdateUsers(params.Service.UpdateUser))
	secureRouter.OPTIONS(userIDPath, handler.HandleUpdateUsers(params.Service.UpdateUser))
	secureRouter.DELETE(userIDPath, handler.HandleDeleteUser(params.Service.DeleteUser))
	secureRouter.OPTIONS(userIDPath, handler.HandleDeleteUser(params.Service.DeleteUser))

}
