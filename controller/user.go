package controller

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/hyperjiang/gin-skeleton/manager"
	"github.com/hyperjiang/gin-skeleton/manager/jwt"
	requestSchema "github.com/hyperjiang/gin-skeleton/manager/schema/request"
	"github.com/hyperjiang/gin-skeleton/manager/steam"
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
	const (
		success = "Get User Private Profile successfully"
		failed  = "Get User Private Profile unsuccessfully"
	)
	authU, err := jwt.HandleUserCookie(c.Writer, c.Request)
	if err != nil && authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("Can't parse token"))
		return
	}

	var user model.User
	err2 := user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))
	if err2 != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("Not valid user token"))
		return
	}

	avatars, err := model.GetUserAvatars(user.ID)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	created_events, err := model.GetUserCreatedEvents(user.ID)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	created_root_subworld_templates, err := model.GetRootFromCreator(int(user.ID))
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	created_deriv_subworld_templates, err := model.GetAllDerivFromCreator(int(user.ID))
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"user":                             user,
		"avatars":                          avatars,
		"created_events":                   created_events,
		"created_root_subworld_templates":  created_root_subworld_templates,
		"created_deriv_subworld_templates": created_deriv_subworld_templates,
	})
}

func (ctrl *UserController) UpdateUserProfile(c *gin.Context) {
	const (
		success = "Update User Private Profile successfully"
		failed  = "Update User Private Profile unsuccessfully"
	)
	authU, err := jwt.HandleUserCookie(c.Writer, c.Request)
	if err != nil && authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	var req requestSchema.UpdateUserProfileReq
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	var user model.User
	err2 := user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))
	if err2 != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("Not valid user token"))
		return
	}

	user.Name = req.Name
	if err := user.Update(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"user": user,
	})
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
	const (
		success = "Get Or Create User successfully"
		failed  = "Get Or Create User unsuccessfully"
	)

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
						abortWithStatusError(c, http.StatusBadRequest, failed, err)
						return
					} else {
						lr.UpdateUserId(user.ID)
					}
				}
				JSONReturn(c, http.StatusOK, success, gin.H{
					"user":         user,
					"require_auth": false,
				})
				return
			}
		}

		var user model.User

		if err := user.GetOrCreateByWallet(req.WalletAddress); err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}

		JSONReturn(c, http.StatusOK, success, gin.H{
			"user":         user,
			"require_auth": true,
		})
	} else if req.LoginMode == "GOOGLE" {
		if authU != nil && authU.SocialEmail == req.GoogleEmail { // nonce message was signed
			var user model.User
			err := user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))
			if err == nil {
				// store login Link approval if session key exists
				if req.SessionKey != "" {
					var lr model.LoginRequest
					if err := lr.GetByKey(req.SessionKey); err != nil {
						abortWithStatusError(c, http.StatusBadRequest, failed, err)
						return
					} else {
						lr.UpdateUserId(user.ID)
					}
				}
				JSONReturn(c, http.StatusOK, success, gin.H{
					"user":         user,
					"require_auth": false,
				})
				return
			}
		}

		var user model.User

		if err := user.GetOrCreateBySocialEmail(req.GoogleEmail); err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}

		JSONReturn(c, http.StatusOK, success, gin.H{
			"user":         user,
			"require_auth": true,
		})
	} else {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("unsupported login mode"))
		return
	}
}

func (ctrl *UserController) Update(c *gin.Context) {
	const (
		success = "Update User successfully"
		failed  = "Update User unsuccessfully"
	)
	var req requestSchema.UpdateUserReq
	c.BindJSON(&req)

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	if authU == nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("can't parse token"))
		return
	}

	var user model.User
	if err := user.GetUserById(strconv.FormatUint(uint64(authU.ID), 10)); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.GoogleEmail != "" {
		user.SocialEmail = req.GoogleEmail
	}
	if req.WalletAddress != "" {
		user.WalletAddress = req.WalletAddress
		user.WalletNonce = util.GenerateRandomString(10)
	}
	if err := user.Update(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	JSONReturn(c, http.StatusOK, success, gin.H{
		"user": user,
	})
}

func (ctrl *UserController) MockAuth(c *gin.Context) {
	const (
		success = "Mock Auth successfully"
		failed  = "Mock Auth unsuccessfully"
	)
	var req requestSchema.MockAuthUserReq
	c.BindJSON(&req)

	if req.LoginMode == "METAMASK" {
		var user model.User
		if err := user.GetUserByWalletAddress(req.WalletAddress); err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}

		jwt.WriteUserCookie(c.Writer, &user)
		JSONReturn(c, http.StatusOK, success, gin.H{
			"user": user,
		})
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
	const (
		success = "Auth Login Link successfully"
		failed  = "Auth Login Link unsuccessfully"
	)

	var req requestSchema.AuthLoginLink
	c.BindJSON(&req)

	authU, _ := jwt.HandleUserCookie(c.Writer, c.Request)
	fmt.Println("AuthLoginLink AuthU")
	fmt.Println(authU)
	if authU != nil {
		var origUser model.User
		err := origUser.GetUserById(strconv.FormatUint(uint64(authU.ID), 10))
		if err == nil {
			// store login Link approval if session key exists
			if req.SessionKey != "" {
				var lr model.LoginRequest
				if err := lr.GetByKey(req.SessionKey); err != nil {
					abortWithStatusError(c, http.StatusBadRequest, failed, err)
					return
				} else {
					lr.UpdateUserId(origUser.ID)
				}
			}
			JSONReturn(c, http.StatusOK, success, gin.H{
				"user": origUser,
			})
			return
		}
	}

	var user model.User

	if req.LoginMode == "METAMASK" {
		if err := user.GetUserByWalletAddress(req.WalletAddress); err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}

		msg := "I am signing my one-time nonce: " + user.WalletNonce
		verifyResult := manager.VerifyWalletSig(user.WalletAddress, req.WalletSignature, []byte(msg))

		if !verifyResult {
			abortWithStatusError(c, http.StatusBadRequest, failed, errors.New("signature not verified"))
			return
		}

		jwt.WriteUserCookie(c.Writer, &user)

		if req.SessionKey != "" {
			var lr model.LoginRequest
			if err := lr.GetByKey(req.SessionKey); err != nil {
				abortWithStatusError(c, http.StatusBadRequest, failed, err)
				return
			} else {
				lr.UpdateUserId(user.ID)
			}
		}

		JSONReturn(c, http.StatusOK, success, gin.H{
			"user": user,
		})
		user.UpdateWalletNonce()
	} else if req.LoginMode == "GOOGLE" {
		if err := user.GetUserByGoogleEmail(req.GoogleEmail); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		claims, err := manager.ValidateGoogleJWT(req.GoogleToken)
		if err != nil {
			abortWithStatusError(c, http.StatusUnauthorized, failed, err)
			return
		}
		if claims.Email != user.SocialEmail {
			abortWithStatusError(c, http.StatusUnauthorized, failed, errors.New("emails mismatch"))
			return
		}

		jwt.WriteUserCookie(c.Writer, &user)

		if req.SessionKey != "" {
			var lr model.LoginRequest
			if err := lr.GetByKey(req.SessionKey); err != nil {
				abortWithStatusError(c, http.StatusBadRequest, failed, err)
				return
			} else {
				lr.UpdateUserId(user.ID)
			}
		}

		JSONReturn(c, http.StatusOK, success, gin.H{
			"user": user,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported login mode"})
		return
	}
}

func (ctrl *UserController) HandleSteamLoginOpenID(c *gin.Context) {
	const (
		success = "Handle Steam Login through OpenID successfully"
		failed  = "Handle Steam Login through OpenID unsuccessfully"
	)
	w, r := c.Writer, c.Request
	fmt.Println("Start Handling Steam Login " + r.Host)
	sessionKey := c.Request.URL.Query().Get("session_key")
	fmt.Println("Session Key " + sessionKey)
	opId := steam.NewOpenId(r, os.Getenv("API_HOST"))
	switch opId.Mode() {
	case "":
		http.Redirect(w, r, opId.AuthUrl(), 302)
	case "cancel":
		w.Write([]byte("Authorization cancelled"))
	default:
		fmt.Println("Steam Login 3rd party success")
		steamId, err := opId.ValidateAndGetId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		fmt.Println("Received steamID " + steamId)
		// Do whatever you want with steam id
		//w.Write([]byte(steamId))

		var user model.User
		if err := user.GetOrCreateBySteamId(steamId); err != nil {
			abortWithStatusError(c, http.StatusBadRequest, failed, err)
			return
		}
		if sessionKey != "" {
			var lr model.LoginRequest
			if err := lr.GetByKey(sessionKey); err != nil {
				fmt.Println("login request session key not found")
			} else {
				lr.UpdateUserId(user.ID)
			}
		}
		fmt.Println("UI Domain " + os.Getenv("UI_DOMAIN") + os.Getenv("UI_HOST"))
		jwt.WriteUserCookie(w, &user)
		JSONReturn(c, http.StatusOK, success, gin.H{
			"user": user,
		})
	}
}

func (ctrl *UserController) HandleSteamLoginSessionTicket(c *gin.Context) {
	const (
		success = "Handle Steam Login through InApp Session successfully"
		failed  = "Handle Steam Login through InApp Session unsuccessfully"
	)
	w := c.Writer
	ticket := c.Request.URL.Query().Get("ticket")
	steamId, err := manager.ValidateAndGetSteamIdByTicket(ticket)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	fmt.Println("Received steamID from session ticket " + steamId)
	var user model.User
	if err := user.GetOrCreateBySteamId(steamId); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	jwt.WriteUserCookie(w, &user)
	JSONReturn(c, http.StatusOK, success, gin.H{
		"user": user,
	})
}

func (ctrl *UserController) Logout(c *gin.Context) {
	const (
		success = "Logout successfully"
		failed  = "Logout unsuccessfully"
	)
	jwt.DeleteUserCookie(c.Writer)
	JSONReturn(c, http.StatusOK, success, gin.H{})
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

	avatars, err := model.GetUserAvatars(user.ID)
	if err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	JSONReturn(c, http.StatusOK, success, gin.H{
		"avatars": avatars,
	})
}
