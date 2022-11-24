package notification

import (
	"backFolderToCos/config"
	"fmt"
)

type NotifyInterface interface {
	SendNotify(string, string, string) error
}

func SendNotify(title string, message string, short string) {
	config := config.GetConfig()
	driver := config.NotifyDriver
	var notifyTool NotifyInterface
	switch driver {
	case "", "server_chan":
		notifyTool = &ServerChanNotify{}
	case "gotify":
		notifyTool = &Gotify{}
	default:
		panic("notify driver error")
	}
	err := notifyTool.SendNotify(title, message, short)
	if err != nil {
		fmt.Println("发送ServerChan通知失败,err = ", err)
	}
}

func GetNotifyData(needUploadFiles map[string][]string) (string, string, string) {
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
	return title, desc, short
}
