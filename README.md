# 使用方法

前提：安装并配置好golang，并开启GO111MODULE=on

1.将Mysql_Markdown下载至目录。

2.用命令行定位至目录，输入命令：

```go
go mod tidy
```

下载mod中需要的依赖包。

```go
go run main.go 
```

3.使用手册

```cmd
G:\传输文件\Mysql_Markdown>go run main.go
数据库默认的host:127.0.0.1 port:3306 charset=utf8mb4
请输入要连接的mysql数据库的账户(默认: root):
root
请输入要连接的mysql数据库的密码(默认: 123456):
123456
请输入要使用的数据库:
study
请输入转码后存放文件的地址(例:C:/jojo):
G:\传输文件\Mysql_Markdown
保存成功,请到G:\传输文件\Mysql_Markdown地址处查看文件
```

