package notification

import (
	"cos_backup/config"
	"fmt"
	"io"
	"net/http"
	urllib "net/url"
)

type Gotify struct {
}

func (g *Gotify) SendNotify(title string, message string, short string) error {
	config := config.GetConfig()
	token := config.GotifyToken
	if len(token) == 0 {
		panic("gotify token 未设置")
	}
	url := config.GotifyUrl
	if len(token) == 0 {
		panic("gotify url 未设置")
	}
	url += "/message?token=" + token
	rsp, err := http.PostForm(url, urllib.Values{"title": {title}, "message": {message}})
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(rsp.Body)
	if err != nil {
		fmt.Println("SendGotifyMsg err ---- ", err)
		return err
	}
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println("SendGotifyMsg err ---- ", err)
	}
	fmt.Println("SendGotifyMsg -- ", string(body))
	return err
}
