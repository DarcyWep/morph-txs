package morph

import (
	"database/sql"
	"fmt"
)

func OpenSqlServers() error {
	for range tables {
		sqlServer := openSqlServer()
		if sqlServer == nil {
			for _, ss := range sqlServers {
				_ = ss.Close()
			}
			return fmt.Errorf("Open sql server failed")
		}
		sqlServers = append(sqlServers, sqlServer)
	}
	return nil
}

func openSqlServer() *sql.DB {
	sqlServer, err := sql.Open(driver, dataSource) // open不会检验用户名和密码
	if err != nil {
		fmt.Println("Connect Mysql failed", err)
		return nil
	}

	_, err = sqlServer.Exec(useDatabase) // 选择数据库
	if err != nil {
		fmt.Println("Error: Use database failed", err)
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

func CloseSqlServers() {
	for _, ss := range sqlServers {
		closeSqlServer(ss)
	}
}
