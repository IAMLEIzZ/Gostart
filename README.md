# Gostart
## 一个基于 Go 实现的简单 Web 框架
### version 0.1.0
实现了路由映射表， 用户可以自己自定义注册静态路由功能
main.go 为用户启动服务，默认端口为 9999
### version 0.1.1
对路由表相关操作和上下文进行了封装。这里的上下文指的是 Request 和 Response 以及其中携带的信息吗，如状态码等。
增加了对 Json 和 HTML 等返回类型的支持
### version 0.2.1
利用 Trie 树实现了动态路由，动态路由具有参数匹配和通配两种方式
router_test 为 Trie 路由树的测试方法