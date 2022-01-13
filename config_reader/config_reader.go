package config_reader

import (
	"../database"
	"fmt"
	"gopkg.in/ini.v1"
)

type IniFileString string

func (file IniFileString) Load() IniFile {
	cfg, err := ini.Load(string(file))
	if err != nil {
		panic(err.Error())
	}
	return IniFile{cfg}
}

type IniFile struct {
	file *ini.File
}

func (iniFile IniFile) DatabaseConfig() database.DatabaseConfig {
	return database.DatabaseConfig{
		iniFile.String("database", "host", true),
		iniFile.String("database", "database", true),
		iniFile.String("database", "user", true),
		iniFile.String("database", "password", true),
		iniFile.String("database", "protocol", true),
		iniFile.String("database", "driver", true),
	}
}

func (iniFile IniFile) String(section string, key string, required bool) string {
	value := iniFile.file.Section(section).Key(key).String()

	if required && value == "" {
		panic(fmt.Sprintf("Empty value for %v %v", section, key))
	}

	return value
}
