package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wonderivan/logger"
	"k8s-platform-api/config"
	"time"
)

var (
	isInit bool
	GORM   *gorm.DB
	err    error
)

func Init() {

	//	 判断是否已经初始化
	if isInit {
		return
	}

	//	 组装连接的配置
	//   parseTime是查询结果是否自动解析为时间
	//	 loc是mysql的时区设置

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DbUser,
		config.DbPwd,
		config.DbHost,
		config.DbPort,
		config.DbName)

	fmt.Println(dsn)
	// 建立数据库连接，生成一个*gorm.DB类型的对象
	GORM, err := gorm.Open(config.DbType, dsn)

	if err != nil {
		panic("数据库连接失败" + err.Error())
	}
	//  打印sql语句
	GORM.LogMode(config.LogMode)

	//	开启连接池
	GORM.DB().SetConnMaxLifetime(time.Duration(config.MaxLifeTime))
	GORM.DB().SetMaxOpenConns(config.MaxOpenConns)
	GORM.DB().SetMaxIdleConns(config.MaxIdleConns)

	isInit = true

	logger.Info("数据库初始化成功")

}

// 关闭数据库连接
func Close() error {

	return GORM.Close()

}
