package cp_dc

type DcCacheConfig struct {
	Alias    string `json:"alias"`
	Type     string `json:"type"`
	Server   string `json:"server"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

