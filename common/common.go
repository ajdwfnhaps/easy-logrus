package common

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	uuid "github.com/satori/go.uuid"
)

var (
	//Pid 应用进程id
	Pid = os.Getpid()
)

// NewTraceID 创建追踪ID
func NewTraceID() string {
	return fmt.Sprintf("pid:%d-%s",
		Pid,
		UniqueID())
}

//NewUUID 创建通用唯一识别码
func NewUUID() string {
	ret := uuid.NewV4()
	return ret.String()
}

//PathExists 判断一个文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//GetMd5String 生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//UniqueID 生成Guid字串
func UniqueID() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}
