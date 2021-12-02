package source_music

type ISourceMusic interface {
	GetHot100Songs(date string) ([]string, error)
}
