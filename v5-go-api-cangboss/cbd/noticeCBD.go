package cbd

import (
	"warehouse/v5-go-component/cp_obj"
)

//------------------------ req ------------------------
type AddNoticeReqCBD struct {
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	Title		string		`json:"title"  binding:"required,lte=255"`
	Content		string		`json:"content"  binding:"required,lte=4000"`
	IsTop		uint8		`json:"is_top"  binding:"omitempty,eq=0|eq=1"`
	Display		uint8		`json:"display"  binding:"omitempty,eq=0|eq=1"`
	Sort		int		`json:"sort"  binding:"omitempty,gte=0"`
}

type ListNoticeReqCBD struct {
	SellerID	uint64		`json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`
	VendorID	uint64		`json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	VendorIDList	[]string

	IsPaging	bool		`json:"is_paging" form:"is_paging"`
	PageIndex	int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize	int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditNoticeReqCBD struct {
	ID		uint64		`json:"id"  binding:"required,gte=1"`
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	Title		string		`json:"title"  binding:"required,lte=255"`
	Content		string		`json:"content"  binding:"required,lte=4000"`
	IsTop		uint8		`json:"is_top"  binding:"omitempty,eq=0|eq=1"`
	Display		uint8		`json:"display"  binding:"omitempty,eq=0|eq=1"`
	Sort		int		`json:"sort"  binding:"omitempty,gte=0"`
}

type DelNoticeReqCBD struct {
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	ID		uint64		`json:"id"  binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type ListNoticeRespCBD struct {
	ID		uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID	uint64		`json:"vendor_id"  xorm:"vendor_id"`
	Title		string		`json:"title"  xorm:"title"`
	Content		string		`json:"content"  xorm:"content"`
	IsTop		uint8		`json:"is_top"  xorm:"is_top"`
	Display		uint8		`json:"display"  xorm:"display"`
	Sort		int		`json:"sort"  xorm:"sort"`
	CreateTime	cp_obj.Datetime	`json:"create_time" xorm:"create_time"`
}
