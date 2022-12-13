package main

import (
	"fmt"
	"time"

	"github.com/karmadon/remauth"
)

const authHeader = "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6ImY3NjJjYWE3M2ExOGYwNTAyYjVmMWMxMWNkNDllNTUy" +
	"ODNhZTY1ZmUwYTZkMDY3OTgwYThlMzg1OWM5ZmVhMDhmZjViY2U2NDBlOWMxY2E0In0.eyJhdWQiOiIxIiwianRpIjoiZjc2MmNhYTczYTE4Zj" +
	"A1MDJiNWYxYzExY2Q0OWU1NTI4M2FlNjVmZTBhNmQwNjc5ODBhOGUzODU5YzlmZWEwOGZmNWJjZTY0MGU5YzFjYTQiLCJpYXQiOjE2NDA2MjA2M" +
	"zcsIm5iZiI6MTY0MDYyMDYzNywiZXhwIjoxNjcyMTU2NjM3LCJzdWIiOiI1Iiwic2NvcGVzIjpbXX0.m7jPDmXN0gI_lP0xjJE3anoVKVtZJJvzg" +
	"G0Jl4uWPI2EJllJMNfHBTJV8_5-IVKXvnwjWxIY6AtnyEaLXoS1sHjO-Op0Rh-IXsonfVEeHK8P0pG_8LhEQ-34IoJMnL5fsSQYOCy05WCKfAxATT" +
	"qOfO1-OW3DcNvfx5N3iTpL1jKtxw3eU462kVEWf-n0My1u-6kUqeiKT9uf-iYQV07u9G9QpBNg4otUFKk_-CDCK1z1VwNoYTUGidl0VzIFRqTwKFnc" +
	"voSYMiyTuOv57oxEfpmF_HwD1uv0ckWlKv-e6MADMMmAqR7BqRmcU5ccU8mZyeoA3k0VUPYVATb7UfK03olkasRHmE7cljtJUDnOzb20DpptbeIYkx" +
	"TB2FiuiNPE_ITL86KvuVtqA6HWf73ZYa-mlF3yfZK5tMaMIUMUXKUm0JycBeRB0GE2JevKfohz1-39GvlRBSGj2SOY9QljK41ezwgY8hj60KFNIbHFm" +
	"pPiYctRldvnrIC-4uMxqAQ9iqCZ2zNgZH6Pd6pzksmD8JtsYgvCVsaKBjDJXDJ1JICwm5-yku_oleoSlfURcslxwzRnkg22RGgs8sbhIl6cB_yyR7o" +
	"yvk2T5Q83aTQuX-XTiJ_08Ht9rzjvgynBQOaSJOtBCmpyeZmmGIoW-uV4kYlT8G7ZuFhXiR7pamQkVdY10o"

func main() {
	options := &remauth.Options{
		Debug:     true,
		CheckUrl:  "https://example.com/client/operator/whoami",
		Timeout:   10,
		CacheTime: 5 * time.Minute,
	}

	remoteAuth, _ := remauth.New(options)

	if remoteAuth.Valid(authHeader) {
		println("ok")
	} else {
		println("nok")
	}

	for i := 0; i < 10; i++ {
		remoteAuth.Valid(authHeader)
	}

	time.Sleep(1 * time.Second)
	for i := 0; i < 10; i++ {
		remoteAuth.Check(authHeader, func(response interface{}, err error) {
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("%#v", response)
		})
	}
}
