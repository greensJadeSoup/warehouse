package cbd

//------------------------ req ------------------------
type AddConsumableReqCBD struct {
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	Name		string		`json:"name"  binding:"required,lte=255"`
	Note		string		`json:"note"  binding:"omitempty,lte=255"`
}

type ListConsumableReqCBD struct {
	VendorID	uint64		`json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	IsPaging	bool		`json:"is_paging" form:"is_paging"`
	PageIndex	int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize	int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditConsumableReqCBD struct {
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	ID		uint64		`json:"id"  binding:"required,gte=1"`
	Name		string		`json:"name"  binding:"required,lte=255"`
	Note		string		`json:"note"  binding:"omitempty,lte=255"`
}

type DelConsumableReqCBD struct {
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	ID		uint64		`json:"id"  binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type ListConsumableRespCBD struct {
	ID		uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID	uint64		`json:"vendor_id"  xorm:"vendor_id"`
	Name		string		`json:"name"  xorm:"name"`
	Note		string		`json:"note"  xorm:"note"`
}
