package cbd

import (
	"warehouse/v5-go-api-cangboss/model"
	"warehouse/v5-go-component/cp_obj"
)

// ------------------------ req ------------------------
type AddMidConnectionReqCBD struct {
	VendorID   	uint64 `json:"vendor_id"  binding:"required,gte=1"`
	ConnectionID	uint64	`json:"connection_id"  binding:"omitempty,gte=1"`
	MidConnectionID	uint64	`json:"mid_connection_id" binding:"omitempty,gte=1"` //如果没有传，则自动生成新的中包
	MidType		string	`json:"mid_type"  binding:"omitempty,eq=normal|eq=special_a|eq=special_b"`
	Weight		float64 `json:"weight"  binding:"omitempty,gte=0"`
	Platform   	string `json:"platform"  binding:"omitempty,eq=shopee|eq=manual|eq=stock_up"`
	Note       	string `json:"note"  binding:"omitempty,lte=1024"`
}

type ListMidConnectionReqCBD struct {
	VendorID   	uint64 `json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	ConnectionID   	uint64 `json:"connection_id" form:"connection_id" binding:"omitempty,gte=1"`
	CustomsNum 	string `json:"customs_num" form:"customs_num" binding:"omitempty,lte=32"`
	MidNum 		string `json:"mid_num" form:"mid_num" binding:"omitempty,lte=32"`
	MidType		string `json:"mid_type" form:"mid_type" binding:"omitempty,eq=normal|eq=special_a|eq=special_b"`
	Status     	string `json:"status" form:"status" binding:"omitempty,eq=init|eq=stock_out|eq=customs|eq=arrive"`
	NoteKey	   	string `json:"note_key" form:"note_key" binding:"omitempty,lte=64"`

	IsPaging  	bool `json:"is_paging" form:"is_paging"`
	PageIndex 	int  `json:"page_index" form:"page_index" binding:"required"`
	PageSize  	int  `json:"page_size" form:"page_size" binding:"required"`
}

type EditMidConnectionReqCBD struct {
	ID         	uint64 		`json:"id"  binding:"required,gte=1"`
	VendorID   	uint64 		`json:"vendor_id"  binding:"required,gte=1"`
	CustomsNum 	string 		`json:"customs_num"  binding:"required,lte=32"`
	MidNum 		string 		`json:"mid_num"  binding:"required,lte=32"`
	Note       	string		`json:"note"  binding:"omitempty,lte=1024"`
	Weight		float64 	`json:"weight"  binding:"omitempty,gte=0"`

	ConnectionID   	uint64
	MdMidConn		*model.MidConnectionMD
}

type EditMidConnectionWeightReqCBD struct {
	ID         	uint64 		`json:"id"  binding:"required,gte=1"`
	VendorID   	uint64 		`json:"vendor_id"  binding:"required,gte=1"`
	Weight		float64 	`json:"weight"  binding:"omitempty,gte=0"`
}

type ChangeMidConnectionReqCBD struct {
	ID         uint64 `json:"id"  binding:"omitempty,gte=1"`
	CustomsNum string `json:"customs_num"  binding:"omitempty,lte=64"`
	VendorID   uint64 `json:"vendor_id"  binding:"required,gte=1"`
	Status     string `json:"status"  binding:"required,eq=stock_out|eq=customs|eq=arrive"`
}

type DeductMidConnectionReqCBD struct {
	ID       uint64 `json:"id"  binding:"required,gte=1"`
	VendorID uint64 `json:"vendor_id"  binding:"required,gte=1"`
}

type GetMidConnectionReqCBD struct {
	VendorID uint64 `json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	ID       uint64 `json:"id" form:"id" binding:"required,gte=1"`
}

type DelMidConnectionReqCBD struct {
	VendorID uint64 `json:"vendor_id"  binding:"required,gte=1"`
	ID       uint64 `json:"id"  binding:"required,gte=1"`

	MdMidConn *model.MidConnectionMD
}

type BatchMidConnectionOrderReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"omitempty,gte=1"`
	CustomsNum		string		`json:"customs_num"  binding:"omitempty,lte=32"`
	MidConnectionID		uint64		`json:"mid_connection_id" binding:"omitempty,gte=1"` //如果没有传，则自动生成新的中包
	MidType			string		`json:"mid_type"  binding:"omitempty,eq=normal|eq=special_a|eq=special_b"`
	Weight       		float64 	`json:"weight"  binding:"omitempty,gte=0"`
	AddKeyDetail		[]string	`json:"add_key" binding:"required"`
}
// --------------------resp-------------------------------
type ListMidConnectionRespCBD struct {
	ID             uint64                        `json:"id"  xorm:"id pk autoincr"`
	VendorID       uint64                        `json:"vendor_id"  xorm:"vendor_id"`
	MidNum	       string                        `json:"mid_num"  xorm:"mid_num"`
	MidNumCompany  string                        `json:"mid_num_company"  xorm:"mid_num_company"`
	MidType	       string                        `json:"mid_type"  xorm:"mid_type"`
	ConnectionID   uint64			     `json:"connection_id"  xorm:"connection_id"`
	CustomsNum     string                        `json:"customs_num"  xorm:"customs_num"`
	DescribeNormal     string                    `json:"describe_normal"  xorm:"describe_normal"`
	DescribeSpecial    string                    `json:"describe_special"  xorm:"describe_special"`
	Platform       string                        `json:"platform"  xorm:"platform"`
	OrderCount     int                           `json:"order_count"  xorm:"order_count"`
	//Status         string                        `json:"status"  xorm:"status"`
	//WarehouseID    uint64                        `json:"warehouse_id"  xorm:"warehouse_id"`
	//WarehouseName  string                        `json:"warehouse_name"  xorm:"warehouse_name"`
	//LineID         uint64                        `json:"line_id"  xorm:"line_id"`
	//Source         uint64                        `json:"source"  xorm:"source"`
	//To             uint64                        `json:"to"  xorm:"to"`
	//SourceName     string                        `json:"source_name"  xorm:"source_name"`
	//ToName         string                        `json:"to_name"  xorm:"to_name"`
	//SendWayID      uint64                        `json:"sendway_id"  xorm:"sendway_id"`
	//SendWayType    string                        `json:"sendway_type"  xorm:"sendway_type"`
	//SendWayName    string                        `json:"sendway_name" xorm:"sendway_name"`
	//Note           string                        `json:"note"  xorm:"note"`
	Weight         float64                       `json:"weight"  xorm:"weight"`
	CreateTime     cp_obj.Datetime               `json:"create_time"  xorm:"create_time"`
	DeductFailList map[string]*DeductFailListCBD `json:"deduct_fail_list"  xorm:"create_time"`
}

type MidConnectionInfoResp struct {
	Num			string		`json:"num"  xorm:"num"`
	NumCompany		string		`json:"num_company"  xorm:"num_company"`
	Header			string		`json:"header"  xorm:"header"`
	Invoice			string		`json:"invoice"  xorm:"invoice"`
	SendAddr		string		`json:"send_addr"  xorm:"send_addr"`
	SendName		string		`json:"send_name"  xorm:"send_name"`
	RecvName		string		`json:"recv_name"  xorm:"recv_name"`
	RecvAddr		string		`json:"recv_addr"  xorm:"recv_addr"`
	Condition		string		`json:"condition"  xorm:"condition"`
	Item			string		`json:"item"  xorm:"item"`
	Describe		string		`json:"describe"  xorm:"describe"`
	Pcs			string		`json:"pcs"  xorm:"pcs"`
	Total			string		`json:"total"  xorm:"total"`
	ProduceAddr		string		`json:"produce_addr"  xorm:"produce_addr"`
	TimeNow			int64		`json:"time_now"`
}

type GetInfoByConnectionRespCBD struct {
	MidCount		int		`json:"mid_count"  xorm:"mid_count"`
	MidWeight		float64		`json:"mid_weight"  xorm:"mid_weight"`
}
