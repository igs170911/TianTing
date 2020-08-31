package Settings

type AppConf struct {
	CodeName       string `default:""`
	AppMode        string `default:"debug"`
	SSL            bool   `default:"false"`
	Domain         string `default:""`
	AdminEmail     string `default:""`
	LogLevel       string `default:"info"`
	RpcCommandMode bool   `default:"false"`
	RpcBindPort    int    `default:"9999"`
	RpcEndpoint    string `default:"0.0.0.0"`
	HttpPort       int    `default:"8080"`
	ReadTimeout    int64  `default:"2000"`
	WriteTimeout   int64  `default:"15000"`
}
