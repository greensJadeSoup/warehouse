package cbd

//------------------------ req ------------------------
type AddMidConnectionNormalReqCBD struct {
	VendorID		uint64		`json:"vendor_id"  binding:"required,gte=1"`
	Num			string		`json:"num"  binding:""`
	Header			string		`json:"header"  binding:""`
	Invoice			string		`json:"invoice"  binding:""`
	SendAddr		string		`json:"send_addr"  binding:""`
	SendName		string		`json:"send_name"  binding:""`
	RecvName		string		`json:"recv_name"  binding:""`
	RecvAddr		string		`json:"recv_addr"  binding:""`
	Condition		string		`json:"condition"  binding:""`
	Item			string		`json:"item"  binding:""`
	Describe		string		`json:"describe"  binding:""`
	Pcs			string		`json:"pcs"  binding:""`
	Total			string		`json:"total"  binding:""`
	ProduceAddr		string		`json:"produce_addr"  binding:""`
}

type ListMidConnectionNormalReqCBD struct {
	ID			uint64		`json:"id" form:"id" binding:"required,gte=1"`
	VendorID		uint64		`json:"vendor_id" form:"vendor_id" binding:"required,gte=1"`
	IsPaging		bool		`json:"is_paging" form:"is_paging"`
	PageIndex		int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize		int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditMidConnectionNormalReqCBD struct {
	ID		uint64		`json:"id"  binding:"required,gte=1"`
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	Num		string		`json:"num"  binding:""`
	Header		string		`json:"header"  binding:""`
	Invoice		string		`json:"invoice"  binding:""`
	SendAddr	string		`json:"send_addr"  binding:""`
	SendName	string		`json:"send_name"  binding:""`
	RecvName	string		`json:"recv_name"  binding:""`
	RecvAddr	string		`json:"recv_addr"  binding:""`
	Condition	string		`json:"condition"  binding:""`
	Item		string		`json:"item"  binding:""`
	Describe	string		`json:"describe"  binding:""`
	Pcs		string		`json:"pcs"  binding:""`
	Total		string		`json:"total"  binding:""`
	ProduceAddr	string		`json:"produce_addr"  binding:""`
}

type DelMidConnectionNormalReqCBD struct {
	ID		uint64		`json:"id"  binding:"required,gte=1"`
}

//--------------------resp-------------------------------
type ListMidConnectionNormalRespCBD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	Num			string		`json:"num"  xorm:"num"`
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
}
