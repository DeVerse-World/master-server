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
	systemSettingController := controller.NewSystemSettingController()

	app.GET(
		"/", indexController.GetIndex,
	)

	api := app.Group("/api")
	{
		api.GET("/version", indexController.GetVersion)

		api.PUT("/user/:id", userController.Update)
		api.GET("/user/profile", userController.GetUserPrivateProfile)
		api.POST("/user/profile", userController.UpdateUserProfile)
		api.GET("/user/profile/:id/getAvatars", userController.GetAvatars)
		//api.GET("/user/profile/:id", userController.GetUserPublicProfile)
		api.GET("/user/getByWallet/:wallet_address", userController.GetUserByWalletAddress)
		//api.POST("/wallet/updateAssets", userController.UpdateAssets)
		//api.GET("/wallet/fetchAssets/:address", userController.FetchAssets)
		api.POST("/user/getOrCreate", userController.GetOrCreate)
		api.POST("/user/mockAuth", userController.MockAuth)
		api.POST("/user/createLoginLink", userController.CreateLoginLink)
		api.POST("/user/authLoginLink", userController.AuthLoginLink)
		api.POST("/user/logout", userController.Logout)
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

		api.GET("/event", eventController.GetAll)
		api.POST("/event", eventController.Create)
		api.GET("/event/:id", eventController.Get)
		api.PUT("/event/:id", eventController.Update)
		api.DELETE("/event/:id", eventController.Delete)
		api.GET("/event/:id/checkParticipant", eventController.CheckParticipant)
		api.POST("/event/:id/start", eventController.Start)
		api.POST("/event/:id/pause", eventController.Pause)
		api.POST("/event/:id/stop", eventController.Stop)
		api.POST("/event/:id/join", eventController.Join)
		api.POST("/event/:id/updateScore", eventController.UpdateScore)
		api.POST("/event/:id/rejoin", eventController.Rejoin)
		api.POST("/event/:id/exit", eventController.Exit)

		api.GET("/subworld/template/:id", subworldTemplateController.GetById)

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

		api.GET("/setting", systemSettingController.GetByInfo)
		api.POST("/setting", systemSettingController.CreateOrUpdate)

	}

	utilsFe := app.Group("/pages")
	{
		utilsFe.GET("/user/steam_login", userController.HandleSteamLogin)
	}
}
