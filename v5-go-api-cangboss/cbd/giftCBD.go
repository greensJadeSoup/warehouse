package cbd

//------------------------ req ------------------------
type AddGiftReqCBD struct {
	SellerID		uint64		`json:"seller_id"  binding:"required,gte=1"`
	SourceShopID		uint64		`json:"source_shop_id"  xorm:"source_shop_id"`
	SourcePlatformShopID	string		`json:"source_platform_shop_id"  xorm:"source_platform_shop_id"`
	SourceItemID		uint64		`json:"source_item_id"  xorm:"source_item_id"`
	SourcePlatformItemID	string		`json:"source_platform_item_id"  xorm:"source_platform_item_id"`
	SourceModelID		uint64		`json:"source_model_id"  xorm:"source_model_id"`
	SourcePlatformModelID	string		`json:"source_platform_model_id"  xorm:"source_platform_model_id"`
	ToShopID		uint64		`json:"to_shop_id"  xorm:"to_shop_id"`
	ToPlatformShopID	string		`json:"to_platform_shop_id"  xorm:"to_platform_shop_id"`
	ToItemID		uint64		`json:"to_item_id"  xorm:"to_item_id"`
	ToPlatformItemID	string		`json:"to_platform_item_id"  xorm:"to_platform_item_id"`
	ToModelID		uint64		`json:"to_model_id"  xorm:"to_model_id"`
	ToPlatformModelID	string		`json:"to_platform_model_id"  xorm:"to_platform_model_id"`
}

type ListGiftReqCBD struct {
	VendorID	uint64		`json:"vendor_id" form:"vendor_id" binding:"omitempty,gte=1"`
	SellerID	uint64		`json:"seller_id" form:"seller_id" binding:"omitempty,gte=1"`
	WarehouseID	uint64		`json:"warehouse_id" form:"warehouse_id" binding:"omitempty,gte=1"`
	ModelIDList	string		`json:"model_id_list" form:"model_id_list"  binding:"required,gte=1"`
	ModelIDStrList	[]string

	IsPaging	bool		`json:"is_paging" form:"is_paging"`
	PageIndex	int		`json:"page_index" form:"page_index" binding:"required"`
	PageSize	int		`json:"page_size" form:"page_size" binding:"required"`
}

type EditGiftReqCBD struct {
	ID		uint64		`json:"id"  binding:"required,gte=1"`
	SellerID	uint64		`json:"seller_id"  binding:"required,gte=1"`
	ModelID		uint64		`json:"model_id"  binding:"required,gte=1"`
	VendorID	uint64		`json:"vendor_id"  binding:"required,gte=1"`
	WarehouseID	uint64		`json:"warehouse_id"  binding:"required,gte=1"`
	StockID		uint64		`json:"stock_id"  binding:"required,gte=1"`
}

type DelGiftReqCBD struct {
	ID		uint64		`json:"id"  binding:"required,gte=1"`
}

type BindGiftReqCBD struct {
	SellerID	uint64		`json:"seller_id" binding:"required,gte=1"`
	ModelID		uint64		`json:"model_id,string" binding:"required,gte=1"` 	//源
	ModelIDList	[]string	`json:"model_id_list" binding:"required,gte=1"`		//目的
}

type SetAutoImportCBD struct {
	SellerID	uint64		`json:"seller_id" binding:"required,gte=1"`
	ModelID		uint64		`json:"model_id,string" binding:"required,gte=1"` 	//源
	AutoImport	bool		`json:"auto_import" binding:"omitempty"` 		//预报的时候是否自动导入
}

type UnBindGiftReqCBD struct {
	SellerID	uint64		`json:"seller_id" binding:"required,gte=1"`
	ModelID		uint64		`json:"model_id,string" binding:"required,gte=1"` 	//源
	ModelIDList	[]string	`json:"model_id_list" binding:"required,gte=1"`		//目的
}

//--------------------resp-------------------------------
type ListGiftRespCBD struct {
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	DependID		uint64		`json:"depend_id,string"  xorm:"source_model_id"`
	ItemID			uint64		`json:"item_id,string"  xorm:"item_id"`
	ModelID			uint64		`json:"model_id,string"  xorm:"to_model_id"`
	Platform		string		`json:"platform"  xorm:"platform"`
	ShopID			uint64		`json:"shop_id"  xorm:"shop_id"`
	PlatformShopID		string		`json:"platform_shop_id"  xorm:"platform_shop_id"`
	PlatformItemID		string		`json:"platform_item_id"  xorm:"platform_item_id"`
	PlatformModelID		string		`json:"platform_model_id"  xorm:"platform_model_id"`
	ShopName		string		`json:"shop_name"  xorm:"shop_name"`
	ItemName		string		`json:"item_name"  xorm:"item_name"`
	ItemStatus		string		`json:"item_status"  xorm:"item_status"`
	ModelSku		string		`json:"model_sku"  xorm:"model_sku"`
	ModelIsDelete		uint8		`json:"model_is_delete"  xorm:"model_is_delete"`
	Images			string		`json:"model_images"  xorm:"model_images"`
	StockID			uint64		`json:"stock_id,string"  xorm:"stock_id"`
	Total			int		`json:"total"  xorm:"total"`
	Freeze			int		`json:"freeze"  xorm:"freeze"`

	RackDetail		[]RackDetailCBD `json:"rack_detail" xorm:"rack_detail"`
}
