package cp_dc

type DcNosqlConfig struct {
	Mongo struct{
		Hosts   []string 	`json:"hosts"`
	}	`json:"mongo"`
}

