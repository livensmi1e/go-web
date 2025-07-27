package ports

type Hasher interface {
	Hash(password string) (string, error)
	Compare(hash string, plain string) error
}
