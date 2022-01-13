package file_system

import (
	"fmt"
	"os"
)

type QueueFilePath string

func (queuePath QueueFilePath) Open() QueueFile {
	fileForWrite, err := os.OpenFile(string(queuePath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic(err.Error())
	}
	return QueueFile{fileForWrite}
}

type QueueFile struct {
	File *os.File
}

func (currentFile QueueFile) WritePath(path string) {
	pathLength := len(path)
	if pathLength < 30 {
		return
	}

	recordString := fmt.Sprintf("%v\n", path)
	_, err := currentFile.File.WriteString(recordString)
	if err != nil {
		panic(err.Error())
	}
}