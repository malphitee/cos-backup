package notification

import (
	config2 "backFolderToCos/config"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
)

type NotifyTool struct {
	ServerChanSendKey string
}

func (n *NotifyTool) SendServerChanMsg(config config2.Config, title string, description string, short string) error {
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

func (n *NotifyTool) DoServerChanNotify(config config2.Config, needUploadFiles map[string][]string) {
	title := "BitWarden文件变动上传至COS通知"
	desc := ""
	if len(needUploadFiles["success"]) > 0 {
		desc = desc + fmt.Sprintf("\r\n### 共成功上传 %d 个文件\n\n成功文件列表\n\n", len(needUploadFiles["success"]))
		for _, item := range needUploadFiles["success"] {
			desc = desc + "+ " + item + "\r\n"
		}
		desc = desc + "---\r\n"
	}
	if len(needUploadFiles["failure"]) > 0 {
		desc = desc + fmt.Sprintf("\r\n### 共上传失败 %d 个文件\n\n失败文件列表\n\n", len(needUploadFiles["failure"]))
		for _, item := range needUploadFiles["failure"] {
			desc = desc + "+ " + item + "\r\n"
		}
		desc = desc + "---\r\n"
	}
	short := fmt.Sprintf("共上传成功 %d 个\n\n上传失败%d个\n\n", len(needUploadFiles["success"]), len(needUploadFiles["failure"]))
	err := n.SendServerChanMsg(config, title, desc, short)
	if err != nil {
		fmt.Println("发送ServerChan通知失败，err = ", err)
	}
}
