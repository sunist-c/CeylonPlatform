package database

import (
	"xorm.io/xorm"
)

var (
	syncEntityList []interface{} = make([]interface{}, 0, 16)
)

// Sync 与数据库同步表结构
func Sync(dbConn *xorm.Engine) error {
	return dbConn.Sync2(syncEntityList)
}

// AddEntity 向数据库表结构同步列表添加一个实体
func AddEntity(interfaces ...interface{}) {
	syncEntityList = append(syncEntityList, interfaces)
}

// ==================================== //
// ========这下面的内容暂时没有实现======== //
// ==================================== //

// IDatabase DAO抽象接口，暂时没有实现的打算
type IDatabase interface {
	Sync(objects []interface{}) (err error)
	Connect() (err error)
	Check() (ok bool, err error)
	Close() (err error)
	Query(object interface{}) (ok bool, err error)
	Update(oldObject, newObject interface{}) (ok bool, err error)
	Create(object interface{}) (ok bool, err error)
	Delete(object interface{}) (ok bool, err error)
}

type Mysql struct {
	dbConn *xorm.Engine
}

func (m Mysql) Sync(objects []interface{}) (err error) {
	return m.dbConn.Sync2(objects)
}

func (m Mysql) Connect() (err error) {
	//TODO implement me
	panic("implement me")
}

func (m Mysql) Check() (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (m Mysql) Close() (err error) {
	//TODO implement me
	panic("implement me")
}

func (m Mysql) Query(object interface{}) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (m Mysql) Update(oldObject, newObject interface{}) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (m Mysql) Create(object interface{}) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (m Mysql) Delete(object interface{}) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

type Redis struct {
}

func (r Redis) Sync(objects []interface{}) (err error) {
	//TODO implement me
	panic("implement me")
}

func (r Redis) Connect() (err error) {
	//TODO implement me
	panic("implement me")
}

func (r Redis) Check() (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r Redis) Close() (err error) {
	//TODO implement me
	panic("implement me")
}

func (r Redis) Query(object interface{}) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r Redis) Update(oldObject, newObject interface{}) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r Redis) Create(object interface{}) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r Redis) Delete(object interface{}) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

type File struct {
}

func (f File) Sync(objects []interface{}) (err error) {
	//TODO implement me
	panic("implement me")
}

func (f File) Connect() (err error) {
	//TODO implement me
	panic("implement me")
}

func (f File) Check() (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (f File) Close() (err error) {
	//TODO implement me
	panic("implement me")
}

func (f File) Query(object interface{}) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (f File) Update(oldObject, newObject interface{}) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (f File) Create(object interface{}) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (f File) Delete(object interface{}) (ok bool, err error) {
	//TODO implement me
	panic("implement me")
}
