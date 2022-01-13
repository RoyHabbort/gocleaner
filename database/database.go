package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type CacheCountType map[string]int

type Xml struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type AdapterDatabase struct {
	DB *sql.DB
	CacheCount CacheCountType
}

func (adapter AdapterDatabase) Cached(query string, count int) {
	adapter.CacheCount[query] = count
}

func (adapter AdapterDatabase) GetByCache(query string) (int, bool) {
	cnt, ok := adapter.CacheCount[query]
	return cnt, ok
}

type XmlMapType map[string]Xml

func (adapter AdapterDatabase) GetParserXmlActiveRecords() XmlMapType {
	results, err := adapter.DB.Query("SELECT id, name, url FROM xml WHERE status = 1")
	Pause()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer results.Close()

	xmlMaps := make(XmlMapType)
	for results.Next() {
		var xml Xml
		// for each row, scan the result into our tag composite object
		err = results.Scan(&xml.ID, &xml.Name, &xml.Url)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		xmlMaps[xml.ID] = xml
	}

	countActiveXmlMap := len(xmlMaps)
	if countActiveXmlMap == 0 {
		panic("Failed script. xml from db is empty")
	}

	return xmlMaps
}

func (adapter AdapterDatabase) CountItemsByParserId(tableName string, parserId int) int {
	query := fmt.Sprintf("SELECT count(id) FROM %v WHERE parser_id = %v", tableName, parserId)
	if cnt, ok := adapter.GetByCache(query); ok {
		return cnt
	}

	fmt.Printf("Do query: %v\n", query)
	result, err := adapter.DB.Query(query)
	Pause()

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var count int
	result.Next()
	if errScan := result.Scan(&count); errScan != nil {
		panic(err.Error())
	}

	adapter.Cached(query, count)

	return count
}

func InitDatabase(databaseConfig DatabaseConfig) AdapterDatabase {
	db := ConnectDatabase(databaseConfig)
	return AdapterDatabase{DB: db, CacheCount: CacheCountType{}}
}

func ConnectDatabase(databaseConfig DatabaseConfig) *sql.DB {
	db, err := sql.Open(databaseConfig.Driver, databaseConfig.ConnectionString())
	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	return db
}

func Pause() {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(1000)
	i = i+100
	time.Sleep(time.Duration(i) * time.Millisecond)
}
