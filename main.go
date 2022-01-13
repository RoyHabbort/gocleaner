package main

import (
	"./config_reader"
	"./database"
	"./file_system"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"sync"
	"time"
)

type DeletingMapType map[string]string

func (deletingMap DeletingMapType) ShowDeleteQueueInConsole() bool {
	for _, path := range deletingMap {
		fmt.Println(path)
	}
	fmt.Printf("Total count paths: %v\n", len(deletingMap))
	return true
}

func (deletingMap DeletingMapType) WriteQueueIntoFile() bool {
	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")

	filePathString := fmt.Sprintf("delete_queue-%v.txt", date)
	fmt.Println(filePathString)

	fileForWritePath := file_system.QueueFilePath(filePathString)
	fileForWrite := fileForWritePath.Open()

	defer fileForWrite.File.Close()
	for _, path := range deletingMap {
		fileForWrite.WritePath(path)
	}
	return true
}

func (deletingMap DeletingMapType) Append(checkDir string, parserId string, path string) {
	key := fmt.Sprintf("%v_%v", checkDir, parserId)
	deletingMap[key] = path
}

func main()  {
	var wg sync.WaitGroup

	configIniFilePath := config_reader.IniFileString("config/config.ini")
	configIni := configIniFilePath.Load()

	adapterDatabase := database.InitDatabase(configIni.DatabaseConfig())
	// defer the close till after the main function has finished
	// executing
	defer adapterDatabase.DB.Close()

	//записи из xml - только активные
	xmlMaps := adapterDatabase.GetParserXmlActiveRecords()

	deletingMap := make(DeletingMapType)
	wg.Add(1)
	go CheckParsingDirectory("1", &adapterDatabase, xmlMaps, &deletingMap, &wg)
	wg.Add(1)
	go CheckParsingDirectory("2", &adapterDatabase, xmlMaps, &deletingMap, &wg)

	wg.Wait()

	var agree string
	fmt.Print("Чтобы записать в файл введите `yes`: ")

	if fmt.Scan(&agree); agree == "yes" {
		fmt.Println("Вы согласны")
		deletingMap.WriteQueueIntoFile()
	} else {
		fmt.Println("Вы не согласны")
		deletingMap.ShowDeleteQueueInConsole()
	}
}

func CheckParsingDirectory (
	checkDir string,
	adapterDatabase *database.AdapterDatabase,
	xmlMaps database.XmlMapType,
	deletingMap *DeletingMapType,
	wg *sync.WaitGroup,
	) {

	currentDir := file_system.GetCurrentDir()
	checkDirPath := fmt.Sprintf("%v/%v", currentDir, checkDir)

	//директории с файлового хранилища. их имя ID записи в xml
	directoryMaps := file_system.GetParserXmlDirectory(checkDirPath)

	for fileName, _ := range directoryMaps {
		if _, ok := xmlMaps[fileName]; !ok {
			parserId := GetParserId(fileName)

			countItems := adapterDatabase.CountItemsByParserId("items", parserId)
			countItemsArchive := adapterDatabase.CountItemsByParserId("items_archive", parserId)

			if countItems == 0 && countItemsArchive == 0 {
				parserPath := fmt.Sprintf("%v/%v/", checkDirPath, fileName)
				deletingMap.Append(checkDir, fileName, parserPath)
			} else {
				fmt.Printf("Xml NOT deleted %v \n", fileName)
			}
		}
	}

	defer wg.Done()
}

func GetParserId(fileName string) int  {
	parserId, err := strconv.Atoi(fileName)
	if err != nil {
		panic(err.Error())
	}
	return parserId
}

