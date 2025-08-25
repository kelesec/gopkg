# gopkg

# About
Golang 常用库封装，便于快速导入使用

# 目录
## [httpx](doc/httpx.md)
基于 fasthttp 二次封装的请求库，支持链式调用

TODO:
- [x] 自定义重定向支持：禁止重定向/允许重定向N次
- [x] 重定向响应记录：支持获取每次重定向请求的相应内容

## [fileutils](doc/fileutils.md)
对常见文件读写能力进行封装，方便快捷使用

TODO:
- [x] xlsx读写支持
- [ ] 压缩文件支持

## [iputils](doc/iputils.md)
ip/url资产整理：

TODO:
- [x] 资产去重+排序整理
- [x] IP地址生成，如 x.x.x.1/16 生成对应的IP地址数组
- [x] 从乱序的IP地址中提取出C段信息

# Changelog
## 2025-08
- update: 增加 httpx 模块，基于 fasthttp 封装链式调用 http 请求库
- update: httpx 增加代理支持，可配置 HTTP/HTTPS/SOCKS5 等代理类型
- update: httpx 增加重定向能力支持，默认不启用重定向，并且支持记录每次重定向请求的响应
- update: 增加fileutils模块
  - 支持基本文件读写能力，可通过 chan 进行文件实时读写；
  - 增加 xlsx 文件读写能力支持
- update: 增加 iputils 模块
  - 支持CIDR/IP范围生成对应的IP地址
  - 支持将IP数组转为CIDR格式
  - 支持IP地址整理：去重、格式校验、排序、提取公网/内网IP地址等
