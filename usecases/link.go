package usecases

type LinkUseCasesInterface interface {
	CutLink(link string) (string, error)
	// todo: add more functions for link use cases
}

type LinkUseCases struct{}

func (LinkUseCases) CutLink(link string) (string, error) {
	panic("implement me")
}