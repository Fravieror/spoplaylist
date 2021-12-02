package source_music

type ISourceMusic interface {
	GetHot100Songs(date string) ([]string, error)
}

type MockSourceMusic struct {
	MockGetHot100Songs func(date string) ([]string, error)
}

func (m MockSourceMusic) GetHot100Songs(date string) ([]string, error) {
	return m.MockGetHot100Songs(date)
}
