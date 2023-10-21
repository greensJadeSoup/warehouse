package cp_dc


type DcTraceLog struct {
	OnOff 		bool	`json:"onoff"`
	ActionTopic	string	`json:"actionTopic"`
	RuntimeTopic	string	`json:"runtimeTopic"`
}
