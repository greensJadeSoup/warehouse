package cbd

type ModelDetailCBD struct {
	ID			uint64		`json:"id"  xorm:"id pk autoincr"`
	SellerID		uint64		`json:"seller_id"  xorm:"seller_id"`
	VendorID		uint64		`json:"vendor_id"  xorm:"vendor_id"`
	StockID			uint64		`json:"stock_id"  xorm:"stock_id"`
	ShopID			uint64		`json:"shop_id"  xorm:"shop_id"`
	Platform		string		`json:"platform"  xorm:"platform"`

	PlatformItemID		string		`json:"platform_item_id"  xorm:"platform_item_id"`
	PlatformModelID		string		`json:"platform_model_id"  xorm:"platform_model_id"`

	ItemName		string		`json:"item_name"  xorm:"item_name"`
	ItemSku			string		`json:"item_sku"  xorm:"item_sku"`
	ItemStatus		string		`json:"item_status"  xorm:"item_status"`
	ItemImages		string		`json:"item_images"  xorm:"item_images"`

	ModelSku		string		`json:"model_sku"  xorm:"model_sku"`
	ModelIsDelete		uint8		`json:"model_is_delete"  xorm:"model_is_delete"`
	ModelImages		string		`json:"model_images"  xorm:"model_images"`
}
