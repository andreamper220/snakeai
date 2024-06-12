package storages

var Storage StorageInterface

type StorageInterface interface {
	AddUser(user User) error
}

type User struct {
	Email    string
	Password string
}
