# Gostart
## 一个基于 Go 实现的简单 Web 框架
### version 0.1.0
实现了路由映射表， 用户可以自己自定义注册静态路由功能
main.go 为用户启动服务，默认端口为 9999
### version 0.1.1
对路由表相关操作和上下文进行了封装。这里的上下文指的是 Request 和 Response 以及其中携带的信息吗，如状态码等。
增加了对 Json 和 HTML 等返回类型的支持
### version 0.2.0
利用 Trie 树实现了动态路由，动态路由具有参数匹配和通配两种方式,
router_test 为 Trie 路由树的测试方法
### version 0.2.1
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
### version 0.3.0
支持全局或分组的中间件使用，给某一组路由添加统一的中间件。
整个流程大致是，现在 main 中注册好分组，然后给这个分组指定要添加的中间件方法，handler 方法和中间件方法都会被视为 handler（本质上都是业务处理）。
当一个请求传来服务器的时候，服务器先会判断当前这个请求在哪个分组，然后把这个分组和这个分组的父组所包含的中间件都复制给当前请求的 context.handler 中，
然后把请求路径对应的 handler 方法也加入到 context.handler 中，由 context.next() 统一执行。
### version 0.4.0
支持 HTML 模板渲染， 支持静态资源服务
### version 0.5.0
支持作物恢复，将错误恢复功能制作成中间件的形式，然后成默认启动的全局中间件即可
```
// handler 执行
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)

	// 按顺序执行中间件操作
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}
这块写的很有意思，可以按照以下情况讨论：
	1. 如果所有的 handler 和中间件都必须按照规范写 Next()方法，会是什么情况？代码能否退化？
	2. 如果有的 handler 和中间件不写 Next() 方法，for 循环的作用是什么？
这里复制上评论区的一个网友写的理解，我的理解和他大致一样。

在本章中，对于中间件的串行执行，我们使用了一个Next方法，这里有必要对这个方法进行详细解释：
func (c *Context) Next() {
	c.index++
	s := len(c.middleware)
	for ; c.index < s; c.index++ {
		c.middleware[c.index](c)
	}
}
如上，如果我们假设中间件的写法都是：

func(ctx *gee.Context) {
    //do something
    ctx.Next()
    //do something
}
也即每个中间件都调用ctx.Next()，那我们其实可以将上述Next改为

func (c *Context) Next() {
	c.index++
	if(c.index < len(c.middleware)){
		c.middleware[c.index](c)
	}
}
如上不需要循环，因为每个middleWare里都调用了ctx.Next()，这个链式调用会一直走下去
当走到最后一个节点，又会由于函数调用入栈的原因，调用结束后会出栈反向执行

但很明显，我们得考虑中间件中无ctx.Next()的情况，如果中间件中无ctx.Next()，上述版本就会出现调用断掉
因此我们得加for，手动遍历循环这些middleware，避免middleware中断掉的情况

这里存在一个容易疑惑的点是：
假设T1 开始执行Next()，index变为0，然后进入for，for内看着似乎会遍历所有的middleware，假设执行到T2 某个middleware内也包含Next()调用，此时再次进入Next()，又开启了一次for循环，那会不会导致部分middleware重复执行呢？

我们做个假设，现在有4个middleware A、B、C、D,那加上原本的路由处理函数h，我们就得到了5个handler。
从最极端的假设开始，A、B、C、D内部都调用了Next()，它们的结构如下：

func(ctx *gee.Context) {
    before()
    ctx.Next()
	after()
}
首先路由请求，进入Next()，此时index++变为0，然后开始for循环第一次迭代，for循环先执行middleware[0]也即A，A内调用Next()，index++更新为1，然后开启第二个for循环第一次迭代，由于此时index == 1，因此执行B，重复如此，直到D，D内调用Next()的时候，index被更新为4，因此本次for执行的是middleware[4]也即路由处理函数h，

//index为4 s为5
for ; c.index < s; c.index++ {
	//c.middleware[c.index]是h
	c.middleware[c.index](c)
}
h执行完成后，for循环迭代c.index++，此时index被更新为5，不满足for循环退出，然后回到它的调用点，也即回到D的Next()执行后，D执行after()，执行完后，for循环迭代c.index++，此时index被更新为6，同样不满足迭代返回到调用点，此时到C的after()，重复迭代一直到A的after()，最后整个middleware执行完毕。

所以上面这种情况，虽然会多次进入for循环，但每个for其实只会迭代一次，不会重复执行middleware

再考虑另一个情况，假设B和C内都无Next()调用，我们再分析会发生什么：

首先路由请求，进入Next()，此时index++变为0，然后开始for循环第一次迭代，for循环先执行middleware[0]也即A，A内调用Next()，index++更新为1，然后开启第二个for循环第一次迭代，由于此时index == 1，因此执行B，B内无Next()因此B执行完成后控制权回到for，for循环迭代c.index++，此时index被更新为2，然后满足for继续迭代，此时执行C，同样C内无Next()，执行完后，for循环迭代c.index++，index被更新为3，然后满足for继续迭代，开始执行D，
D内包含Next()调用Next()的时候index被更新为4，执行f，f执行完成，index被更新为5，不满足继续迭代，控制权回到D的after(),D的after()执行完成，index被更新为6，
不满足继续迭代，控制权此时回到A，A执行after(),index被更新为6，不满足继续迭代，整个流程执行完成。

可以看到，我们总结下就是，对于中间件中不包含Next()调用的，是由for循环的迭代c.index++来实现调用下一个middleware的，而中间件如果包含Next()，则是通过进入Next一开始的那行c.index++实现调用下一middleware的
```