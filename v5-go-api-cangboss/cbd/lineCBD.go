package cbd

//------------------------ line ------------------------
type AddLineReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"required"`

	Source			uint64		`json:"source" binding:"required"`
	To			uint64		`json:"to" binding:"required"`
	Note			string		`json:"note" binding:"lt=255"`
	Sort			int		`json:"sort" binding:"omitempty,gte=0"`
}

type ListLineReqCBD struct {
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	SellerID		uint64		`json:"seller_id"  form:"seller_id"  binding:"omitempty,gte=1"`
	Source			uint64		`json:"source"  form:"source"  binding:"omitempty,gte=1"`
	To			uint64		`json:"to"  form:"to"  binding:"omitempty,gte=1"`
	WarehouseID		uint64		`json:"warehouse_id"  form:"warehouse_id"  binding:"omitempty,gte=1"`

	VendorIDList		[]string
	WarehouseIDList		[]string	//sso和check有用到
	LineIDList		[]string

	IsPaging		bool     	`json:"is_paging" form:"is_paging"`
	PageIndex		int      	`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int      	`json:"page_size" form:"page_size" binding:"required"`
}

type EditLineReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"required"`

	LineID			uint64		`json:"line_id" binding:"required"`
	Source			uint64		`json:"source" binding:"required"`
	To			uint64		`json:"to" binding:"required"`
	Note			string		`json:"note" binding:"lt=255"`
	Sort			int		`json:"sort" binding:"omitempty,gte=0"`
}

type DelLineReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"required"`
	LineID			uint64		`json:"line_id" binding:"required"`
}

type GetLineCBD struct {
	ID 				uint64		`json:"id" xorm:"id pk autoincr"`
	VendorID 			uint64		`json:"vendor_id" xorm:"vendor_id"`
	Source	 			uint64		`json:"source" xorm:"source"`
	To	 			uint64		`json:"to" xorm:"to"`
	SourceName                   	string        	`json:"source_name" xorm:"source_name"`
	ToName                      	string        	`json:"to_name" xorm:"to_name"`
	Sort				int		`json:"sort" xorm:"sort"`
	Note				string		`json:"note" xorm:"note"`
}

//--------------------resp-------------------------------
type ListLineRespCBD struct {
	ID 			uint64			`json:"id" xorm:"id"`
	VendorID 		uint64			`json:"vendor_id" xorm:"vendor_id"`
	Source			uint64			`json:"source" xorm:"source"`
	To			uint64			`json:"to" xorm:"to"`
	SourceWhr		string			`json:"source_whr" xorm:"source_whr"`
	ToWhr			string			`json:"to_whr" xorm:"to_whr"`
	Note			string			`json:"note" xorm:"note"`
	Sort			int			`json:"sort" xorm:"sort"`
	Detail			[]ListSendWayRespCBD `json:"detail"`
}
