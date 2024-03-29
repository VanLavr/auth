package usecase

type Usecase interface {
	GetNewTokenPair()
	RefreshToken()
}
