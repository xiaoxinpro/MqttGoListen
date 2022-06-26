package main

import (
	"bytes"
	"math/rand"
	"os"
	"time"
)

// 检查文件或文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil //文件存在
	}
	if os.IsNotExist(err) {
		return false, nil //文件不存在
	}
	return false, err //不确定是否存在
}

//随机字符生成
func RandChar(size int) string {
	const char = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.NewSource(time.Now().UnixNano()) // 产生随机种子
	var s bytes.Buffer
	for i := 0; i < size; i++ {
		s.WriteByte(char[rand.Int63()%int64(len(char))])
	}
	return s.String()
}
