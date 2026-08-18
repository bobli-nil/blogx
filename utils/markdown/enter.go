package markdown

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func MdToHTML(md string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

func ExtractContent(content string, length int) (abs string, err error) {
	htmlStr := MdToHTML(content)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htmlStr)))
	if err != nil {
		return
	}
	htmlText := doc.Text()
	if len(htmlText) > length {
		abs = string([]rune(htmlText)[:length])
	}
	return
}
