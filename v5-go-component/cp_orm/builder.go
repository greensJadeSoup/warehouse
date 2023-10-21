package cp_orm

import (
	"errors"

	"xorm.io/builder"
)

type Builder struct {
	*builder.Builder
}

//Select 重写Builder.Select方法，总是创建新对象
func (b *Builder) Select(cols ...string) *Builder {
	// if b.Builder == nil {
	b.Builder = builder.Select(cols...)
	// } else {
	// b.Builder.Select(cols...)
	// }

	return b
}

//Insert 重写Builder.Insert方法，总是创建新对象
func (b *Builder) Insert(eq builder.Eq) *Builder {
	// if b.Builder == nil {
	b.Builder = builder.Insert(eq)
	// } else {
	// b.Builder.Insert(eq)
	// }

	return b
}

//Delete 重写Builder.Delete方法，总是创建新对象
func (b *Builder) Delete(conds ...builder.Cond) *Builder {
	// if b.Builder == nil {
	b.Builder = builder.Delete(conds...)
	// } else {
	// b.Builder.Delete(conds...)
	// }

	return b
}

func (b *Builder) getTableName(table interface{}) string {
	switch val := table.(type) {
	case ModelInterface:
		return val.TableName()
	case string:
		return val
	default:
		panic(errors.New("table 只支持 ModelInterface 接口和 string 类型"))
	}
}

//From 重写Builder.From方法，支持model传入
func (b *Builder) From(table interface{}) *Builder {
	b.Builder.From(b.getTableName(table))

	return b
}

//Into 重写Builder.Into方法，支持model传入
func (b *Builder) Into(table interface{}) *Builder {
	b.Builder.Into(b.getTableName(table))

	return b
}

//Join 重写Builder.Join方法，支持model传入
func (b *Builder) Join(joinType string, joinTable interface{}, joinCond interface{}) *Builder {
	b.Builder.Join(joinType, b.getTableName(joinTable), joinCond)

	return b
}

//InnerJoin 重写Builder.InnerJoin方法，支持model传入
func (b *Builder) InnerJoin(joinTable interface{}, joinCond interface{}) *Builder {
	b.Builder.InnerJoin(b.getTableName(joinTable), joinCond)

	return b
}

//LeftJoin 重写Builder.LeftJoin方法，支持model传入
func (b *Builder) LeftJoin(joinTable interface{}, joinCond interface{}) *Builder {
	b.Builder.LeftJoin(b.getTableName(joinTable), joinCond)

	return b
}

//RightJoin 重写Builder.RightJoin方法，支持model传入
func (b *Builder) RightJoin(joinTable interface{}, joinCond interface{}) *Builder {
	b.Builder.RightJoin(b.getTableName(joinTable), joinCond)

	return b
}

//CrossJoin 重写Builder.CrossJoin方法，支持model传入
func (b *Builder) CrossJoin(joinTable interface{}, joinCond interface{}) *Builder {
	b.Builder.CrossJoin(b.getTableName(joinTable), joinCond)

	return b
}

//FullJoin 重写Builder.FullJoin方法，支持model传入
func (b *Builder) FullJoin(joinTable interface{}, joinCond interface{}) *Builder {
	b.Builder.FullJoin(b.getTableName(joinTable), joinCond)

	return b
}
