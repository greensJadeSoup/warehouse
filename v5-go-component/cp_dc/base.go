package cp_dc

type DcBaseConfig struct {
	IP		string	`json:"bind"` 		//绑定ip
	HttpPort	int	`json:"httpPort"`
	ServingLimit    int	`json:"servingLimit"`
	IsLeader	bool	`json:"isLeader"`
	IsLocal		bool	`json:"isLocal"`
	IsNeedSign	bool	`json:"isNeedSign"`
	IsTest		bool	`json:"isTest"`
}

