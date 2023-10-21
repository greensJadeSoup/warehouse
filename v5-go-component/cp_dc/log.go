package cp_dc

type DcLogItemConfig struct {
	Enable   	bool		`json:"enable"`
	Path	   	string		`json:"filePath"`
	Level		string		`json:"level"`
}

type DcLogConfig struct {
	ConsoleWriter   	DcLogItemConfig		`json:"consoleWriter"`
	InfoFileWriter    	DcLogItemConfig		`json:"infoFileWriter"`
	ErrorFileWriter		DcLogItemConfig		`json:"errorFileWriter"`
}
