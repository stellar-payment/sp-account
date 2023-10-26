package router

const (
	basePath = "/account/api/v1"
	PingPath = basePath + "/ping"

	// ----- Auth
	authLogin  = basePath + "/auth/login"
	authSignup = basePath + "/auth/signup"

	// ----- Users
	userMePath   = basePath + "/me"
	userBasepath = basePath + "/users"
	userIDPath   = userBasepath + "/:userID"
)
