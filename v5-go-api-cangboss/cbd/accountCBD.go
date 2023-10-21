package cbd

//--------------------req-------------------------------
type AddManagerReqCBD struct {
	VendorID	uint64		`json:"vendor_id" binding:"required,gte=1"`

	WarehouseID	string		`json:"warehouse_id" binding:"required,gte=1"`
	Account 	string		`json:"account" binding:"required,lt=32"`
	Type		string		`json:"type" binding:"required,eq=manager|eq=super_manager|eq=service"`
	Password     	string		`json:"password" form:"password" binding:"required,lte=32"`
	RealName	string		`json:"real_name" binding:"omitempty,lt=16"`
	Phone		string		`json:"phone" binding:"omitempty,lt=32"`
	Email		string		`json:"email" binding:"omitempty,lt=32"`
	AllowLogin	uint8		`json:"allow_login" binding:"omitempty,eq=0|eq=1"`
	Note		string		`json:"note" binding:"omitempty,lt=255"`
}

type AddSellerReqCBD struct {
	VendorID	uint64		`json:"vendor_id" binding:"required,gte=1"`

	Account 	string		`json:"account" binding:"required,lt=32"`
	Password     	string		`json:"password" binding:"required,lte=32"`
	RealName	string		`json:"real_name" binding:"omitempty,lt=16"`
	Phone		string		`json:"phone" binding:"omitempty,lt=32"`
	Email		string		`json:"email" binding:"omitempty,lt=32"`
	AllowLogin	uint8		`json:"allow_login" binding:"omitempty,eq=0|eq=1"`
	Note		string		`json:"note" binding:"omitempty,lt=255"`
}

type EditManagerReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"required,gte=1"`

	ManagerID		uint64		`json:"manager_id" binding:"required,gte=1"`
	WarehouseID		string		`json:"warehouse_id" binding:"required,gte=1"`
	RealName		string		`json:"real_name" binding:"omitempty,lt=32"`
	Phone			string		`json:"phone" binding:"omitempty,lt=32"`
	Email			string		`json:"email" binding:"omitempty,lt=32"`
	AllowLogin		uint8		`json:"allow_login" binding:"omitempty,eq=0|eq=1"`
	Note			string		`json:"note" binding:"omitempty,lt=255"`
}

type EditSellerReqCBD struct {
	VendorID	uint64		`json:"vendor_id" binding:"required,gte=1"`

	SellerID	uint64		`json:"seller_id" binding:"required,gte=1"`
	RealName	string		`json:"real_name" binding:"omitempty,lt=32"`
	Phone		string		`json:"phone" binding:"omitempty,lt=32"`
	Email		string		`json:"email" binding:"omitempty,lt=32"`
	AllowLogin	uint8		`json:"allow_login" binding:"omitempty,eq=0|eq=1"`
	Note		string		`json:"note" binding:"omitempty,lt=255"`
}

type ListManagerReqCBD struct {
	VendorID	uint64		`json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`

	Type		string		`json:"type" form:"type" binding:"omitempty,eq=manager|eq=super_manager|eq=service"`
	WarehouseID	string		`json:"warehouse_id" form:"warehouse_id"`
	Account 	string		`json:"account" form:"account" binding:"omitempty,lt=32"`
	RealName	string		`json:"real_name" form:"real_name" binding:"omitempty,lt=16"`

	IsPaging	bool     	`json:"is_paging" form:"is_paging"`
	PageIndex	int      	`json:"page_index" form:"page_index" binding:"required"`
	PageSize	int      	`json:"page_size" form:"page_size" binding:"required"`
}

type ListSellerReqCBD struct {
	VendorID	uint64		`json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`

	Account 	string		`json:"account" form:"account" binding:"omitempty,lt=32"`
	RealName	string		`json:"real_name" form:"real_name" binding:"omitempty,lt=16"`

	IsPaging	bool     	`json:"is_paging" form:"is_paging"`
	PageIndex	int      	`json:"page_index" form:"page_index" binding:"required,gte=0"`
	PageSize	int      	`json:"page_size" form:"page_size" binding:"required,gte=0"`
}

type DelManagerReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"required,gte=1"`
	ManagerID		uint64		`json:"manager_id" binding:"required,gte=1"`
}

type DelSellerReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"required,gte=1"`
	SellerID		uint64		`json:"seller_id" binding:"required,gte=1"`
}

type ListBalanceReqCBD struct {
	SellerID		uint64		`json:"seller_id" form:"seller_id" binding:"required,gte=1"`
}

type ModifyPasswordReqCBD struct {
	VendorID	uint64		`json:"vendor_id" binding:"omitempty,gte=1"`
	SellerID	uint64		`json:"seller_id" binding:"omitempty,gte=1"`
	OldPassword	string		`json:"old_password" binding:"omitempty,lte=32"`
	NewPassword	string		`json:"new_password" binding:"required,lte=32"`

	Account 	string
	Type		string
	Salt		string
}

type CheckPasswordReqCBD struct {
	Account 	string
	HashPassword	string
	InPassword	string
	Salt		string
}

type EditBalanceReqCBD struct {
	VendorID	uint64		`json:"vendor_id" binding:"required,gte=1"`
	SellerID	uint64		`json:"seller_id" binding:"required,gte=1"`
	Type		string		`json:"type" binding:"required,eq=add|eq=sub"`
	Num		float64		`json:"num" binding:"required,gte=0"`
	Note		string		`json:"note" binding:"omitempty,lte=255"`
}

type EditProfileReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"omitempty,gte=1"`
	SellerID		uint64		`json:"seller_id" binding:"omitempty,gte=1"`
	RealName 		string		`json:"real_name" binding:"omitempty,lte=64"`
	Phone 			string		`json:"phone" binding:"omitempty,lte=32"`
	Email 			string		`json:"email" binding:"omitempty,lte=64"`
	CompanyName 		string		`json:"company_name" binding:"omitempty,lte=32"`
	WechatNum 		string		`json:"wechat_num" binding:"omitempty,lte=32"`
}

//--------------------resp-------------------------------
type ListSellerRespCBD struct {
	ID		uint64		`json:"id"`
	Account 	string		`json:"account"`
	RealName	string		`json:"real_name" xorm:"real_name"`
	Phone		string		`json:"phone"`
	Email		string		`json:"email"`
	AllowLogin	uint8		`json:"allow_login" xorm:"allow_login"`

	Balance		float64		`json:"balance" xorm:"balance"`
	DiscountID	uint64		`json:"discount_id" xorm:"discount_id"`
	DiscountEnable	uint8		`json:"-" xorm:"discount_enable"`
	DiscountName	string		`json:"discount_name" xorm:"discount_name"`
	Note		string		`json:"note" xorm:"note"`
}

type ListManagerRespCBD struct {
	ID		uint64		`json:"id"`
	Account 	string		`json:"account"`
	WarehouseID	string		`json:"warehouse_id" xorm:"warehouse_id"`
	WarehouseName	string		`json:"warehouse_name" xorm:"warehouse_name"`
	Type		string		`json:"type" xorm:"type"`
	WarehouseRole   string		`json:"warehouse_role" xorm:"warehouse_role"`
	RealName	string		`json:"real_name" xorm:"real_name"`
	Phone		string		`json:"phone"`
	Email		string		`json:"email"`
	AllowLogin	uint8		`json:"allow_login" xorm:"allow_login"`
	Note		string		`json:"note" xorm:"note"`
}
