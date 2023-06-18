package wavgenerator

type Service interface {
	Generate(s string) ([]byte, error)
}
