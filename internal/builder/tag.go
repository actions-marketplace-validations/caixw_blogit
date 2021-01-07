// SPDX-License-Identifier: MIT

package builder

import "github.com/caixw/blogit/internal/data"

type tags struct {
	XMLName struct{} `xml:"tags"`
	Tags    []*tag   `xml:"tag"`
}

type tag struct {
	XMLName struct{} `xml:"tag"`

	Permalink string      `xml:"permalink"`
	Title     string      `xml:"title"`
	Color     string      `xml:"color,attr"`
	Created   datetime    `xml:"created"`
	Modified  datetime    `xml:"modified"`
	Posts     []*postMeta `xml:"post,omitempty"`
	Content   innerhtml   `xml:"summary"`
}

func newTag(t *data.Tag, d *data.Data) *tag {
	ps := make([]*postMeta, 0, len(t.Posts))
	for _, p := range t.Posts {
		ps = append(ps, &postMeta{
			Permalink: d.BuildURL(p.Slug),
			Title:     p.Title,
			Created:   toDatetime(p.Created, d),
			Modified:  toDatetime(p.Modified, d),
			Summary:   innerhtml{Content: p.Summary},
		})
	}

	return &tag{
		Permalink: d.BuildURL(t.Path),
		Title:     t.Title,
		Color:     t.Color,
		Created:   toDatetime(t.Created, d),
		Modified:  toDatetime(t.Modified, d),
		Posts:     ps,
		Content:   innerhtml{Content: t.Content},
	}
}

func (b *Builder) buildTags(d *data.Data) error {
	tags := make([]*tag, 0, len(d.Tags))

	for _, t := range d.Tags {
		tt := newTag(t, d)
		xsl := d.BuildThemeURL("tag.xsl")
		if err := b.appendXMLFile(t.Path, xsl, t.Modified, tt); err != nil {
			return err
		}

		tt.Posts = nil
		tags = append(tags, tt)
	}

	xsl := d.BuildThemeURL("tags.xsl")
	return b.appendXMLFile("tags.xml", xsl, d.Modified, tags)
}
