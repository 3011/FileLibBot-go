# FileLibBot-go

文件存储Bot，利用TG群组转存文件，所有文件存储在TG服务器上，可利用链接分享文件。

## 使用方法
### 1. 创建 config.yml

```yml
bottoken: xxxxxxxxxx:AAGxxxxxxxxxxxxxxx # bot的token
forwardgroupid: -xxxxxxxx # 转发至群组的id
dbfilename: db.db # 数据库文件名 可选 默认 FileLibBot.db
```
### 2.运行 FileLibBot
```shell
./FileLibBot
```