package dbconfig

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Username      string `gorm:"unique"`
	Password      string `gorm:"check: length(password) > 4"`
	FollowCount   int64  `gorm:"check: follow_count >= 0"`
	FollowerCount int64  `gorm:"check: follower_count >= 0"`
}

type mysqlConn struct {
	user      string
	pwd       string
	protocol  string
	address   string
	port      string
	dbname    string
	charset   string
	parseTime string
}

func (c *mysqlConn) GetDSN() string {
	return fmt.Sprintf("%v:%v@%v(%v:%v)/%v?charset=%v&parseTime=%v",
		c.user, c.pwd, c.protocol, c.address, c.port, c.dbname, c.charset, c.parseTime)
}

var DB *gorm.DB

func InitDB() {
	dsn := (&mysqlConn{
		user:      "tkadmin",
		pwd:       "123456",
		protocol:  "tcp",
		address:   "127.0.0.1",
		dbname:    "tiktok",
		port:      "3306",
		charset:   "utf8mb4",
		parseTime: "True",
	}).GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
}