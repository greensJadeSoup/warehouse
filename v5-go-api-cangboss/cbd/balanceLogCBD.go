package cbd

import "warehouse/v5-go-component/cp_obj"

//------------------------ req ------------------------
type AddBalanceLogReqCBD struct {
	VendorID		uint64		`json:"vendor_id" binding:"omitempty,gte=1"`
	UserType		string		`json:"user_type"  binding:"required,lte=16"`
	UserID			uint64		`json:"user_id"  binding:"required,gte=1"`
	UserName		string		`json:"user_name"  binding:"required,lte=32"`
	ManagerID		uint64
	ManagerName		string
	EventType		string		`json:"event_type"  binding:"required,lte=32"`
	Change			float64		`json:"change"  binding:""`
	Balance			float64		`json:"balance"  binding:""`
	Status			string		`json:"status"  binding:""`
	Content			string		`json:"content"  binding:""`
	ObjectType		string		`json:"object_type"  binding:""`
	ObjectID		string		`json:"object_id"  binding:""`
	PriDetail		string		`json:"pri_detail"  binding:""`
	ToUser			uint64		`json:"to_user"  binding:""`
	Note			string		`json:"note"  binding:""`

}

type ListBalanceLogReqCBD struct {
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	SellerID		uint64		`json:"seller_id"  form:"seller_id"  binding:"omitempty,gte=1"`
	UserType		string		`json:"user_type" form:"user_type" binding:"omitempty,eq=seller|eq=super_manager"`
	UserID			uint64		`json:"user_id" form:"user_id" binding:"omitempty,gte=1"`
	SellerKey		string		`json:"seller_key" form:"seller_key" binding:"omitempty,lte=32"`
	EventType		string		`json:"event_type" form:"event_type" binding:"omitempty,lte=32"`
	Status			string		`json:"status" form:"status" binding:"omitempty,eq=success|eq=fail"`

	IsPaging		bool		`json:"is_paging" form:"is_paging"`
	PageIndex		int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int		`json:"page_size" form:"page_size" binding:"required"`
}

//--------------------resp-------------------------------
type ListBalanceLogRespCBD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	UserType		string		`json:"user_type"  xorm:"user_type"`
	UserID			uint64		`json:"user_id"  xorm:"user_id"`
	SellerName		string		`json:"seller_name"  xorm:"seller_name"`
	ManagerID		uint64		`json:"manager_id"  xorm:"manager_id"`
	ManagerName		string		`json:"manager_name"  xorm:"manager_name"`
	EventType		string		`json:"event_type"  xorm:"event_type"`
	Change			float64		`json:"change"  xorm:"change"`
	Balance			float64		`json:"balance"  xorm:"balance"`
	Status			string		`json:"status"  xorm:"status"`
	Content			string		`json:"content"  xorm:"content"`
	PriDetail		string		`json:"pri_detail"  xorm:"pri_detail"`
	ToUser			uint64		`json:"to_user"  xorm:"to_user"`
	Note			string		`json:"note"  xorm:"note"`
	CreateTime		cp_obj.Datetime `json:"create_time"  xorm:"create_time"`
}
