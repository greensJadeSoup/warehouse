package cp_dc

type DcMQConfig struct {
	Alias    	string `json:"alias"`
	Type     	string `json:"type"`
	Server   	[]string `json:"server"`
}
