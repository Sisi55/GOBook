package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"path/filepath"
	"text/template"
)

type Context struct {
	Params map[string]interface{}

	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type HandlerFunc func(*Context)

// struct Context 멤버 함수
func (c *Context) RenderJson(v interface{}) {
	// HTTP Status를 StatusOK 로 설정 -- json render 할 때 ?
	c.ResponseWriter.WriteHeader(http.StatusOK)
	// 아 WriterHeader --> Writer

	// Content-Type 를 application/json 으로 지정
	c.ResponseWriter.Header().Set(
		"Cntent-Type", "application/json; charset=utf-8")

	// v 값을 json으로 출력
	if err := json.NewEncoder(c.ResponseWriter).Encode(v); err != nil {
		// NewEncoder 가 rw 이용해서 Encoder 객체를 생성하고
		// v 객체를 marshal 한다. -- 직렬화

		// 아무튼 정상이 아니면
		c.RenderErr(http.StatusInternalServerError, err)
	}
}

func (c *Context) RenderXml(v interface{}) { // render json 과 같은 로직
	c.ResponseWriter.WriteHeader(http.StatusOK)
	c.ResponseWriter.Header().Set("Context-Type", "application/xml; charset=utf-8")

	if err := xml.NewEncoder(c.ResponseWriter).Encode(v); err != nil {
		c.RenderErr(http.StatusInternalServerError, err)
	}
}

func (c *Context) RenderErr(code int, err error) {
	if err != nil { // 밖에서도 검사했는데 또 한다 !
		if code > 0 {
			// 정상적인 code를 전달한 경우, HTTP Status를 해당 code로 지정한다
			http.Error(c.ResponseWriter, http.StatusText(code), code)
		} else {
			// 정상적인 code가 아닌 경우
			// HTTP Status를 StatusInternalServerError 로 지정한다.
			defaultErr := http.StatusInternalServerError
			http.Error(c.ResponseWriter, http.StatusText(defaultErr), defaultErr)
			// defaultErr 매개3은 Header에 쓰고
			// string 은 화면에 출력한다 fmt.Print
		}
	}
}

func (c *Context) Redirect(url string) {
	http.Redirect(c.ResponseWriter, c.Request, url, http.StatusMovedPermanently)
}

// templates: 템플릿 객체를 보관하기 위한 map
var templates = map[string]*template.Template{}

func (c *Context) RenderTemplate(path string, v interface{}) {
	// path 에 해당하는 템플릿이 존재하는지 확인
	t, ok := templates[path]
	if !ok {
		// 정상이 아니면 -- 이 표현은 의미론적이지 않은것 같다!
		// path에 템플릿이 존재하지 않으면, 템플릿 객체 생성
		t = template.Must(template.ParseFiles(
			filepath.Join(".", path))) // ??
		templates[path] = t
	}

	// v 값을 템플릿 내부로 전달하여 만들어진 최종 결과를
	// c.ResponseWriter 에 출력
	t.Execute(c.ResponseWriter, v)
}
