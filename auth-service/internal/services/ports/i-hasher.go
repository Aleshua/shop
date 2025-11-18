package ports

type IHasherService interface {
	Hash(value string) (string, error)
	Compare(value, hash string) bool
}
