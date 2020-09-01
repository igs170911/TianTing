package Settings

type DocDbConf struct {
	Protocol          string `default:""`
	Host              string `default:""`
	DefaultDb         string `default:""`
	Username          string `default:""`
	Password          string `default:""`
	ReplicaSet        string `default:""`
	ReadPreference    string `default:""`
	ConnectTimeoutMs  int    `default:""`
	SocketTimeoutMs   int    `default:""`
	ReconnectInterval int    `default:""`
	PoolSize          int    `default:""`
	BufferMaxEntries  int    `default:""`
	KeepAlive         bool   `default:""`
	BufferCommands    bool   `default:""`
	AutoReconnect     bool   `default:""`
	SSL               bool   `default:""`
	CaFilePath        string `default:""`
}
