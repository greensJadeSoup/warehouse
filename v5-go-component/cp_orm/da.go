package cp_orm

import (
	"fmt"
	"strconv"
	"warehouse/v5-go-component/cp_error"
	"warehouse/v5-go-component/cp_log"
	"xorm.io/xorm"
)

type DA struct {
	*xorm.Session

	NotComm		bool	//如果为true，这表示目前正在一个事务的执行过程中，不允许Commit
	Transacting 	bool	//如果为true，这表示目前正在一个事务的执行过程中，不允许Close

	Builder   	*Builder
	model     	ModelInterface
	Engine    	*xorm.Engine
}

type DAInterface interface {
	Init(*xorm.Engine, ModelInterface)
	Begin() error
	Rollback() error
	Commit() error
	NotCommit()
	AllowCommit() *DA
}

// 初始化DA对象
func InitDA(da DAInterface, model ModelInterface) error {
	engine, err := Engineer.Get(model)
	if err != nil {
		return err
	}

	da.Init(engine, model)
	return nil
}

func (da *DA) NotCommit() {
	da.NotComm = true
}

func (da *DA) AllowCommit() *DA {
	da.NotComm = false
	return da
}

func (da *DA) Begin() error {
	err := da.Session.Begin()
	if err != nil {
		return err
	}
	da.Transacting = true //开启事务标志，当close()遇到该标志，不执行close动作
	return nil
}

func (da *DA) Rollback() error {
	da.Transacting = false //关启事务标志，允许close()
	err := da.Session.Rollback()
	if err != nil {
		return err
	}
	return nil
}

func (da *DA) Commit() error {
	if da.NotComm {
		return nil
	}

	da.Transacting = false //关启事务标志，允许close()
	err := da.Session.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (da *DA) Close() {
	if da.Transacting {
		return
	}
	da.Session.Close()
}

func (da *DA) GetModel() ModelInterface {
	return da.model
}

func (da *DA) Init(engine *xorm.Engine, model ModelInterface) {
	da.Engine = engine
	da.model = model

	//如果不为nil，则表示已经继承会话，不需要开启新的会话
	if da.Session == nil {
		da.Session = engine.NewSession()
	}
}

//添加记录
func (da *DA) Insert(model interface{}) (int64, error) {
	return da.Session.Insert(model)
}

//删除记录
func (da *DA) Delete(model interface{}) (int64, error) {
	return da.Session.Delete(model)
}

//查询单条记录
func (da *DA) Get(model interface{}) (bool, error) {
	return da.Session.Get(model)
}

//同xorm.Engine.In
func (da *DA) In(column string, args ...interface{}) *DA {
	da.Session.In(column, args...)
	return da
}

//同xorm.Engine.Where
func (da *DA) Where(querystring string, args ...interface{}) *DA {
	da.Session.Where(querystring, args...)
	return da
}

//同xorm.Engine.And
func (da *DA) And(querystring string, args ...interface{}) *DA {
	da.Session.And(querystring, args...)
	return da
}

//同xorm.Engine.Or
func (da *DA) Or(querystring string, args ...interface{}) *DA {
	da.Session.Or(querystring, args...)
	return da
}

//同xorm.Engine.Desc
func (da *DA) Desc(colNames ...string) *DA {
	da.Session.Desc(colNames...)
	return da
}

//同xorm.Engine.Asc
func (da *DA) Asc(colNames ...string) *DA {
	da.Session.Asc(colNames...)
	return da
}

//同xorm.Engine.Cols
func (da *DA) Cols(colNames ...string) *DA {
	da.Session.Cols(colNames...)
	return da
}

func (da *DA)MysqlModelList(listSQL string, isPaging bool, pageIndex, pageSize int, item... interface{}) (*ModelList, error) {
	ml := &ModelList{
		NoPaging: !isPaging,
		Items: make([]interface{}, 0),
	}

	if pageIndex == 0 {
		pageIndex = 1
	}

	if pageSize == 0 {
		pageSize = 10
	}

	if len(item) > 0 {
		ml.Items = item[0]
	}

	var searchSQL string
	if isPaging {
		searchSQL = listSQL + " limit " + strconv.Itoa(pageSize) + " offset " + strconv.Itoa((pageIndex-1)*pageSize)
	} else {
		searchSQL = listSQL
	}

	cp_log.Debug(searchSQL)
	err := da.SQL(searchSQL).Find(ml.Items)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	countSQL := "select count(0) from (" + listSQL + ")tc"

	cp_log.Debug(countSQL)
	_, err = da.SQL(countSQL).Get(&ml.Total)
	if err != nil {
		return nil, cp_error.NewSysError(err)
	}

	if isPaging {
		ml.PageSize = pageSize
		ml.PageIndex = pageIndex
		ml.PageCount = int(ml.Total) / ml.PageSize
		if int(ml.Total) % ml.PageSize > 0 {
			ml.PageCount += 1
		}
	} else {
		ml.PageSize = int(ml.Total)
		ml.PageIndex = 1
		ml.PageCount = 1
	}

	return ml, nil
}

func (da *DA) DeferHandle(err *error) {
	if r := recover(); r != nil || *err != nil {
		da.Rollback()

		if r != nil {
			*err = cp_error.NewSysError(fmt.Sprintf("%v", r)) //需要赋值给err，否则前端无报错
		}
	}
}