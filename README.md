# reptile
Tencent Education

>go get github.com/zhaoshouli/reptile/

该项目必须要连接数据库使用，数据库配置文件在config目录下的config.go中配置

```golang
//配置数据库的信息
	MysqlDbInfo = MySqlInfo{
		UserName:       "root",		    //用户名
		DataBaseName:   "mytest",	    //库名
		Addr:	        "localhost",	   //数据库地址
		Password:       "123456",       //数据库密码
	}
```

dependency_pack目录为该项目的外部依赖包

logic目录中logic.go主要为项目的主要逻辑

mysql目录中conndb.go为连接数据文件包含数据的初始化连接

该项目主要设计思路为一层一层的爬取数据，先将最外层的科目信息爬取出来，再将每个科目中每个年级阶段的链接爬取出来，在讲每个年级阶段中的所有课程爬取出来，最后将每个课程的信息并发爬出来并写进数据库中，

解析页面主要用到了[goquery](https://github.com/goquery)这个包来解决的

注意！mac环境下
```
switch fmt.Sprint(runtime.GOOS) {
	//mac环境下如果编译报了errer：unknown field 'HideWindow' in struct literal of type syscall.SysProcAttr
	//这个错误就将下面这个代码注释掉没有的话就不需要注释
	//case "windows":
		//cmd := exec.Command(`cmd`, `/c`, `start`, `http://localhost:8080`)
		//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		//cmd.Start()
	case "darwin":
		exec.Command(`open`, `http://localhost:8080`).Start()
	case "linux":
		exec.Command(`xdg-open`, `http://localhost:8080`).Start()
	}
```

爬取完数据会自动打开浏览器，如果没有自动跳转则打开[localhost:8080](http://localhost:8080)
