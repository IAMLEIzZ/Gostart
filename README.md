# Gostart
## 一个基于 Go 实现的简单 Web 框架
### version 0.1.0
实现了路由映射表， 用户可以自己自定义注册静态路由功能
main.go 为用户启动服务，默认端口为 9999
### version 0.1.1
对路由表相关操作和上下文进行了封装。这里的上下文指的是 Request 和 Response 以及其中携带的信息吗，如状态码等。
增加了对 Json 和 HTML 等返回类型的支持
### version 0.2.1
利用 Trie 树实现了动态路由，动态路由具有参数匹配和通配两种方式,
router_test 为 Trie 路由树的测试方法
### version 0.2.2
支持分组控制路由
```
	关于路由组的理解，首先在 Engine 中嵌套了 RouterGroup，这样 Engine 也可以作为一个路由组
	一般来说，用 engine 来创建一个总的路由组，再用 RouterGroup，来创建单个的路由分组
	例如： /index/v1/hello/doc 和 /index/v2/hello/doc，先用 engine 将 index 作为一个大的分组
	再在其中建立两个叫 v1 和 v2 的小的分组，而 v1 和 v2 没有 groups 属性，可以认为是一个较弱的分组对象，
	而只有 Engine 有 groups 属性，是一较强的分组对象，他包含了他之下所有的分组信息。
	• Engine 作为根路由组，拥有 groups 字段来存储所有的路由组。通过这个字段，Engine 知道了自己管理的所有子路由组的信息。
	  其他的 RouterGroup 不拥有 groups 字段，它们不需要存储全局的路由组信息。
	• RouterGroup 作为子路由组，并不管理其他路由组，它们只是用来为某些路由定义特定的前缀或局部中间件。
	  每个子 RouterGroup 都通过 parent 字段与其上级 RouterGroup 或 Engine 关联，形成一个层次结构。
```