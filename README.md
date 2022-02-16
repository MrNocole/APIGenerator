## APIGenerator

主要技术：

Golang--Web主体

HTML

Gin--Web框架

MySQL

JSON



核心功能：

允许开发者/用户在本地/服务器快速创建一个Web应用，管理各种资源，生成API接口，便于本地模拟线上环境/线上环境资源管理。

支持登录注册，可集中化管理多来源的资源。

资源集中化存储，用md5标识，多用户有相同文件仓库只存一份。



开发计划：

优化文件仓库的存储逻辑。

在线修改，实现copy on write

支持更多种文件的预览/API生成 （目前仅支持预览JSON）

做一个选择本地（无数据库使用的接口）
