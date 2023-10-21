package cp_api

import "warehouse/v5-go-component/cp_obj"

type SellerShopDetail struct {
	Platform	string	`json:"platform"`
	ShopCount	int	`json:"shop_count"`
	CbCount		int	`json:"cb_count"`
}

type WarehouseDetail struct {
	WarehouseID		uint64			`json:"warehouse_id"`
	Name			string			`json:"name"`
	Role			string			`json:"role"`
	Store			bool			`json:"store"`
}

type LineDetail struct {
	LineID			uint64		`json:"line_id"`
	Source			uint64		`json:"source"`
	To			uint64		`json:"to"`
	SourceWhr		string		`json:"source_whr"`
	ToWhr			string		`json:"to_whr"`
}

type VendorDetail struct {
	VendorID		uint64				`json:"vendor_id"`
	VendorName		string				`json:"vendor_name"`
	VendorAlias		string				`json:"vendor_alias"`
	LineDetail		[]LineDetail			`json:"line_detail"`
	WarehouseDetail		[]WarehouseDetail		`json:"warehouse_detail"` //如果是仓库, 仓管管理了哪些仓库
}

type CheckSessionInfo struct {
	UserID			uint64		`json:"user_id"`
	ParentID		uint64		`json:"parent_id"`
	Account 		string		`json:"account"`
	AccountType		string		`json:"account_type"`
	WareHouseRole		string		`json:"warehouse_role"`
	RealName 		string		`json:"real_name"`
	Email 			string		`json:"email"`
	Phone 			string		`json:"phone"`
	CompanyName 		string		`json:"company_name"`
	WechatNum 		string		`json:"wechat_num"`
	DeviceType		string		`json:"-"`
	DeviceInfo		string		`json:"-"`
	LoginType      		string    	`json:"-"`// account:账号登录 email:邮箱 phone:手机 third:第三方 face:扫脸登陆

	ManagerID		uint64		`json:"-"`
	IsManager		bool		`json:"-"`
	IsSuperManager		bool		`json:"-"`

	SessionKey		string		`json:"session_key"`
	AllowLogin		uint8		`json:"allow_login"`

	VendorDetail		[]VendorDetail		`json:"vendor_detail"`
	SellerShopDetail	[]SellerShopDetail	`json:"seller_shop_detail"`

	Kick			uint8		`json:"-"`
	DeviceID		string		`json:"-"`
	AppVersion		string		`json:"-"`
	LoginTime		cp_obj.Datetime	`json:"login_time"`
	LogoutTime		cp_obj.Datetime	`json:"-"`
	LastActiveDate		cp_obj.Datetime	`json:"last_active_time"`
	LastActiveIP		string		`json:"last_active_ip"`
	ExpireTime		cp_obj.Datetime	`json:"expire_time"`
}

type SingleOrder struct {
	OrderID			uint64		`json:"order_id,string" binding:"required,gte=1"`
	OrderTime		int64		`json:"order_time" binding:"required,gte=1"`
}

type BatchOrderReq struct {
	SellerID		uint64		`json:"seller_id"`
	VendorID		uint64		`json:"vendor_id"`
	OrderList		[]SingleOrder	`json:"order_list"`
}

type BatchOrderResp struct {
	OrderID			uint64		`json:"order_id,string"`
	SN			string		`json:"sn"`
	Success			bool		`json:"success"`
	Reason			string		`json:"reason"`
}

type GetTrackInfoItemResp struct {
	UpdateTime			int64		`json:"update_time"`
	Description			string		`json:"description"`
	LogisticsStatus			string		`json:"logistics_status"`
}

type GetTrackInfoReqCBD struct {
	SellerID		uint64		`json:"seller_id" form:"seller_id"`
	VendorID		uint64		`json:"vendor_id" form:"vendor_id"`
	OrderID			uint64		`json:"order_id" form:"order_id"`
	OrderTime		int64		`json:"order_time" form:"order_time"`
}
