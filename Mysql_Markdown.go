package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/* *
Structed Query language
*/
const (
	sqlTables  = "SELECT table_name AS name,table_comment AS comment from tables WHERE table_schema=?"
	sqlColumns = "SELECT column_name AS ColumnName,column_type AS ColumnType,is_nullable AS IsNullable,column_comment AS ColumnComment from columns WHERE table_schema=? AND table_name=?"
)

/* *
Struct for tables
*/
type tableInfo struct {
	Name    string `table_name`    //表名
	Comment string `table_comment` //表备注
}

/* *
Struct for columns
*/
type tableColumn struct {
	ColumnName    string `column_name`    //字段名
	ColumnType    string `column_type`    //字段类型
	IsNullable    string `is_nullable`    //是否为空
	ColumnComment string `column_comment` //备注
}

/* *
main func
*/
func main() {
	//Query Mysql_User_Message
	var username string
	var pwd string
	var database string
	fmt.Println("数据库默认的hose:127.0.0.1 port:3306 charset=utf8mb4")
	fmt.Println("请输入要连接的mysql数据库的账户(默认: root):")
	fmt.Scan(&username)
	fmt.Println("请输入要连接的mysql数据库的密码(默认: 123456):")
	fmt.Scan(&pwd)
	fmt.Println("请输入要使用的数据库:")
	fmt.Scan(&database)
	//Connect Mysql Service
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/information_schema?charset=utf8mb4&parseTime=True&loc=Local", username, pwd)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	tableCollect := []tableInfo{}
	columnCollect := []tableColumn{}
	//Query tableInfo
	db.Raw(sqlTables, database).Scan(&tableCollect)
	for _, value := range tableCollect {
		if value.Comment == "" {
			value.Comment = "无信息"
		}
		//Query tableColumn
		db.Raw(sqlColumns, database, value.Name).Scan(&columnCollect)
		//Markdown Head
		fmt.Printf("# * %v *(* %v *)* \n", value.Name, value.Comment)
		fmt.Println("| 字段名 | 字段类型 | 是否为空 | 备注 |")
		fmt.Println("| ------ | ------- | --------| ------ |")
		for _, value1 := range columnCollect {
			if value1.ColumnComment == "" {
				value1.ColumnComment = "无信息"
			}
			fmt.Printf("| %v | %v | %v | %v |\n", value1.ColumnName, value1.ColumnType, value1.IsNullable, value1.ColumnComment)
		}
		fmt.Println("")
	}
}
