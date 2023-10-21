package cbd

import "time"

//------------------------ req ------------------------
type AuthShopReqCBD struct {
	Platform		string		`json:"platform"  binding:"required,eq=shopee|eq=manual"`
	SellerID		uint64		`json:"seller_id"   binding:"required,gte=1"`
	SpecialID		string
	Host			string
}

type BindingShopReqCBD struct {
	PlatformShopID		string		`json:"shop_id"  form:"shop_id"`
	MainAccountID		uint64		`json:"main_account_id"  form:"main_account_id"`
	Code			string		`json:"code"  form:"code"  binding:"required"`
	SpecialID		string
}

type ShopDetail struct {
	Platform		string		`json:"platform"  binding:"required,eq=shopee|eq=manual"`
	ID			uint64		`json:"id" binding:"required"`
}

type SyncShopReqCBD struct {
	SellerID		uint64		`json:"seller_id"  binding:"omitempty,gte=1"`
	VendorID		uint64		`json:"vendor_id"  binding:"omitempty,gte=1"`
	ShopDetail		[]ShopDetail	`json:"shop_detail" binding:"required"`
}

type AddShopReqCBD struct {
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	Platform		string		`json:"platform"  xorm:"platform"`
	PlatformShopID		string		`json:"platform_shop_id"  xorm:"platform_shop_id"`
	AccessToken		string		`json:"access_token"  xorm:"access_token"  binding:""`
	RefreshToken		string		`json:"refresh_token"  xorm:"refresh_token"  binding:""`
	AccessExpire		time.Time	`json:"access_expire" xorm:"access_expire"`
	RefreshExpire		time.Time	`json:"refresh_expire" xorm:"refresh_expire"`
	ShopExpire		time.Time	`json:"shop_expire" xorm:"shop_expire"`
	Name			string		`json:"name"  xorm:"name"`
	Status			string		`json:"status"  xorm:"status"`
	Region			string		`json:"region"  xorm:"region"`
	IsCB			int8		`json:"is_cb"  xorm:"is_cb"`
	IsCNSC			int8		`json:"is_cnsc"  xorm:"is_cnsc"`
	IsSIP			int8		`json:"is_sip" xorm:"is_sip"`
	Logo			string		`json:"logo"  xorm:"logo"`
	Description		string		`json:"description"  xorm:"description"`
}

type EditShopReqCBD struct {
	ID			uint64		`json:"id"  xorm:"id"  binding:"required,gte=1"`
	AccessToken		string		`json:"access_token"  xorm:"access_token"  binding:""`
	RefreshToken		string		`json:"refresh_token"  xorm:"refresh_token"  binding:""`
	AccessExpire		time.Time	`json:"access_expire" xorm:"access_expire"`
	RefreshExpire		time.Time	`json:"refresh_expire" xorm:"refresh_expire"`
	ShopExpire		time.Time	`json:"shop_expire" xorm:"shop_expire"`
	Name			string		`json:"name"  xorm:"name"`
	Status			string		`json:"status"  xorm:"status"`
	Region			string		`json:"region"  xorm:"region"`
	IsCB			int8		`json:"is_cb"  xorm:"is_cb"`
	IsCNSC			int8		`json:"is_cnsc"  xorm:"is_cnsc"`
	IsSIP			int8		`json:"is_sip" xorm:"is_sip"`
	Logo			string		`json:"logo"  xorm:"logo"`
	Description		string		`json:"description"  xorm:"description"`
}

type RefreshShopReqCBD struct {
	Platform		string		`json:"platform"  binding:"required,eq=shopee|eq=manual"`
	ID			uint64		`json:"id"  xorm:"id"  binding:"required,gte=1"`
	AccessToken		string		`json:"access_token"  xorm:"access_token"  binding:""`
	RefreshToken		string		`json:"refresh_token"  xorm:"refresh_token"  binding:""`
	AccessExpire		time.Time	`json:"access_expire" xorm:"access_expire"`
	RefreshExpire		time.Time	`json:"refresh_expire" xorm:"refresh_expire"`
}

type DelShopReqCBD struct {
	Platform		string		`json:"platform"  binding:"required,eq=shopee|eq=manual"`
	ID			uint64		`json:"id"  xorm:"id"  binding:"required,gte=1"`
}