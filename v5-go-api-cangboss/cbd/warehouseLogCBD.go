package cbd

import "warehouse/v5-go-component/cp_obj"

//------------------------ req ------------------------
type AddWarehouseLogReqCBD struct {
	VendorID		uint64		`json:"vendor_id"  binding:"required,gte=1"`
	UserType		string		`json:"user_type"  binding:"required,lte=16"`
	UserID			uint64		`json:"user_id"  binding:"required,gte=1"`
	RealName		string		`json:"real_name"  binding:"required,lte=32"`
	EventType		string		`json:"event_type"  binding:"required,lte=16"`
	WarehouseID		uint64		`json:"warehouse_id"  binding:"required,gte=1"`
	WarehouseName		string		`json:"warehouse_name"  binding:"required,lte=32"`
	ObjectType		string		`json:"object_type"  binding:"required,lte=16"`
	ObjectID		string		`json:"object_id"  binding:"required,gte=1"`
	Content			string		`json:"content" binding:"required,lte=255"`
}

type ListWarehouseLogReqCBD struct {
	VendorID	uint64		`json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	WarehouseID 	uint64		`json:"warehouse_id" form:"warehouse_id" binding:"gte=1"`
	EventType	string		`json:"event_type" form:"event_type" binding:"lte=16"`
	ObjectType	string		`json:"object_type"  binding:"omitempty,lte=16"`
	ObjectID	string		`json:"object_id"  binding:"omitempty,lte=255"`
	WarehouseIDList []string

	IsPaging	bool     	`json:"is_paging" form:"is_paging"`
	PageIndex	int      	`json:"page_index" form:"page_index" binding:"required"`
	PageSize	int      	`json:"page_size" form:"page_size" binding:"required"`
}

type ListWarehouseLogByObjIDListReqCBD struct {
	UserType		string		`json:"user_type"  binding:"required,lte=16"`
	ObjectType		string		`json:"object_type" form:"object_type" binding:"required,lte=16"`
	ObjectID		[]string	`json:"object_id" form:"object_id" binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type ListWarehouseLogRespCBD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	WarehouseID		uint64		`json:"warehouse_id"  xorm:"warehouse_id"`
	WarehouseName		string		`json:"warehouse_name"  xorm:"warehouse_name"`
	UserType		string		`json:"user_type"  xorm:"user_type"`
	UserID			uint64		`json:"user_id"  xorm:"user_id"`
	RealName		string		`json:"real_name"  xorm:"real_name"`
	EventType		string		`json:"event_type"  xorm:"event_type"`
	ObjectType		string		`json:"object_type"  xorm:"object_type"`
	ObjectID		string		`json:"object_id"  xorm:"object_id"`
	Content			string		`json:"content"  xorm:"content"`
	CreateTime		cp_obj.Datetime `json:"create_time"  xorm:"create_time"`
}
