package database

//NoSQLDbHandler expõe os métodos disponíveis para as abstrações de banco
type NoSQLDbHandler interface {
	Store(string, interface{}) error
	Update(string, interface{}, interface{}) error
	FindAll(string, interface{}, interface{}) error
	FindOne(string, interface{}, interface{}, interface{}) error
}
