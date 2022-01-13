package file_system

import (
	"io/ioutil"
	"os"
)

type ParserDirectoryList map[string]os.FileInfo

func GetCurrentDir() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}

	return path
}

func GetParserXmlDirectory(checkDirectory string) ParserDirectoryList {
	files, err := ioutil.ReadDir(checkDirectory)
	if err != nil {
		panic(err)
	}

	directoryMaps := make(ParserDirectoryList)
	for _, file := range files {
		if file.IsDir() {
			directoryMaps[file.Name()] = file
		}
	}

	return directoryMaps
}

