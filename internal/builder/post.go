// SPDX-License-Identifier: MIT

package builder

import "github.com/caixw/blogit/internal/data"

type posts struct {
	XMLName struct{}    `xml:"posts"`
	Posts   []*postMeta `xml:"post"`
}

type postMeta struct {
	Permalink string     `xml:"permalink,attr"`
	Title     string     `xml:"title"`
	Created   *datetime  `xml:"created,omitempty"`
	Modified  *datetime  `xml:"modified,omitempty"`
	Tags      []*tagMeta `xml:"tag,omitempty"`
	Summary   innerhtml  `xml:"summary,omitempty"`
}

type tagMeta struct {
	Permalink string `xml:"permalink,attr"`
	Title     string `xml:"title"`
}

type post struct {
	XMLName   struct{}   `xml:"post"`
	Permalink string     `xml:"permalink,attr"`
	Title     string     `xml:"title"`
	Created   *datetime  `xml:"created,omitempty"`
	Modified  *datetime  `xml:"modified,omitempty"`
	Tags      []*tagMeta `xml:"tag"`
	Language  string     `xml:"language,attr"`
	Outdated  *outdated  `xml:"outdated,omitempty"`
	Authors   []*author  `xml:"author"`
	License   *link      `xml:"license"`
	Summary   innerhtml  `xml:"summary,omitempty"`
	Content   innerhtml  `xml:"content"`
	Prev      *link      `xml:"prev"`
	Next      *link      `xml:"next"`
}

type author struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url,omitempty"`
	Email  string `yaml:"email,omitempty"`
	Avatar string `yaml:"avatar,omitempty"`
}

type link struct {
	URL   string `xml:"url,attr"`
	Title string `xml:"title,attr,omitempty"`
	Text  string `xml:"text"`
}

type outdated struct {
	Outdated *datetime `xml:"outdated,omitempty"` // 过期的时间
	Content  string    `xml:",chardata"`
}

func (b *Builder) buildPosts(d *data.Data) error {
	index := &posts{Posts: make([]*postMeta, 0, len(d.Posts))}

	for _, p := range d.Posts {
		tags := make([]*tagMeta, 0, len(p.Tags))
		for _, t := range p.Tags {
			tags = append(tags, &tagMeta{
				Permalink: d.BuildURL(t.Path),
				Title:     t.Title,
			})
		}

		authors := make([]*author, 0, len(p.Authors))
		for _, a := range p.Authors {
			authors = append(authors, &author{
				Name:   a.Name,
				URL:    a.URL,
				Email:  a.Email,
				Avatar: a.Avatar,
			})
		}

		var od *outdated
		if p.Outdated != nil {
			od = &outdated{
				Content:  p.Outdated.Content,
				Outdated: newDatetime(p.Outdated.Outdated, d),
			}
		}

		pp := &post{
			Permalink: d.BuildURL(p.Path),
			Title:     p.Title,
			Created:   newDatetime(p.Created, d),
			Modified:  newDatetime(p.Modified, d),
			Tags:      tags,
			Language:  p.Language,
			Outdated:  od,
			Authors:   authors,
			License: &link{
				URL:   p.License.URL,
				Title: p.License.Title,
				Text:  p.License.Text,
			},
			Content: innerhtml{Content: p.Content},
			Summary: innerhtml{Content: p.Summary},
		}
		if p.Prev != nil {
			pp.Prev = &link{
				URL:   d.BuildURL(p.Prev.Path),
				Title: "上一篇文章",
				Text:  p.Prev.Title,
			}
		}
		if p.Next != nil {
			pp.Next = &link{
				URL:   d.BuildURL(p.Next.Path),
				Title: "后一篇文章",
				Text:  p.Next.Title,
			}
		}
		err := b.appendXMLFile(p.Slug+".xml", d.BuildThemeURL(p.Template), p.Modified, pp)
		if err != nil {
			return err
		}

		index.Posts = append(index.Posts, &postMeta{
			Permalink: d.BuildURL(p.Path),
			Title:     p.Title,
			Created:   newDatetime(p.Created, d),
			Modified:  newDatetime(p.Modified, d),
			Tags:      tags,
			Summary:   innerhtml{Content: p.Summary},
		})
	}

	return b.appendXMLFile("index.xml", d.BuildThemeURL("index.xsl"), d.Modified, index)
}
