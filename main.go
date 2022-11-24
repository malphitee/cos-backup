package main

import (
	"backFolderToCos/file_upload"
	"fmt"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("发生错误, err = ", err)
		}
	}()
	fileTool := &file_upload.FileTool{}
	fileTool.SyncToCos()
}
