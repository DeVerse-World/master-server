package response

import "github.com/hyperjiang/gin-skeleton/manager/schema"

type FetchAssetResp struct {
	Nfts []schema.NftStruct `json:"nfts"`
}
