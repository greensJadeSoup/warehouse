package cp_orm

type ModelInterface interface {
	DatabaseAlias() string
	TableName() string
}

type ModelList struct {
	NoPaging  bool        `json:"noPaging"`
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
	PageCount int         `json:"pageCount"`
	Total     int64       `json:"total"`
	Items     interface{} `json:"items"`
}
