package source_music

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anaskhan96/soup"
)

type Billboard struct {
}

func NewBillboard() ISourceMusic {
	return &Billboard{}
}

func (b *Billboard) GetHot100Songs(date string) ([]string, error) {
	url := fmt.Sprintf("https://www.billboard.com/charts/hot-100/%s", date)
	resp, err := soup.Get(url)
	if err != nil {
		fmt.Errorf("error getting from url %s: detail %w", url, err)
		return nil, errors.New("error getting top 100, check logs for more details")
	}
	doc := soup.HTMLParse(resp)
	elements := doc.Find("div", "id", "post-1479786").FindAll("h3", "id", "title-of-a-story")
	songs := make([]string, 0)
	for _, element := range elements {
		song := strings.Replace(element.Text(), "\"", "", -1)
		song = strings.Replace(song, "\n", "", -1)

		if len(song) > 0 {
			fmt.Println(song)
			songs = append(songs, song)
		}
	}
	return songs, nil
}
