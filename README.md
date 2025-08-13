# gopkg

# About
Golang 常用库封装，便于快速导入使用

# 目录
## httpx
基于 fasthttp 二次封装的请求库，支持链式调用

TODO:
- [x] 自定义重定向支持：禁止重定向/允许重定向N次
- [x] 重定向响应记录：支持获取每次重定向请求的相应内容

# Changelog
## 2025-08
- update: 增加 httpx 模块，基于 fasthttp 封装链式调用 http 请求库
- update: httpx 增加代理支持，可配置 HTTP/HTTPS/SOCKS5 等代理类型
- update: httpx 增加重定向能力支持，默认不启用重定向，并且支持记录每次重定向请求的响应