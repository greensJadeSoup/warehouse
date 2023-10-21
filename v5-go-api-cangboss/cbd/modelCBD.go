package cbd

import "mime/multipart"

// ------------------------ req ------------------------
type AddModelReqCBD struct {
	SellerID uint64 `json:"seller_id"  form:"seller_id" binding:"required,gte=1"`
	ItemID   uint64 `json:"item_id"  form:"item_id" binding:"required,gte=1"`
	SKUList  string `json:"sku_list"  form:"sku_list" binding:"required,lte=255"`

	ShopID         uint64
	PlatformShopID string
	PlatformItemID string
	Detail         []ModelImageDetailCBD
}

type EditModelReqCBD struct {
	SellerID    uint64 `json:"seller_id"  form:"seller_id"  binding:"required,gte=1"`
	ID          uint64 `json:"id"  form:"id"  binding:"required,gte=1"`
	ModelSku    string `json:"model_sku"  form:"model_sku"  binding:"required,lte=255"`
	ImageChange bool   `json:"image_change"  form:"image_change"  binding:"omitempty"`

	Url     string
	Image   *multipart.FileHeader
	TmpPath string
}

type DelModelReqCBD struct {
	SellerID uint64 `json:"seller_id" binding:"required,gte=1"`
	ID       uint64 `json:"id,string"  binding:"required,gte=1"`
}

type ListModelReqCBD struct {
	SellerID       uint64 `json:"seller_id" form:"seller_id" xorm:"seller_id" binding:"required,gte=1"`
	ShopID         uint64 `json:"shop_id" form:"shop_id" xorm:"shop_id" binding:"required,gte=1"`
	PlatformItemID uint64 `json:"platform_item_id" form:"platform_item_id" xorm:"platform_item_id" binding:"required,gte=1"`
	ModelID        uint64 `json:"model_id" form:"model_id" xorm:"model_id" binding:"required,gte=1"`
	ModelSku       string `json:"model_sku" form:"model_sku" xorm:"model_sku" binding:"required,lte=255"`
	Remark         string `json:"remark" form:"remark" xorm:"remark"`

	IsPaging  bool `json:"is_paging" form:"is_paging"`
	PageIndex int  `json:"page_index" form:"page_index" binding:"required"`
	PageSize  int  `json:"page_size" form:"page_size" binding:"required"`
}

type ModelImageDetailCBD struct {
	Sku             string `json:"sku" form:"sku"`
	Url             string `json:"url" form:"url"`
	PlatformModelID string `json:"platform_model_id" form:"platform_model_id"`
	ModelID         uint64 `json:"model_id,string"`
	Image           *multipart.FileHeader
	TmpPath         string
}

type ModelDetailCBD struct {
	ID       uint64 `json:"id"  xorm:"id"`
	Platform string `json:"platform"  xorm:"platform"`
	SellerID uint64 `json:"seller_id"  xorm:"seller_id"`

	ShopID         uint64 `json:"shop_id"  xorm:"shop_id"`
	ShopName       string `json:"shop_name"  xorm:"shop_name"`
	PlatformShopID string `json:"platform_shop_id"  xorm:"platform_shop_id"`
	Region         string `json:"region"  xorm:"region"`

	ItemID         uint64 `json:"item_id"  xorm:"item_id"`
	PlatformItemID string `json:"platform_item_id"  xorm:"platform_item_id"`

	PlatformModelID string `json:"platform_model_id"  xorm:"platform_model_id"`

	ItemName   string `json:"item_name"  xorm:"item_name"`
	ItemStatus string `json:"item_status"  xorm:"item_status"`

	ItemSku  string `json:"item_sku"  xorm:"item_sku"`
	ModelSku string `json:"model_sku"  xorm:"model_sku"`

	ModelIsDelete uint8 `json:"model_is_delete"  xorm:"model_is_delete"`

	ItemImages  string `json:"item_images"  xorm:"item_images"`
	ModelImages string `json:"model_images"  xorm:"model_images"`
	Remark      string `json:"remark"  xorm:"remark"`
}

// --------------------resp-------------------------------
type ListModelRespCBD struct {
	ID              uint64 `json:"id"  xorm:"id pk autoincr"  binding:"required,gte=1"`
	SellerID        uint64 `json:"seller_id"  xorm:"seller_id"  binding:"required,gte=1"`
	ShopID          uint64 `json:"shop_id"  xorm:"shop_id"  binding:"required,gte=1"`
	PlatformShopID  uint64 `json:"platform_shop_id"  xorm:"platform_shop_id"`
	PlatformItemID  uint64 `json:"platform_item_id"  xorm:"platform_item_id"  binding:"required,gte=1"`
	PlatformModelID uint64 `json:"platform_model_id"  xorm:"platform_model_id"  binding:"required,gte=1"`
	ModelSku        string `json:"model_sku"  xorm:"model_sku"  binding:"required,lte=255"`
	Remark          string `json:"remark"  xorm:"remark"`
}
