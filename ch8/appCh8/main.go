package main

import (
	"crypto/hmac"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type User struct {
	Id        string
	AddressId string
}

const VerifyMessage = "verified"

func main() {
	s := NewServer()

	s.HandleFunc("GET", "/", func(c *Context) {
		c.RenderTemplate("/public/index.html",
			map[string]interface{}{"time": time.Now()})
	})

	s.HandleFunc("GET", "/about", func(c *Context) {
		fmt.Fprintln(c.ResponseWriter, "about")
	})

	s.HandleFunc("GET", "/user/:id", func(c *Context) {
		u := User{Id: c.Params["id"].(string)}
		c.RenderXml(u) // 왜 xml 으로 렌더링할까?
	})

	s.HandleFunc("POST", "/users", func(c *Context) {
		// Context에 결과가 담기나 보다
		c.RenderJson(c.Params)
	})

	s.HandleFunc("GET", "/login", func(c *Context) {
		c.RenderTemplate("/public/login.html",
			map[string]interface{}{"message": "로그인이 필요합니다."})
	})

	s.HandleFunc("POST", "/login", func(c *Context) {
		// 로그인 정보를 확인하여 쿠키에 인증 토큰 값 기록
		if CheckLogin(c.Params["username"].(string), c.Params["password"].(string)) {
			http.SetCookie(c.ResponseWriter, &http.Cookie{
				Name:  "X_AUTH",
				Value: Sign(VerifyMessage),
				Path:  "/",
			})
			c.Redirect("/")
		}
		// id, password 맞지 않으면 다시 /login 페이지 렌더링
		c.RenderTemplate("/public/login.html",
			map[string]interface{}{"message": "id 또는 password가 일치하지 않습니다."})
	})

	s.Use(AuthHandler)
	s.Run(":8080")
}

func CheckLogin(username, password string) bool {
	// 로그인 처리
	const (
		USERNAME = "tester"
		PASSWORD = "12345"
	)
	return username == USERNAME && password == PASSWORD
}

func AuthHandler(next HandlerFunc) HandlerFunc {
	ignore := []string{"/login", "public/index.html"}
	return func(c *Context) {
		// url prefix가 /, /login, public/index.html 인 경우
		// auth를 체크하지 않음
		for _, s := range ignore {
			if strings.HasPrefix(c.Request.URL.Path, s) {
				next(c) // 마지막 진짜 Handler 과정 전에
				// AuthHandler 가 미들웨어 역할 하는건가 ?
				// 아 사이 handler 만드는게 미들웨어인가 ?
				return
			}
		}

		if v, err := c.Request.Cookie("X_AUTH"); err == http.ErrNoCookie {
			// X_AUTH 쿠키 값이 없으면 /login 으로 이동
			c.Redirect("/login")
			return
		} else if err != nil {
			// 정상이 아니면
			c.RenderErr(http.StatusInternalServerError, err)
			return
		} else if Verify(VerifyMessage, v.Value) {
			// 쿠키 값으로 인증이 확인되면 다음 핸들러로 넘어감
			next(c)
			return
		}

		// 위 해당사항 없으면
		// /login으로 이동
		c.Redirect("/login")
	}
}

// 인증 토큰 생성
func Sign(message string) string {
	secretKey := []byte("golang-book-secret-key2")
	if len(secretKey) == 0 {
		return ""
	}
	mac := hmac.New(shal.New, secretKey)
	io.WriteString(mac, message)
	return hex.EncodeToString(mac.Sum(nil))
}

// 인증 토큰 확인
func Verify(message, sig string) bool {
	return hmac.Equal([]byte(sig), []byte(Sign(message)))
}
