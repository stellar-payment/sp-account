package router

const (
	basePath = "/account/api/v1"
	PingPath = basePath + "/ping"

	// ----- Register
	registerCustomerPath = basePath + "/register/customer"
	registerMerchantPath = basePath + "/register/merchant"

	// ----- Auth
	authLoginPath  = basePath + "/auth/login"
	authSignupPath = basePath + "/auth/signup"

	// ----- Users
	userMePath   = basePath + "/me"
	userBasepath = basePath + "/users"
	userIDPath   = userBasepath + "/:userID"
)
