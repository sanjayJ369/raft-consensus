package types

type DB interface {
	Get(string) any
	Put(string, any) error
	Del(string) error
}
