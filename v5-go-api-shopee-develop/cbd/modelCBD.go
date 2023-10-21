package cbd

//------------------------ req ------------------------

//--------------------resp-------------------------------
type ModelBaseInfoCBD struct {
	ModelID		string		`json:"model_id"`
	ModelSku	string		`json:"model_sku"`
	Images		string		`json:"-"`
}
