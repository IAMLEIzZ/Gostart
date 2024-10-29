package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// 错误恢复，把错误恢复功能制作成中间件，然后定义在全局
func trace(message string) string {
	var pcs [32]uintptr
	// 通过 runtime.Callers 获取调用栈信息
	n := runtime.Callers(3, pcs[:])		

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		// 通过 runtime.FuncForPC 获取函数信息
		fn := runtime.FuncForPC(pc)
		// 通过 fn.FileLine 获取文件名和行号
		file, line := fn.FileLine(pc)
		// 将文件名和行号写入到 str 中
		str.WriteString(fmt.Sprintf("\n\t%s:%s", file, line))
	}

	return str.String()
}

// 通过 defer + recover 来捕获 panic，然后返回 500 状态码
// 通过 log 输出错误信息
// 通过 c.Fail 返回错误信息
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover();err != nil {
				msg := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(msg))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		// 如果没有错误，则会执行 Next(), 如果发生错误，则会执行 defer 中的 recover()
		c.Next()
	}
}