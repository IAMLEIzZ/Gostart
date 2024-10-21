package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// Context 
type Context struct {
	Writer http.ResponseWriter
	Req *http. Request

	Path string
	Method string
	StatusCode int
}

// newContext 
// @param w 
// @param req 
// @return *Context  
// @author IAMLEIzZ   
// @date 2024-10-21 03:15:09 
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	context := &Context{
		Writer: w,
		Req: req,
		Path: req.URL.Path,
		Method: req.Method,
	}

	return context
}

// PostForm 
// @receiver c 
// @param key 
// @return string  
// @author IAMLEIzZ   
// @date 2024-10-21 03:15:12 
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 
// @receiver c 
// @param key 
// @return string  
// @author IAMLEIzZ   
// @date 2024-10-21 03:15:14 
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 
// @receiver c 
// @param code 
// @author IAMLEIzZ   
// @date 2024-10-21 03:15:16 
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 
// @receiver c 
// @param key 
// @param value 
// @author IAMLEIzZ   
// @date 2024-10-21 03:15:18 
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String 
// @receiver c 
// @param code 
// @param format 
// @param values 
// @author IAMLEIzZ   
// @date 2024-10-21 03:15:20 
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 
// @receiver c 
// @param code 
// @param obj 
// @author IAMLEIzZ   
// @date 2024-10-21 03:15:22 
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 
// @receiver c 
// @param code 
// @param data 
// @author IAMLEIzZ   
// @date 2024-10-21 03:15:24 
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML 
// @receiver c 
// @param code 
// @param html 
// @author IAMLEIzZ   
// @date 2024-10-21 03:15:26 
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(html))
}