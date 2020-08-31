package Settings

type DocDbConf struct {
	Protocol          string `envField:"docdb:Protocol"`
	Host              string `envField:"docdb:Host"`
	DefaultDb         string `envField:"docdb:DefaultDb"`
	Username          string `envField:"docdb:Username"`
	Password          string `envField:"docdb:Password"`
	ReplicaSet        string `envField:"docdb:ReplicaSet"`
	ReadPreference    string `envField:"docdb:ReadPreference"`
	ConnectTimeoutMs  int    `envField:"docdb:ConnectTimeoutMs"`
	SocketTimeoutMs   int    `envField:"docdb:SocketTimeoutMs"`
	ReconnectInterval int    `envField:"docdb:ReconnectInterval"`
	PoolSize          int    `envField:"docdb:PoolSize"`
	BufferMaxEntries  int    `envField:"docdb:BufferMaxEntries"`
	KeepAlive         bool   `envField:"docdb:KeepAlive"`
	BufferCommands    bool   `envField:"docdb:BufferCommands"`
	AutoReconnect     bool   `envField:"docdb:AutoReconnect"`
	SSL               bool   `envField:"docdb:SSL"`
	CaFilePath        string `envField:"docdb:CaFilePath"`
}
