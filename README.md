# FileLibBot-go

文件存储 Bot，利用 TG 群组转存文件，所有文件存储在 TG 服务器上，可利用链接分享文件。

## 原理
- Bot 接收到文件时，会将文件转发至指定群组，会得到该转发消息在指定群组的 message_id。
- 获取文件时，Bot 查找文件对应的 message_id，通过指定群组转发文件消息。

## 使用方法
### 1. 创建 config.yml

```yml
bottoken: xxxxxxxxxx:AAGxxxxxxxxxxxxxxx # bot的token
forwardgroupid: -xxxxxxxx # 转发至群组的id (Bot在该群需有收发消息权限)
dbfilename: db.db # 数据库文件名 可选 默认 FileLibBot.db
```
### 2.运行 FileLibBot
```shell
./FileLibBot
```