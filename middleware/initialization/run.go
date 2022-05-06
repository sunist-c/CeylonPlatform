package initialization

import (
	"CeylonPlatform/middleware/api"
	"CeylonPlatform/middleware/database"
	"xorm.io/xorm/names"
)

func StartUp() (err error) {
	// 同步数据库
	DbConnection.SetMapper(names.GonicMapper{})
	if err = database.Sync(DbConnection); err != nil {
		return err
	}

	// 绑定服务路由
	api.Bind()
	api.LoadHtml()

	return nil
}

func Serve() (err error) {
	// 启动API服务
	return starFunction()
}
