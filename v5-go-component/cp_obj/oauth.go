package cp_obj


//easyjson:json
type OauthAccessTokenDetail struct {
	IP		string		`json:"ip"`
	Memo		string		`json:"app_secret"`
	Terminal	string		`json:"terminal"`
	CreateTime	Datetime	`json:"create_time"`
}
