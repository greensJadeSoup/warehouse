package cbd

//------------------------ req ------------------------
type AddVendorReqCBD struct {
	VendorName			string		`json:"vendor_name" binding:"required,lte=32"`
	SuperAdminAccount		string		`json:"super_admin_account" binding:"required,lte=32"`
}

type ListVendorReqCBD struct {
	IsPaging		bool		`json:"is_paging" form:"is_paging"`
	PageIndex		int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditVendorReqCBD struct {
	ID			uint64		`json:"id"  binding:"required,gte=1"`
	VendorID		uint64		`json:"vendor_id"  binding:"required,gte=1"`
	SellerID		uint64		`json:"seller_id"  binding:"required,gte=1"`
}

type DelVendorReqCBD struct {
	ID			uint64		`json:"id"  binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type ListVendorRespCBD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	Name			string		`json:"name" xorm:"name"`
	Balance			float64		`json:"balance" xorm:"balance"`
	OrderFee		float64		`json:"order_fee" xorm:"order_fee"`
	SuperManagerAccount	string		`json:"super_manager_account" xorm:"super_manager_account"`
	SuperManagerRealName	string		`json:"super_manager_real_name" xorm:"super_manager_real_name"`
}
