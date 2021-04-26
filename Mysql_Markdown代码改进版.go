package main

import (
	"bufio"
	"fmt"
	"os"

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
Struct for UserMessage
*/
type Choose struct {
	username string
	pwd      string
	database string
}

var choose Choose

/* *
QueryUserMessage
*/
func (c *Choose) QueryUserMessage() {

	fmt.Println("数据库默认的host:127.0.0.1 port:3306 charset=utf8mb4")
	fmt.Println("请输入要连接的mysql数据库的账户(默认: root):")
	fmt.Scan(&choose.username)

	fmt.Println("请输入要连接的mysql数据库的密码(默认: 123456):")
	fmt.Scan(&choose.pwd)

	fmt.Println("请输入要使用的数据库:")
	fmt.Scan(&choose.database)

}

/* *
Connect Mysql
*/
func Connect() *gorm.DB {
	//Connect Mysql Service
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/information_schema?charset=utf8mb4&parseTime=True&loc=Local", choose.username, choose.pwd)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return db
	}
	return db
}

/* *
main func
*/
func main() {
	choose.QueryUserMessage()
	db := Connect()
	tableCollect := []tableInfo{}
	columnCollect := []tableColumn{}
	// Creat newFile AND Read
	fmt.Println("请输入转码后存放文件的地址(例:C:/jojo):")
	var Path string
	fmt.Scan(&Path)
	filePath := fmt.Sprintf("%v/%v.md ", Path, choose.database)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("创建文件错误= %v \n", err)
		return
	}
	writer := bufio.NewWriter(file)
	defer file.Close()

	//Mysql_Markdown
	db.Raw(sqlTables, choose.database).Scan(&tableCollect)
	for _, value := range tableCollect {
		if value.Comment == "" {
			value.Comment = "无信息"
		}
		//Query tableColumn
		db.Raw(sqlColumns, choose.database, value.Name).Scan(&columnCollect)
		//Markdown Head
		writer.WriteString(fmt.Sprintf("# * %v *(* %v *)* \n", value.Name, value.Comment))
		writer.WriteString(fmt.Sprintln("| 字段名 | 字段类型 | 是否为空 | 备注 |"))
		writer.WriteString(fmt.Sprintln("| ------ | ------- | --------| ------ |"))
		for _, value1 := range columnCollect {
			if value1.ColumnComment == "" {
				value1.ColumnComment = "无信息"
			}
			writer.WriteString(fmt.Sprintf("| %v | %v | %v | %v |\n", value1.ColumnName, value1.ColumnType, value1.IsNullable, value1.ColumnComment))
		}
		writer.WriteString(fmt.Sprintln(""))
		writer.Flush()
	}
	fmt.Printf("保存成功,请到%v地址处查看文件", Path)
}
