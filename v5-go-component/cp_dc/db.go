package cp_dc

type DcDBConfig struct {
	Alias    string `json:"alias"`
	Type     string `json:"type"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}
