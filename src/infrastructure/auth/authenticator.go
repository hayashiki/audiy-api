package auth

// AuthService TODO: idToken検証を切り出す
type AuthService interface {
	Verify(string) (uid string, err error)
}

type authService struct {}

func (a authService) Verify(token string) {

}
