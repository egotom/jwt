package initializers
import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnDb(){
	var err error
	dsn := Config.DSN
  	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("连接数据库失败！")
	}
}