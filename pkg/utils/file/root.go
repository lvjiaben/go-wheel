package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func GetRootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(dir); err == nil {
			return dir, nil
		}

		// 如果已经到达文件系统的根目录，则返回错误
		if dir == filepath.Dir(dir) {
			return "", fmt.Errorf("go.mod not found")
		}

		dir = filepath.Dir(dir)
	}
}

// 生成go文件，参数为：文件路径名，文件内容，原有的是否覆盖（false为不覆盖）
func MakeFile(filePath string, fileName string, content string, cover bool) {
	// 如果不存在就递归创建目录
	if !checkPathIsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	var f *os.File
	var err error
	if checkPathIsExist(fileName) {
		if !cover {
			fmt.Println(fileName + " 文件已存在 -> 当前未覆盖")
			return
		}
		f, err = os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0666) //打开文件
		if err != nil {
			panic(err)
		}
	} else {
		f, err = os.Create(fileName)
		if err != nil {
			panic(err)
		}
	}
	defer f.Close()
	_, err = io.WriteString(f, content)
	if err != nil {
		panic(err)
	}
	fmt.Println("目标： ", fileName, " 已经生成！")
}

// 检查文件是否存在
func checkPathIsExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
