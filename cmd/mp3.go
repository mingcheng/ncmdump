package main

import (
	"github.com/bogem/id3v2"
	"github.com/mingcheng/ncmdump"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type MP3 struct {
	FilePath string
	Meta     ncmdump.Meta
}

func (m *MP3) Update() error {
	tag, err := id3v2.Open(m.FilePath, id3v2.Options{Parse: false})
	if err != nil {
		return err
	}
	defer tag.Close()

	// set default id3 tag encoding as UTF8
	tag.SetDefaultEncoding(id3v2.EncodingUTF8)

	tag.SetTitle(m.Meta.Name)
	tag.SetAlbum(m.Meta.Album.Name)
	tag.SetArtist(func(artists []ncmdump.Artist) string {
		var tmp string
		for _, v := range artists {
			tmp += v.Name
		}
		return strings.TrimSpace(tmp)
	}(m.Meta.Artists))

	// fetch cover images from url
	if u, err := url.Parse(m.Meta.Album.CoverUrl); err == nil {
		pic, err := http.Get(u.String())
		if err == nil && pic.StatusCode == http.StatusOK {
			if artwork, err := ioutil.ReadAll(pic.Body); err != nil {
				return err
			} else {
				pic := id3v2.PictureFrame{
					Encoding:    id3v2.EncodingUTF8,
					MimeType:    "image/jpeg",
					PictureType: id3v2.PTFrontCover,
					Description: m.Meta.Album.Name,
					Picture:     artwork,
				}

				tag.AddAttachedPicture(pic)
			}
		}
	}

	// save and update tags
	return tag.Save()
}
