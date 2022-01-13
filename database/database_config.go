package database

type DatabaseConfig struct {
	Host string
	DatabaseName string
	UserName string
	Password string
	Protocol string
	Driver string
}

func (conf DatabaseConfig) ConnectionString() string {
	return conf.UserName + ":" + conf.Password + "@" + conf.Protocol + "(" + conf.Host + ")/" + conf.DatabaseName
}