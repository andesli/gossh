package web

import (
	"github.com/andesli/gossh/auth"
)

type WebDriver struct {
	url string
	key string
	sql string
}

func init() {
	web := &WebDriver{}
	auth.Register("web", web)
}

func (web WebDriver) GetPassword(host, user string) (string, error) {
	//vist http api to get password
	// .....
	return "password", nil
}
