<!--
 * @Author: xiaozuhui
 * @Date: 2022-12-02 12:22:19
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2023-01-05 22:47:34
 * @Description: 
-->

# chat-shock

聊天工具后端代码
尽可能支持私密、加密的聊天
目前还在项目初期

# 后端需求

1. 用户模块
2. 聊天室模块

## 主要功能

### 文件模块

- [X]  主动接入minio
- [X]  包装、隐藏 url（以及刷新url）和文件名、bucket之间的关系
- [ ]  提供统一的接口和实现，以应对不同文件服务器的功能(oss?)
- [X]  减少redis之类中间件的侵入性
- [ ]  对文件名进行md5加密
- [ ]  对文件进行hash校验
- [ ]  上传下载分片实现异步
- [ ]  上传下载断点续传
- [ ]  重复文件筛查

### 用户模块

- [X]  ~~手机号注册~~
- [X]  接入minio，用于头像
- [X]  JWT
- [X]  删除用户
- [X]  用户好友

### NOTICE模块

- [ ]  推送消息
- [ ]  websocket接入

### 聊天室

- [X]  websocket
- [X]  接入更加广泛的minio，用于发送的图片、视频、音频
- [ ]  缓存未发出的消息
- [ ]  增加根据名称查询聊天室（已经为unque）
- [ ]  名称的模糊查询
- [X]  获取聊天室用户列表

### 私信

### 群组聊天

### 社区话题

### 朋友圈

### 消息模块

- [ ]  消息模型
- [ ]  已读未读
- [ ]  websocket

## 后续正式功能

- [X]  **邮箱接入** (该功能比较重要)
- [ ]  ~~实名认证~~
- [X]  删除手机注册、已经手机号相关功能
- [X]  健壮日志接入
- [ ]  日志管理，能将日志发送到日志管理平台
- [ ]  接入CI/CD
- [ ]  消息系统
- [X]  返回标准化
- [ ]  错误链路追踪
- [ ]  Swagger
- [ ]  加密通讯
- [ ]  **推送**
- [ ]  单点登录
- [ ]  websocket增加token验证
- [ ]  单元测试
- [ ]  **评断是否需要好友功能**
- [ ]  **自部署邮件服务**（尝试使用Apache James？）

## 未来上线功能

- [ ]  视频聊天
- [ ]  语音聊天
