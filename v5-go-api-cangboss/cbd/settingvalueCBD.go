package cbd

//------------------------ req ------------------------
type AddSettingValueReqCBD struct {
	VendorID		uint64		`json:"vendor_id"  binding:"required,gte=1"`
	Type		string		`json:"type"  binding:""`
	Value		string		`json:"value"  binding:""`
}

type ListSettingValueReqCBD struct {
	ID		uint64		`json:"id" form:"id" binding:"required,gte=1"`
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	IsPaging		bool		`json:"is_paging" form:"is_paging"`
	PageIndex		int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditSettingValueReqCBD struct {
	ID		uint64		`json:"id"  binding:"required,gte=1"`
	VendorID		uint64		`json:"vendor_id"  binding:"required,gte=1"`
	Type		string		`json:"type"  binding:""`
	Value		string		`json:"value"  binding:""`
}

type DelSettingValueReqCBD struct {
	ID		uint64		`json:"id"  binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type ListSettingValueRespCBD struct {
	ID		uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	Type		string		`json:"type"  xorm:"type"`
	Value		string		`json:"value"  xorm:"value"`
}
