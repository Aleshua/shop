package ports

type ITokenGeneratorService interface {
	// length - количество симолов
	GenerateToken(length int) (string, error)
}
