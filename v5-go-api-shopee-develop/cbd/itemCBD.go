package cbd

//------------------------ req ------------------------
//--------------------resp-------------------------------

type ModelBaseCBD struct {
	ModelID		string		`json:"model_id"`
	ModelSku	string		`json:"model_sku"`
	Images		string		`json:"-"`
}

type ItemModelListCBD struct {
	ID		uint64	`json:"-"`
	PlatformItemID	string	`json:"-"`
	Model		[]ModelBaseCBD
}

type ItemSimpleListCBD struct {
	ID		uint64	`json:"id" xorm:"id"`
	PlatformItemID	string	`json:"platform_item_id" xorm:"platform_item_id"`
}

type ItemBaseInfoCBD struct {
	ID			uint64		`json:"-"`
	ItemID			string		`json:"item_id"`
	CategoryID		string		`json:"category_id"`
	ItemName		string		`json:"item_name"`
	ItemStatus		string		`json:"item_status"`
	Description		string		`json:"description"`
	ItemSku			string		`json:"item_sku"`
	Weight			string		`json:"weight"`
	HasModel		bool		`json:"has_model"`
	UpdateTime		int64		`json:"update_time"`

	WeightFloat		float64		`json:"-"`
	IntHasModel		uint8		`json:"-"`
	ImageUrlList		[]string	`json:"image"`
}