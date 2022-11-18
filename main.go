package main

import "backFolderToCos/file_upload"

func main() {
	fileTool := &file_upload.FileTool{}
	fileTool.SyncToCos()
}
