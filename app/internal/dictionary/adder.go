package dictionary

type Service interface {
	Add(word, yomi string, accent int) error
}
