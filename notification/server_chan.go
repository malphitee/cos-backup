package notification

import (
	config2 "backFolderToCos/config"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
)

type ServerChanNotify struct {
}

func (n *ServerChanNotify) SendServerChanMsg(title string, description string, short string) error {
	config := config2.GetConfig()
	domain := "https://sctapi.ftqq.com/"
	params := url2.Values{}
	params.Set("title", title)
	params.Set("desp", description)
	params.Set("short", short)
	url := domain + config.ServerChanSendKey + ".send"
	rsp, err := http.PostForm(url, params)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(rsp.Body)
	if err != nil {
		fmt.Println("SendServerChanMsg err ---- ", err)
		return err
	}
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println("SendServerChanMsg err ---- ", err)
	}
	fmt.Println("SendServerChanMsg -- ", string(body))
	return err
}

func (n *ServerChanNotify) SendNotify(title string, desc string, short string) error {
	err := n.SendServerChanMsg(title, desc, short)
	if err != nil {
		fmt.Println("发送ServerChan通知失败，err = ", err)
	}
	return err
}
