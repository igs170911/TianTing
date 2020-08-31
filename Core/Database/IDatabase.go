package Database

type IDatabase interface {
	connect()
	GetClient()
}
