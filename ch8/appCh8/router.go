package main

import (
	"net/http"
	"strings"
)

type router struct {
	// http method - url pattern - HandlerFunc
	handlers map[string]map[string]http.HandlerFunc
}

// struct router 에 멤버 함수 추가
func (r *router) HandleFunc(method, pattern string, h http.HandlerFunc) {
	// http 메서드로 등록된 맵이 있는지 확인
	m, ok := r.handlers[method] // methodName:string 이 key ?
	// method 키로 값을 가져온게 m ? url pattern - HandlerFunc
	// 정상이 아니면 뭘 가져오는거야 m ?
	if !ok { // 정상이 아니면: http method가 없으면 ?
		// 등록된 맵이 없으면 새 맵을 생성
		m = make(map[string]http.HandlerFunc) // create slice
		r.handlers[method] = m
	}
	// method로 등록된 맵에 URL patten과 핸들러 함수 등록
	m[pattern] = h
}

func (r *router) handler() HandlerFunc { // 반환형, 함수 이름만 적어도 되는건가 ?
	return func(c *Context) {
		// HTTP Method에 맞는 모든 handlers 를 반복하면서
		// 요청 url에 해당하는 handler를 찾음
		for pattern, handler := range r.handlers[c.Request.Method] {
			// 요청 http method 에 따라, 해당 url pattern 루프
			if ok, params := match(pattern, c.Request.URL.Path); ok {
				// 일치하면 true, url param 정보 준다
				for k, v := range params {
					c.Params[k] = v
				}
				// 요청 url 에 해당하는 handler 수행
				handler(c) // ?
				return
			}
		}

		// 요청 url에 해당하는 handler 를 찾지 못한 경우
		// Not Found 에러 처리
		http.NotFound(c.ResponseWriter, c.Request)
		return
	}
}

func match(pattern, path string) (bool, map[string]string) {
	// pattern 과 path 가 정확히 일치하는 경우 즉시 true 반환
	if pattern == path {
		return true, nil
	}

	// 패턴과 패쓰를 / 단위로 구분
	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	// 패턴과 패쓰를 /로 구분한 후,
	// 부분 문자열 집합의 갯수가 다르면 false를 리턴
	if len(patterns) != len(paths) {
		return false, nil
	}

	// 패턴에 일치하는 url 파라미터를 담기 위한 params 맵 생성
	params := make(map[string]string)

	// /로 구분된 패턴/패쓰 각 문자열을 하나씩 비교
	for i := 0; i < len(patterns); i++ {
		switch {
		case patterns[i] == paths[i]:
			// 패턴과 패쓰의 부분 문자열이 일치하는 경우, 바로 다음 루프 수행
		case len(patterns[i]) > 0 && patterns[i][0] == ':':
			// 패턴이 : 문자로 시작하는 경우
			// params 에 url param 을 담은 후, 다음 루프 수행
			params[patterns[i][1:]] = paths[i] // 변수 이름 - 값
		default:
			// 일치하는 경우가 없으면 false 리턴
			return false, nil
		}
	}

	return true, params
}
