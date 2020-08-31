package Settings

type CacheDbConf struct {
	Host        string `default:"127.0.0.1"`
	Port        int    `default:"6379"`
	Password    string `default:""`
	MaxIdle     int    `default:"100"`
	MaxActive   int    `default:"4000"`
	IdleTimeout int    `default:"180"`
	Wait        bool   `default:"true"`
	Database    int    `default:"0"`
}
