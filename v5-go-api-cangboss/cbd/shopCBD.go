package cbd

import (
	"warehouse/v5-go-component/cp_obj"
)

//------------------------ req ------------------------
type ListShopReqCBD struct {
	VendorID		uint64		`json:"vendor_id"  form:"vendor_id" binding:"omitempty,gte=1"`
	SellerID		uint64		`json:"seller_id" form:"seller_id" binding:"required,gte=1"`
	SellerIDList		[]string

	ShopKey			string		`json:"shop_key"  form:"shop_key"  binding:"omitempty,lte=64"`
	SellerKey		string		`json:"seller_key"  form:"seller_key"  binding:"omitempty,lte=64"`
	Platform		string		`json:"platform"  form:"platform"  binding:"omitempty,eq=shopee|eq=manual"`

	IsPaging		bool		`json:"is_paging" form:"is_paging"`
	PageIndex		int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int		`json:"page_size" form:"page_size" binding:"required"`
}

type ChangeAccountReqCBD struct {
	ShopID			uint64		`json:"shop_id"  binding:"required,gte=1"`
	NewSellerID		uint64		`json:"seller_id" binding:"required,gte=1"`

	OldSellerID		uint64
}

type DelShopReqCBD struct {
	ID			uint64		`json:"id"  xorm:"id"  binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type ListShopRespCBD struct {
	ID			uint64		`json:"id"  xorm:"id"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	Platform		string		`json:"platform"  xorm:"platform"`
	PlatformShopID		string		`json:"platform_shop_id"  xorm:"platform_shop_id"`
	Name			string		`json:"name"  xorm:"name"`
	Region			string		`json:"region"  xorm:"region"`
	ShopExpire		cp_obj.Datetime	`json:"shop_expire" xorm:"shop_expire"`
	Status			string		`json:"status"  xorm:"status"`
	IsCB			int8		`json:"is_cb"  xorm:"is_cb"`
	IsCNSC			int8		`json:"is_cnsc"  xorm:"is_cnsc"`
	Logo			string		`json:"logo"  xorm:"logo"`
	Description		string		`json:"description"  xorm:"description"`

	RealName		string		`json:"real_name"  xorm:"real_name"`
}