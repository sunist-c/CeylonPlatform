package initialization

import (
	"CeylonPlatform/middleware/api"
	"xorm.io/xorm/names"
)

func StartUp() error {
	// 同步数据库
	DbConnection.SetMapper(names.GonicMapper{})
	for _, v := range syncEntityList {
		err := DbConnection.Sync2(v)
		if err != nil {
			return err
		}
	}

	// 绑定服务路由
	api.Bind()

	// 启动API服务
	return startFunction()
}
