package main

import (
	"gorm.io/driver/mysql"
)

func main() {
	mysql.Open("root:199107lmr1636@tcp(123.57.220.143:3306)")
}
