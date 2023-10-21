package cbd

//------------------------ req ------------------------
type LoginReqCBD struct {
	Account		string		`json:"account" binding:"required,lt=50"`
	Password	string		`json:"password" binding:"required,lt=50"`
	AccountType	string		`json:"account_type" binding:"required,eq=seller|eq=manager"`
	DeviceType	string		`json:"device_type" binding:"lt=64"`
	DeviceInfo	string		`json:"device_info" binding:"lt=255"`
	IP		string		`json:"ip" binding:"lt=32"`
}


type SessionReqCBD struct {
	SessionKey		string		`json:"sessionKey" binding:"required,len=32"`
	IP			string		`json:"ip" binding:"required"`
}

//------------------------ resp ------------------------
