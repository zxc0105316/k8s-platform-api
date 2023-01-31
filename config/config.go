package config

const (
	ListenAddr     = "0.0.0.0:9090"
	Kubeconfig     = "G:\\k8s\\config"
	PodLogTailLine = 50

	// 数据库配置    使用配置文件的方式安全性较低，容易泄漏，建议使用配置中心等工具保存配置信息
	DbUser = "root"
	DbPwd  = "Tan.7764666"
	DbHost = "localhost"
	DbPort = 3306
	DbName = "k8s-api"
	DbType = "mysql"

	// 打印Mysql debug的sql日志
	LogMode = false

	//	连接池的配置
	MaxIdleConns = 10  //最大空闲连接
	MaxOpenConns = 100 //最大连接数
	MaxLifeTime  = 30  // 最大生命周期
)
