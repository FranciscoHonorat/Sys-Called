package out

type PasswordGenerator interface {
	Generate() (string, error)
}
