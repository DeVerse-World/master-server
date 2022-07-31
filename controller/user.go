package controller

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/hyperjiang/gin-skeleton/manager/jwt"
	requestSchema "github.com/hyperjiang/gin-skeleton/manager/schema/request"
	"github.com/hyperjiang/gin-skeleton/manager/util"
	"github.com/hyperjiang/gin-skeleton/model"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (ctrl *UserController) GetUserByWalletAddress(c *gin.Context) {
	var user model.User

	walletAddress := c.Param("wallet_address")

	if err := user.GetUserByWalletAddress(walletAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func (ctrl *UserController) GetUserPrivateProfile(c *gin.Context) {
	authU, err := jwt.HandleUserCookie(c.Writer, c.Request)
	if err == nil || authU != nil { // TODO: Check err of jwt in all different places
		var user model.User
		err2 := user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))
		if err2 == nil {
			c.JSON(http.StatusOK, user)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (ctrl *UserController) UpdateAssets(c *gin.Context) {
	var req requestSchema.UpdateAssetReq
	c.BindJSON(&req)

	var user model.User
	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if err := user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	user.DeleteUserNft()

	// TODO: Create In Batch
	for i := 0; i < len(req.Nfts); i++ {
		var nft model.Nft
		var collection model.NftCollection
		collection.TokenAddress = req.Nfts[i].TokenAddress
		collection.GetOrCreate()

		nft.ImageUrl = req.Nfts[i].ImageUrl
		nft.UserId = user.ID
		nft.CollectionId = collection.ID
		nft.Create()
	}
}

func (ctrl *UserController) FetchAssets(c *gin.Context) {
	walletAddress := c.Param("wallet_address")
	var user model.User
	if err := user.GetUserByWalletAddress(walletAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if assets, err := user.FetchAssetsByUser(user.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, assets)
	}
}

func (ctrl *UserController) GetOrCreate(c *gin.Context) {
	var req requestSchema.SignupUserReq
	c.BindJSON(&req)

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)

	if req.LoginMode == "METAMASK" {
		if authU != nil && authU.WalletAddress == req.WalletAddress { // nonce message was signed
			var user model.User
			err := user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))
			if err == nil {
				// store login Link approval if session key exists
				if req.SessionKey != "" {
					var lr model.LoginRequest
					if err := lr.GetByKey(req.SessionKey); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					} else {
						lr.UpdateUserId(user.ID)
					}
				}
				return
			}
		}

		var user model.User

		if err := user.GetOrCreateByWallet(req.WalletAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, user)
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported login mode"})
	}
}

func (ctrl *UserController) Auth(c *gin.Context) {
	var req requestSchema.AuthUserReq
	c.BindJSON(&req)

	var user model.User

	if req.LoginMode == "METAMASK" {
		if err := user.GetUserByWalletAddress(req.WalletAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			msg := "I am signing my one-time nonce: " + user.WalletNonce
			verifyResult := util.VerifySig(user.WalletAddress, req.WalletSignature, []byte(msg))

			if !verifyResult {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			} else {
				jwt.WriteUserCookie(c.Writer, &user)
				c.JSON(http.StatusOK, user)
				user.UpdateWalletNonce()
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported login mode"})
	}

}

func (ctrl *UserController) MockAuth(c *gin.Context) {
	var req requestSchema.MockAuthUserReq
	c.BindJSON(&req)

	if req.LoginMode == "METAMASK" {
		var user model.User
		user.WalletAddress = req.WalletAddress

		jwt.WriteUserCookie(c.Writer, &user)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported login mode"})
	}
}

func (ctrl *UserController) CreateLoginLink(c *gin.Context) {

	var lr model.LoginRequest
	if err := lr.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"login_url": os.Getenv("UI_HOST") + "/login?key=" + lr.SessionKey})
	}
}

func (ctrl *UserController) AuthLoginLink(c *gin.Context) {
	var req requestSchema.AuthLoginLink
	c.BindJSON(&req)

	var user model.User

	if req.LoginMode == "METAMASK" {
		if err := user.GetUserByWalletAddress(req.WalletAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			msg := "I am signing my one-time nonce: " + user.WalletNonce
			verifyResult := util.VerifySig(user.WalletAddress, req.WalletSignature, []byte(msg))

			if !verifyResult {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			} else {
				jwt.WriteUserCookie(c.Writer, &user)

				if req.SessionKey != "" {
					var lr model.LoginRequest
					if err := lr.GetByKey(req.SessionKey); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
						return
					} else {
						lr.UpdateUserId(user.ID)
					}
				}

				c.JSON(http.StatusOK, user)
				user.UpdateWalletNonce()
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported login mode"})
		return
	}
}

func (ctrl *UserController) PollLoginLink(c *gin.Context) {
	sessionKey := c.Param("session_key")

	var lr model.LoginRequest
	if err := lr.GetByKey(sessionKey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if lr.UserId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user id => this log in link is not authorized yet"})
		return
	}

	var user model.User

	if err := user.GetUserById(strconv.FormatUint(uint64(*lr.UserId), 10)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt.WriteUserCookie(c.Writer, &user)
	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) GetTemporaryEventRewards(c *gin.Context) {
	const (
		success = "GetTemporaryEventRewards successfully"
		failed  = "GetTemporaryEventRewards unsuccessfully"
	)

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't find token"))
		return
	}
	var user model.User
	if err := user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10)); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	var rewardNfts []model.MintedNft
	rewardNfts, err := model.GetTemporaryRewards(user.ID)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"reward_nfts": rewardNfts,
	})
}

func (ctrl *UserController) GetAvatars(c *gin.Context) {
	const (
		success = "Get User Avatars successfully"
		failed  = "Get User Avatars unsuccessfully"
	)

	id := c.Param("id")
	var user model.User
	if err := user.GetUserById(id); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	avatars, err := user.GetUserAvatars(user.ID)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	JSONReturn(c, http.StatusOK, success, gin.H{
		"avatars": avatars,
	})
}
