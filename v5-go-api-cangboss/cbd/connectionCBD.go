package cbd

import (
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_obj"
)

// ------------------------ req ------------------------
type AddConnectionReqCBD struct {
	VendorID   uint64 `json:"vendor_id"  binding:"required,gte=1"`
	CustomsNum string `json:"customs_num"  binding:"required,lte=32"`
	Platform   string `json:"platform"  binding:"omitempty,eq=shopee|eq=manual|eq=stock_up"`
	Note       string `json:"note"  binding:"omitempty,lte=1024"`
}

type ListConnectionReqCBD struct {
	VendorID   uint64 `json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	CustomsNum string `json:"customs_num" form:"customs_num" binding:"omitempty,lte=512"`
	MidNum 	   string `json:"mid_num" form:"mid_num" binding:"omitempty,lte=32"`
	SN 	   string `json:"sn" form:"sn" binding:"omitempty,lte=32"`
	MidType    string `json:"mid_type" form:"mid_type" binding:"omitempty"`
	Status     string `json:"status" form:"status" binding:"omitempty,eq=init|eq=stock_out|eq=customs|eq=arrive"`
	NoteKey	   string `json:"note_key" form:"note_key" binding:"omitempty,lte=64"`
	From	   int64  `json:"from" form:"from" binding:"omitempty,gte=1"`
	To	   int64  `json:"to" form:"to" binding:"omitempty,gte=1"`
	LineID	   uint64 `json:"line_id" form:"line_id" binding:"omitempty,gte=1"`
	SendWayID  uint64 `json:"sendway_id" form:"sendway_id" binding:"omitempty,gte=1"`

	MidTypeList	[]string
	CustomsNumList  []string
	LineIDList 	[]string
	ExcelOutput	bool

	IsPaging  bool `json:"is_paging" form:"is_paging"`
	PageIndex int  `json:"page_index" form:"page_index" binding:"required"`
	PageSize  int  `json:"page_size" form:"page_size" binding:"required"`
}

type EditConnectionReqCBD struct {
	ID         uint64 `json:"id"  binding:"required,gte=1"`
	VendorID   uint64 `json:"vendor_id"  binding:"required,gte=1"`
	CustomsNum string `json:"customs_num"  binding:"required,lte=32"`
	Platform   string `json:"platform"  binding:"omitempty,eq=shopee|eq=manual|eq=stock_up"`
	Note       string `json:"note"  binding:"omitempty,lte=1024"`

	MdConn		*model.ConnectionMD
}

type ChangeConnectionReqCBD struct {
	IDList     []uint64 `json:"id_list"  binding:"omitempty"`
	CustomsNum string `json:"customs_num"  binding:"omitempty,lte=64"`
	VendorID   uint64 `json:"vendor_id"  binding:"required,gte=1"`
	Status     string `json:"status"  binding:"required,eq=stock_out|eq=customs|eq=arrive"`

	ID     	   uint64
}

type GetConnectionReqCBD struct {
	ID       uint64 `json:"id" form:"id" binding:"omitempty,gte=1"`
	CustomsNum string `json:"customs_num" form:"customs_num"  binding:"omitempty,lte=64"`
	VendorID uint64 `json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
}

type DeductConnectionReqCBD struct {
	ID       	uint64 `json:"id"  binding:"omitempty,gte=1"`
	CustomsNum  	string `json:"customs_num"  binding:"omitempty,lte=64"`
	VendorID 	uint64 `json:"vendor_id"  binding:"required,gte=1"`
}

type DelConnectionReqCBD struct {
	VendorID uint64 `json:"vendor_id"  binding:"required,gte=1"`
	ID       uint64 `json:"id"  binding:"required,gte=1"`
}

// --------------------resp-------------------------------
type DeductFailListCBD struct {
	FailOrderCount int     `json:"fail_order_count"  binding:"required,gte=1"`
	FailFee        float64 `json:"fail_fee"  binding:"required,gte=1"`
}

type ListConnectionRespCBD struct {
	ID             uint64                        `json:"id"  xorm:"id pk autoincr"`
	VendorID       uint64                        `json:"vendor_id"  xorm:"vendor_id"`
	CustomsNum     string                        `json:"customs_num"  xorm:"customs_num"`
	Platform       string                        `json:"platform"  xorm:"platform"`
	Status         string                        `json:"status"  xorm:"status"`
	MidType         string                       `json:"mid_type"  xorm:"mid_type"`
	Weight         float64                       `json:"weight"  xorm:"weight"`
	OrderCount     int                           `json:"order_count"  xorm:"order_count"`
	WarehouseID    uint64                        `json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName  string                        `json:"warehouse_name"  xorm:"warehouse_name"`
	LineID         uint64                        `json:"line_id"  xorm:"line_id"`
	Source         uint64                        `json:"source"  xorm:"source"`
	To             uint64                        `json:"to"  xorm:"to"`
	SourceName     string                        `json:"source_name"  xorm:"source_name"`
	ToName         string                        `json:"to_name"  xorm:"to_name"`
	SendWayID      uint64                        `json:"sendway_id"  xorm:"sendway_id"`
	SendWayType    string                        `json:"sendway_type"  xorm:"sendway_type"`
	SendWayName    string                        `json:"sendway_name" xorm:"sendway_name"`
	Note           string                        `json:"note"  xorm:"note"`
	CreateTime     cp_obj.Datetime               `json:"create_time"  xorm:"create_time"`
	DeductFailList map[string]*DeductFailListCBD `json:"deduct_fail_list"  xorm:"create_time"`
}

type GetConnectionRespCBD struct {
	ListConnectionRespCBD
	MidCount		int		`json:"mid_count"`
	MidWeight		float64		`json:"mid_weight"`
}
