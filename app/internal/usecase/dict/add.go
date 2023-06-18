package dict

import (
	"github.com/chun37/greenland-yomiage/internal/dictionary"
	"golang.org/x/xerrors"
)

type Dependencies struct {
	Dictionary dictionary.Service
}

type AddUsecase struct {
	deps Dependencies
}

func NewAddUsecase(deps Dependencies) *AddUsecase {
	return &AddUsecase{
		deps: deps,
	}
}

func (u *AddUsecase) Do(word, yomi string, accent int) error {
	if err := u.deps.Dictionary.Add(word, yomi, accent); err != nil {
		return xerrors.Errorf("cannot add word to dictionary: %w", err)
	}
	return nil
}
