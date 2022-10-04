package morph

import (
	"database/sql"
	"fmt"
)

func openSqlServer() *sql.DB {
	sqlServer, err := sql.Open(driver, dataSource) // open不会检验用户名和密码
	if err != nil {
		fmt.Println("Connect Mysql failed", err)
		return nil
	}

	_, err = sqlServer.Exec(useDatabase) // 选择数据库
	if err != nil {
		fmt.Println("Use database failed", err)
		_ = sqlServer.Close()
		return nil
	}

	return sqlServer
}

func closeSqlServer(sqlServer *sql.DB) {
	if sqlServer != nil {
		_ = sqlServer.Close()
	}
}

func reloadSqlServer(sqlServer *sql.DB) {
	if sqlServer == nil { // 数据库连接建立失败
		for i := 0; i < 5; i++ { // 循环重启连接
			sqlServer = openSqlServer()
			if sqlServer != nil { // 连接成功
				break
			}
		}
	}
}
