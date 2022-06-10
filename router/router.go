package router

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/hyperjiang/gin-skeleton/controller"
	"github.com/hyperjiang/gin-skeleton/manager"
	"github.com/hyperjiang/gin-skeleton/middleware"
	"github.com/hyperjiang/gin-skeleton/model"
)

// Route makes the routing
func Route(app *gin.Engine) {
	inMemoryStoragemanager := manager.NewInMemoryStorageManager()

	indexController := new(controller.IndexController)
	userController := new(controller.UserController)
	walletController := controller.NewWalletController()
	nftController := controller.NewNftController(inMemoryStoragemanager)
	eventController := controller.NewEventController()

	app.GET(
		"/", indexController.GetIndex,
	)

	auth := app.Group("/auth")
	authMiddleware := middleware.Auth()
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", func(c *gin.Context) {
			claims := jwt.ExtractClaims(c)
			user, _ := c.Get("email")
			c.JSON(200, gin.H{
				"email": claims["email"],
				"name":  user.(*model.User).Name,
				"text":  "Hello World.",
			})
		})
	}

	// #TODO: Move to /api
	app.GET(
		"/user/:id", userController.GetUser,
	).GET(
		"/signup", userController.SignupForm,
	).POST(
		"/signup", userController.Signup,
	).GET(
		"/login", userController.LoginForm,
	).POST(
		"/login", authMiddleware.LoginHandler,
	)

	api := app.Group("/api")
	{
		api.GET("/version", indexController.GetVersion)

		api.GET("/wallet/get/:address", walletController.GetWallet)
		api.GET("/wallet/profile", walletController.GetWalletPrivateProfile)
		//api.POST("/wallet/updateAssets", walletController.UpdateAssets)
		//api.GET("/wallet/fetchAssets/:address", walletController.FetchAssets)
		api.POST("/wallet/getOrCreate", walletController.GetOrCreateWallet)
		api.POST("/wallet/auth", walletController.Auth)
		api.POST("/wallet/mockAuth", walletController.MockAuth)
		api.POST("/wallet/createLoginLink", walletController.CreateLoginLink)
		api.POST("/wallet/authLoginLink", walletController.AuthLoginLink)
		api.GET("/wallet/pollLoginLink/:session_key", walletController.PollLoginLink)
		api.GET("/wallet/getTemporaryEventRewards", walletController.GetTemporaryEventRewards)

		api.POST("/nft/createMintNftLink", nftController.CreateMintNftLink)
		api.POST("/nft/notifyMinted", nftController.NotifyMinted)
		api.GET("/nft/checkName", nftController.CheckName)
		api.POST("/nft/lockName", nftController.LockName)
		api.POST("/nft/unlockName", nftController.UnlockName)

		api.POST("/event", eventController.CreateEvent)
		api.POST("/event/:id/start", eventController.StartEvent)
		api.POST("/event/:id/stop", eventController.StartEvent)
		api.POST("/event/:id/join", eventController.JoinEvent)
	}

	api.Use(authMiddleware.MiddlewareFunc())

}
