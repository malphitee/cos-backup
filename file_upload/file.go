package file_upload

import (
	config2 "backFolderToCos/config"
	"backFolderToCos/cos_tool"
	"backFolderToCos/notification"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type FileTool struct {
}

func (f *FileTool) ListDir(dirPath string) []string {
	var allFiles []string
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			list := f.ListDir(dirPath + "/" + file.Name())
			for _, ff := range list {
				allFiles = append(allFiles, ff)
			}
		} else {
			allFiles = append(allFiles, dirPath+"/"+file.Name())
		}
	}
	return allFiles
}

func (f *FileTool) CalcFilesMd5(path string, files []string) map[string]string {
	m := sync.Map{}
	wg := sync.WaitGroup{}
	for _, file := range files {
		wg.Add(1)
		go func(path string, file string, wg *sync.WaitGroup, m *sync.Map) {
			pFile, err := os.Open(file)
			defer pFile.Close()
			if err != nil {
				panic(fmt.Sprintf("计算文件md5值失败, err = %v", err))
			}
			md5h := md5.New()
			_, err = io.Copy(md5h, pFile)
			if err != nil {
				panic(fmt.Sprintf("计算文件md5值失败, err = %v", err))
			}
			if strings.HasPrefix(file, path) {
				newPath := strings.ReplaceAll(file, path+"/", "")
				//fmt.Println(newPath + " -> " + hex.EncodeToString(md5h.Sum(nil)))
				m.Store(newPath+"", hex.EncodeToString(md5h.Sum(nil)))
			} else {
				fmt.Printf("文件夹路径异常，file = %s, path = %s \n", file, path)
			}
			wg.Done()
		}(path, file, &wg, &m)
	}
	wg.Wait()
	res := make(map[string]string)
	m.Range(func(key, value interface{}) bool {
		res[key.(string)] = value.(string)
		return true
	})
	return res
}

func (f *FileTool) SyncToCos() {
	// 获取配置
	config := config2.GetConfigFromYaml()
	// 获取文件列表
	files := f.ListDir(config.Dir + "/" + config.Path)
	// 获取文件md5
	m := f.CalcFilesMd5(config.Dir, files)
	// 获取文件配置
	cosTool := cos_tool.CosTool{
		Prefix:    "",
		Delimiter: "",
		Config:    config,
	}
	bucketResultList := cosTool.GetBucketFileList()
	// 遍历对比
	// 可能存在本地新增文件COS上没有，把COS上文件的 名字->md5 整理成map，进行处理
	contentMap := make(map[string]string)
	for _, content := range bucketResultList.Contents {
		fileName := content.Key
		md5Val := strings.ReplaceAll(content.ETag, "\"", "")
		contentMap[fileName] = md5Val
	}
	wg := sync.WaitGroup{}
	needUploadFiles := make(map[string][]string)
	lock := sync.RWMutex{}
	for localFileName, localMd5 := range m {
		if _, ok := contentMap[localFileName]; ok && contentMap[localFileName] == localMd5 {
			fmt.Println("文件 ", localFileName, " 未发生变动，跳过上传流程")
		} else {
			// 上传的场景：md5变动或者文件没找到
			go func(wg *sync.WaitGroup, dir string, fileName string, needUploadFiles map[string][]string, lock *sync.RWMutex) {
				wg.Add(1)
				d := dir
				n := fileName
				f.Upload(d, n, wg, needUploadFiles, lock)
			}(&wg, config.Dir, localFileName, needUploadFiles, &lock)
		}
	}
	wg.Wait()

	// 如果需要删除本地已经不存在的文件，反向比较一下
	if config.DeleteRemote {
		// @todo
	}

	// 如果有信息需要通知，进行ServerChan推送
	if len(needUploadFiles) > 0 {
		(&notification.NotifyTool{}).DoServerChanNotify(config, needUploadFiles)
	}

}

func (f *FileTool) Upload(path string, fileName string, wg *sync.WaitGroup, needUploadFiles map[string][]string, lock *sync.RWMutex) {
	fmt.Println("将要上传的文件路径 = ", path+"/"+fileName, " 文件名 = ", fileName)
	err := (&cos_tool.CosTool{}).UploadToCos(path+"/"+fileName, fileName)
	wg.Done()
	lock.Lock()
	if err != nil {
		needUploadFiles["failure"] = append(needUploadFiles["failure"], fileName)
	} else {
		needUploadFiles["success"] = append(needUploadFiles["success"], fileName)
	}
	lock.Unlock()
}
