MySQL数据库表生成Markdown工具
简介：
这个小工具我总共写了两次，第一次啊功能写出来了，但是我觉得自己的代码逻辑不够顺畅，感觉就是一坨粑粑在砌墙，乱涂一气~后来经过思索，改善了代码格式，作为学习两个月go的新手，我察觉到，代码怎么改，逻辑是不变的，代码重构是用来提高执行效率和人工体验的~ 

业务描述 

输入 ：

连接MySQL数据库的 地址 端口 库名

输出 ：

扫描输入库名的所有表生成如下格式的Markdown文档

#* 表名 *(*表备注*)* 

字段名 | 字段类型 | 是否为空 | 备注 

‐‐‐|‐‐‐|‐‐‐|‐‐‐ 

id | int | N | 主键 

name | varchar(20) | N | 用户名
首先 拿到一个问题 先分析需要用到的功能：

1.首先我要连接数据库 地址 端口 库名

2.输入库名，将该库名下所有表生成格式为Markdown格式的文档

	2.1 输入

	2.2 获取表

	2.3 按md格式生成

分析完需求，开始寻找第一步如何使用go语言连接数据库？

1.GORM操控数据库
经过GORM操纵数据库的学习，我们通过GORM的最新文档实现了增删改查的操作，GORM更新后是真的又灵活，又好用，又通俗易懂。

2.开始设计工具
1.获取information_schema的所有内容

2.循环遍历按格式输出。

格式：

#* 表名 *(*表备注*)* 

字段名 | 字段类型 | 是否为空 | 备注 

‐‐‐|‐‐‐|‐‐‐|‐‐‐ 

id | int | N | 主键 

name | varchar(20) | N | 用户名          
问题：

在使用GORM的时候，无一例外的都是需要一个结构体的，在这个结构体上创建table表往里面加内容，然后查询。

具体操作：

1.columns表
columns表提供了表中的列信息。详细表述了某张表的所有列以及每个列的信息。是show columns from schemaname(数据库名).tablename(表名)的结果取之此表。

2.tables表
tables表提供了关于数据库中的表的信息（包括视图）。详细表述了某个表属于哪个schema，表类型，表引擎，创建时间等信息。是show tables from db_ca_ods;【注db_ca_ods为数据库名】的结果取之此表。
方法已经出来了:

SELECT table_name,table_comment FROM tables WHERE table_schema = 'study'
table_schema就是数据库名 table_name就是表名

设计思路:

1.创建2个结构体

type tableInfo struct {
   Name string  '表名' 
   Comment string '表备注'
   //这里用语法把除information_schema,mysql,performance_chema的table_schema的其他数据库打印出来就是全部的表.
}
type tableColumn struct {
   ColumnName string '字段名'  // 
   Columntype string '字段类型' 
   IsNullable string '是否为空'
   Extra string '备注'
}
2.数据库连接(这里就得连接information_schema这个数据库)

3.然后使用GORM中的IS NOT 来输出除默认数据库的所有库中的表.格式按照tableColum结构体来输出 .

按(md格式)输出语法 :

fmt.printf("#* %s *(*%s*)*" ,Name,Comment)
fmt.printf("字段名 | 字段类型 | 是否为空 | 备注 \n‐‐‐|‐‐‐|‐‐‐|‐‐‐ \n|%s|%s|%s|%s",ColumnName,Columntype,IsNullable,Extra)
4.把输出的文字保存成md后缀的文件,保存在指定位置.

5.至于其他格式的文件,可以使用typora这个软件自动转码,方便快捷,如果想要WORD格式的,就百度搜索word格式标准,改一下输出格式就行了,不过发现还是md最好用哦,typora什么格式都能转了.

这是思路,具体写下来是这样的:

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
sqlTables="SELECT table_name AS name,table_comment AS comment from tables WHERE table_schema=?"
sqlColumns="SELECT column_name AS ColumnName,column_type AS ColumnType,is_nullable AS IsNullable,column_comment AS ColumnComment from columns WHERE table_schema=? AND table_name=?"
)
/* *
Struct for tables 
*/
type tableInfo struct {
Name string  `table_name`   //表名
Comment string `table_comment`  //表备注
}
/* *
Struct for columns
*/
type tableColumn struct {
ColumnName string `column_name` //字段名
ColumnType string `column_type` //字段类型
IsNullable string `is_nullable`//是否为空
ColumnComment string `column_comment`   //备注
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
dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/information_schema?charset=utf8mb4&parseTime=True&loc=Local",username,pwd)
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
if err != nil {
panic("failed to connect database")
}  

tableCollect :=[]tableInfo{}
columnCollect :=[]tableColumn{}
//Query tableInfo 
db.Raw(sqlTables,database).Scan(&tableCollect)                                                                  
for _, value := range tableCollect{
if value.Comment ==""{
value.Comment="无信息"
}
//Query tableColumn
db.Raw(sqlColumns,database,value.Name).Scan(&columnCollect)
//Markdown Head    
fmt.Printf("# * %v *(* %v *)* \n", value.Name,value.Comment)
fmt.Println("| 字段名 | 字段类型 | 是否为空 | 备注 |")
fmt.Println("| ------ | ------- | --------| ------ |")
for _, value1 := range columnCollect{
if value1.ColumnComment ==""{
value1.ColumnComment="无信息"
}
fmt.Printf("| %v | %v | %v | %v |\n",value1.ColumnName,value1.ColumnType,value1.IsNullable,value1.ColumnComment)
}
fmt.Println("") 
}
}
学习历程:

看前面,打印出来直接就是Markdown的格式.这里还并不完善,可以更人性化,接下来的步骤是如何让连接的数据库输出的字符保存成后缀为md的文档,把输出的内容以txt的格式存为md的文档这才算是完成.不过到这里语法已经完成了.

接下来的就是优化和讲解了,先把源代码放这里.下一步就是完结篇.真的是不要问,这个是为什么?代码都在上面,功能实现完毕.这个小工具想来真的很实用啊.......打开navicat一个个看真的很麻烦,直接打印成表,看起来舒服多了.

这个mysql_markdown的工具目前就设置了用户和密码再就是数据库这三个变量输入端口.默认的本地地址和端口都没有改动.

看起来代码量很少对吧,花了很多心思,需求不难,逻辑结构也不难,难在什么呢???????????

难在如何用go语言去实现这个功能啊!!!!!!!!!!!!!!!!!!!!!由于网上没教程,自己通过原来做的笔记,一点点根据需求弄出来了.

