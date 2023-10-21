package cbd

//------------------------ req ------------------------
type AddSendWayReqCBD struct {
	VendorID		uint64			`json:"vendor_id" binding:"required"`
	LineID			uint64			`json:"line_id"  xorm:"line_id"  binding:"required"`
	Type			string			`json:"type"  xorm:"type"  binding:"required,lt=32"`
	Name			string			`json:"name"  xorm:"name"  binding:"required,lt=32"`
	Sort			int			`json:"sort"  xorm:"sort"  binding:"omitempty,gte=0"`
	Note			string			`json:"note"  xorm:"note"  binding:"lt=255"`
}

type ListSendWayReqCBD struct {
	VendorID	uint64		`json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	SellerID	uint64		`json:"seller_id"  form:"seller_id" binding:"omitempty,gte=1"`
	LineID		uint64		`json:"line_id" form:"line_id" binding:"required,gte=1"`
	LineIDList	[]string

	IsPaging	bool		`json:"is_paging" form:"is_paging"`
	PageIndex	int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize	int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditSendWayReqCBD struct {
	VendorID		uint64			`json:"vendor_id" binding:"required"`
	ID			uint64			`json:"id"  xorm:"id pk autoincr"  binding:"required,gte=1"`
	Name			string			`json:"name"  xorm:"name"  binding:"required,lt=32"`
	Sort			int			`json:"sort"  xorm:"sort"  binding:"omitempty,gte=0"`
	Note			string			`json:"note"  xorm:"note"  binding:"lt=255"`
}

type DelSendWayReqCBD struct {
	VendorID	uint64		`json:"vendor_id" binding:"required"`
	ID		uint64		`json:"id"  xorm:"id pk autoincr"  binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type ListSendWayRespCBD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	LineID			uint64		`json:"line_id"  xorm:"line_id"`
	Type			string		`json:"type"  xorm:"type"`
	Name			string		`json:"name"  xorm:"name"`
	Sort			int		`json:"sort"  xorm:"sort"`
	Note			string		`json:"note"  xorm:"note"`
	RoundUp			uint8		`json:"round_up" xorm:"round_up"`
	AddKg			float64		`json:"add_kg" xorm:"add_kg"`
	PriFirstWeight		float64		`json:"pri_first_weight"  xorm:"pri_first_weight"`
	WeightPriceRules	string		`json:"weight_pri_rules"  xorm:"weight_pri_rules"`
}
