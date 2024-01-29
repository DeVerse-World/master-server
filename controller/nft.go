package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/hyperjiang/gin-skeleton/manager"
	requestSchema "github.com/hyperjiang/gin-skeleton/manager/schema/request"
	"github.com/hyperjiang/gin-skeleton/model"
)

type NftController struct {
	inMemoryStorageManager *manager.InMemoryStorageManager
}

func NewNftController(
	inMemoryStorageManager *manager.InMemoryStorageManager,
) *NftController {
	return &NftController{
		inMemoryStorageManager,
	}
}

func (ctrl *NftController) CreateMintNftLink(c *gin.Context) {
	const (
		success = "Create mint nft link successfully"
		failed  = "Create mint nft link unsuccessfully"
	)

	var req requestSchema.CreateMinkLink
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{
		"mint_nft_url": os.Getenv("UI_HOST") + "/marketplace/create-item?fileUri=" + req.IpfsHash,
	})
}

func (ctrl *NftController) NotifyMinted(c *gin.Context) {
	const (
		success = "Notify Minted successfully"
		failed  = "Notify Minted unsuccessfully"
	)

	var req requestSchema.NotifyMinted
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	var mintedNft = req.MintedNft
	if err := mintedNft.Create(); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, mintedNft)
}

func (ctrl *NftController) CheckName(c *gin.Context) {
	const (
		success = "Check Name successfully"
		failed  = "Check Name unsuccessfully"
	)
	name := c.Request.URL.Query().Get("name")

	var data interface{}
	if err := ctrl.inMemoryStorageManager.Get("nft_"+name, &data); err != nil && err.Error() != "cache miss" {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if data != nil {
		JSONReturn(c, http.StatusOK, success, gin.H{"exist": true})
		return
	}

	var mintedNft model.MintedNft
	if err := mintedNft.GetByName(name); err != nil && err != model.ErrDataNotFound {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}
	if mintedNft.ID > 0 {
		JSONReturn(c, http.StatusOK, success, gin.H{"exist": true})
		return
	}

	JSONReturn(c, http.StatusOK, success, gin.H{"exist": false})
}

func (ctrl *NftController) LockName(c *gin.Context) {
	const (
		success = "Lock Name successfully"
		failed  = "Lock Name unsuccessfully"
	)

	var req requestSchema.LockName
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if err := ctrl.inMemoryStorageManager.Set("nft_"+req.Name, 1); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, nil)
}

func (ctrl *NftController) UnlockName(c *gin.Context) {
	const (
		success = "Unlock Name successfully"
		failed  = "Unlock Name unsuccessfully"
	)

	var req requestSchema.LockName
	if err := c.BindJSON(&req); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if err := ctrl.inMemoryStorageManager.Delete("nft_" + req.Name); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, nil)
}

// MintNft mints a new NFT.
func (ctrl *NftController) MintNft(c *gin.Context) {
	const (
		success = "Mint NFT successfully"
		failed  = "Mint NFT unsuccessfully"
	)
	var (
		request           requestSchema.MintNft
		thxnetRequestBody []byte
		response          interface{}
		err               error
	)

	if err = c.BindJSON(&request); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if thxnetRequestBody, err = json.Marshal(request); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if response, err = requestThxNet("POST", "/l1/nft/mint", thxnetRequestBody); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, response)
}

// TransferNft transfers an NFT to another user.
func (ctrl *NftController) TransferNft(c *gin.Context) {
	const (
		success = "Transfer NFT successfully"
		failed  = "Transfer NFT unsuccessfully"
	)

	var (
		request           requestSchema.TransferNft
		thxnetRequestBody []byte
		response          interface{}
		err               error
	)

	if err = c.BindJSON(&request); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if thxnetRequestBody, err = json.Marshal(request); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if response, err = requestThxNet("POST", "/l1/nft/transfer", thxnetRequestBody); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, response)
}

// BurnNft burns an NFT.
func (ctrl *NftController) BurnNft(c *gin.Context) {
	const (
		success = "Burn NFT successfully"
		failed  = "Burn NFT unsuccessfully"
	)

	var (
		request           requestSchema.BurnNft
		thxnetRequestBody []byte
		response          interface{}
		err               error
	)

	if err = c.BindJSON(&request); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if thxnetRequestBody, err = json.Marshal(request); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if response, err = requestThxNet("POST", "/l1/nft/burn", thxnetRequestBody); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, response)
}

// UpdateDynamicNft updates the dynamic NFT.
func (ctrl *NftController) UpdateDynamicNft(c *gin.Context) {
	const (
		success = "Update Dynamic NFT successfully"
		failed  = "Update Dynamic NFT unsuccessfully"
	)

	var (
		request           requestSchema.UpdateDynamicNft
		thxnetRequestBody []byte
		response          interface{}
		err               error
	)

	if err = c.BindJSON(&request); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if thxnetRequestBody, err = json.Marshal(request); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	if response, err = requestThxNet("PUT", "/l1/nft/metadata/dynamic", thxnetRequestBody); err != nil {
		abortWithStatusError(c, http.StatusBadRequest, failed, err)
		return
	}

	JSONReturn(c, http.StatusOK, success, response)
}

func requestThxNet(method string, path string, requestBody []byte) (interface{}, error) {
	const BaseUrl = "https://api.helpers.testnet.thxnet.org/rest/v0.5"
	var (
		request      *http.Request
		response     *http.Response
		responseBody []byte
		result       map[string]interface{}
		err          error
	)

	// Create a new HTTP request
	if request, err = http.NewRequest(method, BaseUrl+path, bytes.NewReader(requestBody)); err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("BEARER_TOKEN")))

	request.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	if response, err = client.Do(request); err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	// Read response
	defer response.Body.Close()

	if responseBody, err = ioutil.ReadAll(response.Body); err != nil {
		return nil, fmt.Errorf("error getting response: %w", err)
	}

	// Form result
	if err = json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return result, nil
}
