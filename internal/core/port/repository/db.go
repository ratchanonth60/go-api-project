package repository

type DB interface {
	Connect() error
}
