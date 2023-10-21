package cp_orm

import (
	"errors"
	"fmt"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"warehouse/v5-go-component/cp_util"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/yingshengtech/go-mssqldb"
	"strconv"
	"warehouse/v5-go-component/cp_dc"
	"xorm.io/core"
	"xorm.io/xorm"
)

var Engineer *Engine

type Engine struct {
	xormEngine map[string]*xorm.Engine
}

func NewEngine(dblist *[]cp_dc.DcDBConfig) (*Engine, error) {
	Engineer = &Engine{
		xormEngine: make(map[string]*xorm.Engine),
	}

	var err error
	var driverName, dataSourceName string
	var xormEngine *xorm.Engine

	for _, db := range *dblist {
		driverName, dataSourceName, err = getConnect(db.Type, db.Server, db.Port, db.User, db.Password, db.Database)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}

		xormEngine, err = xorm.NewEngine(driverName, dataSourceName)
		if err != nil {
			return nil, cp_error.NewSysError(err)
		}

		if err = xormEngine.Ping(); err != nil {
			return nil, cp_error.NewSysError(err)
		}

		// 连接池的空闲数大小,
		xormEngine.SetMaxIdleConns(20)
		// 最大打开连接数,0:无限制
		xormEngine.SetMaxOpenConns(0)

		//结构体命名与数据库一致
		xormEngine.SetMapper(core.NewCacheMapper(new(core.SameMapper)))

		Engineer.xormEngine[db.Alias] = xormEngine
		cp_log.Info(fmt.Sprintf("数据库[%s][%s] ping success.", db.Type, db.Database))
	}

	cp_util.CheckMoralCharacter()

	return Engineer, nil
}

//获取数据库引擎
func (e *Engine) Get(mi ModelInterface) (*xorm.Engine, error) {
	engine, ok := e.xormEngine[mi.DatabaseAlias()]
	if !ok {
		return nil, errors.New("数据库引擎：" + mi.DatabaseAlias() + "不存在")
	}

	return engine, nil
}

// 获取数据库连接信息
func getConnect(typeDB, server string, port int, user, password, database string) (string, string, error) {
	var err error
	driverName := ""
	dataSourceName := ""
	switch typeDB {
	case "mssql":
		driverName = "mssql"
		dataSourceName += "server=" + server
		dataSourceName += ";port=" + strconv.Itoa(port)
		dataSourceName += ";database=" + database
		dataSourceName += ";user id=" + user
		dataSourceName += ";password=" + password
		//	case "odbc":
		//		//lunny mssql driver
		//		driverName = "odbc"
		//		dataSourceName = "driver={SQL Server}"
		//		dataSourceName += ";Server=" + db["server"] + "," + db["port"]
		//		dataSourceName += ";Database=" + db["database"]
		//		dataSourceName += ";uid=" + db["user"] + ";pwd=" + db["password"] + ";"
	case "mysql":
		driverName = "mysql"
		dataSourceName = user + ":" + password
		dataSourceName += "@(" + server + ":" + strconv.Itoa(port) + ")/"
		dataSourceName += database + "?charset=utf8mb4&multiStatements=true"
	case "sqlite3":
		driverName = "sqlite3"
		dataSourceName = database
	default:
		err = errors.New("不支持的数据库类型：" + typeDB)
	}

	return driverName, dataSourceName, err
}

