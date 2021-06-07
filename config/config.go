package config

var Config Configure

type Configure struct {
    Server      Server
    Log         LogConfig
    MysqlMaster MysqlConfig
}

type MysqlConfig struct {
    Host        string
    Password    string
    User        string
    Databases   Database
    MaxLifetime int
    MaxIdleConn int
    MaxOpenConn int
    Charset     string
}

type Database struct {
    Dss string
    Tag string
}

type Server struct {
    Env      string
    RunMode  string
    HttpPort string
}

type LogConfig struct {
    LogLevel   string
    LogFile    string
    MaxSize    int
    MaxBackups int
    MaxAge     int
    Compress   bool
}

