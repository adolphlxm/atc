package atc

import (
	"net/http"
	//"net/http/httptest"
	"testing"
)
import (
	"encoding/json"
	"net/http/httptest"
	//"path"
	"strings"
)

type TestHandler struct {
	Handler
}

func (c *TestHandler) Get() {
	c.Ctx.SetData("code", 200)
	c.JSON()
}

func (c *TestHandler) Post() {
	a1 := c.Ctx.Query("a")
	b1 := c.Ctx.Query("b")
	c1 := c.Ctx.Query("c")
	if a1 != "1" || b1!= "2" || c1 != "3"{
		c.Error(406,1000).Message("Post receive parameters failure.")
		return
	}
	c.Ctx.SetData("code", "200")
	c.Ctx.SetData("a",a1)
	c.Ctx.SetData("b",b1)
	c.Ctx.SetData("c",c1)
	c.JSON()
}

func TestHttpGet(t *testing.T) {
	// A GET request
	r, _ := http.NewRequest("GET", "/V1/user/test", nil)
	w := httptest.NewRecorder()

	// A GET response
	data := make(map[string]int)
	handler, _ := NewHandlerRouter()
	handler.AddRouter("/V1/user", &TestHandler{})
	handler.ServeHTTP(w, r)
	body := w.Body.Bytes()
	json.Unmarshal(body, &data)
	if data["code"] != 200 {
		t.Errorf("url param set to [%v];", data)
	}
}

func TestHttpPost(t *testing.T) {
	// A GET request
	r, _ := http.NewRequest("POST", "/V1/user/test", strings.NewReader("a=1&b=2&c=3"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//r.Header.Set("Content-Type", "text/xml")
	w := httptest.NewRecorder()

	// A GET response
	data := make(map[string]string)
	handler, _ := NewHandlerRouter()
	handler.AddRouter("/V1/user", &TestHandler{})
	handler.ServeHTTP(w, r)
	body := w.Body.Bytes()
	json.Unmarshal(body, &data)
	if data["code"] != "200" {
		t.Errorf("url param set to [%v];", data)
	}
	if data["a"] != "1" || data["b"] != "2" || data["c"] != "3" {
		t.Errorf("url param post receive [%v]", data)
	}

}

//func TestRouteRegexp(t *testing.T) {
//	r, _ := http.NewRequest("GET", "/V1/user/1", nil)
//	w := httptest.NewRecorder()
//
//	handler, _ := NewHandlerRouter()
//	handler.AddRouter("/V1/user/uid:([0-9]+)", UserHandler)
//	handler.ServeHTTP(w, r)
//	body := w.Body.String()
//	//fmt.Println(r)
//	if body != "UserHandler-1" {
//		t.Errorf("url param set to [%s];", body)
//	}
//}
//
//func TestRouteRestful(t *testing.T) {
//	r, _ := http.NewRequest("GET", "/index", nil)
//	w := httptest.NewRecorder()
//	r1, _ := http.NewRequest("POST", "/index", nil)
//	w1 := httptest.NewRecorder()
//	r2, _ := http.NewRequest("PUT", "/index", nil)
//	w2 := httptest.NewRecorder()
//	r3, _ := http.NewRequest("DELETE", "/index", nil)
//	w3 := httptest.NewRecorder()
//
//	//Adds RESTful routing rules
//	handler, _ := NewHandlerRouter()
//	handler.AddRouter("/index", GetHandler).Get()
//	handler.AddRouter("/index", PostHandler).Post()
//	handler.AddRouter("/index", PutHandler).Put()
//	handler.AddRouter("/index", DeleteHandler).Delete()
//
//	//GET
//	handler.ServeHTTP(w, r)
//	body := w.Body.String()
//	if body != "GetHandler" {
//		t.Errorf("url param set to [%s];", body)
//	}
//
//	//POST
//	handler.ServeHTTP(w1, r1)
//	body1 := w1.Body.String()
//	//fmt.Println(r)
//	if body1 != "PostHandler" {
//		t.Errorf("url param set to [%s];", body1)
//	}
//
//	//PUT
//	handler.ServeHTTP(w2, r2)
//	body2 := w2.Body.String()
//	//fmt.Println(r)
//	if body2 != "PutHandler" {
//		t.Errorf("url param set to [%s];", body2)
//	}
//
//	//DELETE
//	handler.ServeHTTP(w3, r3)
//	body3 := w3.Body.String()
//	//fmt.Println(r)
//	if body3 != "DeleteHandler" {
//		t.Errorf("url param set to [%s];", body3)
//	}
//}
