package router

import (
	"github.com/gin-gonic/gin"

	"github.com/hyperjiang/gin-skeleton/controller"
	"github.com/hyperjiang/gin-skeleton/manager"
)

// Route makes the routing
func Route(app *gin.Engine) {
	inMemoryStoragemanager := manager.NewInMemoryStorageManager()

	indexController := new(controller.IndexController)
	userController := controller.NewUserController()
	avatarController := controller.NewAvatarController()
	nftController := controller.NewNftController(inMemoryStoragemanager)
	eventController := controller.NewEventController()
	subworldTemplateController := controller.NewSubworldTemplateController()
	subworldInstanceController := controller.NewSubworldInstanceController()

	app.GET(
		"/", indexController.GetIndex,
	)

	api := app.Group("/api")
	{
		api.GET("/version", indexController.GetVersion)

		api.GET("/user/profile", userController.GetUserPrivateProfile)
		api.GET("/user/profile/:id/getAvatars", userController.GetAvatars)
		//api.GET("/user/profile/:id", userController.GetUserPublicProfile)
		api.GET("/user/getByWallet/:wallet_address", userController.GetUserByWalletAddress)
		//api.POST("/wallet/updateAssets", userController.UpdateAssets)
		//api.GET("/wallet/fetchAssets/:address", userController.FetchAssets)
		api.POST("/user/getOrCreate", userController.GetOrCreate)
		api.POST("/user/auth", userController.Auth)
		api.POST("/user/mockAuth", userController.MockAuth)
		api.POST("/user/createLoginLink", userController.CreateLoginLink)
		api.POST("/user/authLoginLink", userController.AuthLoginLink)
		api.GET("/user/pollLoginLink/:session_key", userController.PollLoginLink)
		api.GET("/user/getTemporaryEventRewards", userController.GetTemporaryEventRewards)

		api.GET("/avatar/:id", avatarController.Get)
		api.POST("/avatar", avatarController.Create)
		api.PUT("/avatar/:id", avatarController.Update)
		api.DELETE("/avatar/:id", avatarController.Delete)

		api.POST("/nft/createMintNftLink", nftController.CreateMintNftLink)
		api.POST("/nft/notifyMinted", nftController.NotifyMinted)
		api.GET("/nft/checkName", nftController.CheckName)
		api.POST("/nft/lockName", nftController.LockName)
		api.POST("/nft/unlockName", nftController.UnlockName)

		api.POST("/event", eventController.CreateEvent)
		api.POST("/event/:id/start", eventController.StartEvent)
		api.POST("/event/:id/stop", eventController.StartEvent)
		api.POST("/event/:id/join", eventController.JoinEvent)

		api.GET("/subworld/root_template", subworldTemplateController.GetAllRoot)
		api.POST("/subworld/root_template", subworldTemplateController.CreateRoot)
		api.PUT("/subworld/root_template/:root_id", subworldTemplateController.UpdateRoot)
		api.DELETE("/subworld/root_template/:root_id", subworldTemplateController.DeleteRoot)

		api.GET("/subworld/root_template/:root_id/deriv", subworldTemplateController.GetAllDeriv)
		api.POST("/subworld/root_template/:root_id/deriv", subworldTemplateController.CreateDeriv)
		api.PUT("/subworld/root_template/:root_id/deriv/:deriv_id", subworldTemplateController.UpdateDeriv)
		api.DELETE("/subworld/root_template/:root_id/deriv/:deriv_id", subworldTemplateController.DeleteDeriv)

		api.GET("/subworld/instance", subworldInstanceController.GetAll)
		api.POST("/subworld/instance", subworldInstanceController.Create)
		api.PUT("/subworld/instance/:id", subworldInstanceController.Update)
		api.DELETE("/subworld/instance/:id", subworldInstanceController.Delete)
	}

}
