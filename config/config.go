package config

import "os"

type MySqlInfo struct {
	UserName 		string
	DataBaseName	string
	Addr			string
	Password		string
}

var (
	MysqlDbInfo = MySqlInfo{}
	HTMLAddr string
)


func init()  {
	HTMLAddr, _ = os.Getwd()

	//配置数据库的信息
	MysqlDbInfo = MySqlInfo{
		UserName: "root",		//用户名
		DataBaseName: "mytest",	//库名
		Addr:	"localhost",			//数据库地址
		Password: "123456",		//数据库密码
	}
}


