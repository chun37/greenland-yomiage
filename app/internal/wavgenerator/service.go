package wavgenerator

import "io"

type Service interface {
	Generate() io.Reader
}
