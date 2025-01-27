package ui

import (
	"io"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

var (
	buff            strings.Builder
	backGroundColor string

	// Color formatter.
	Color = formatters.Register("Color", chroma.FormatterFunc(func(w io.Writer, s *chroma.Style, iterator chroma.Iterator) error {
		for t := iterator(); t != chroma.EOF; t = iterator() {
			colour := s.Get(t.Type).Colour
			backGroundColor = s.Get(t.Type).Background.String()
			var sb strings.Builder

			sb.WriteString("[")
			sb.WriteString(colour.String())
			sb.WriteString("]")
			sb.WriteString(t.Value)

			if _, err := w.Write([]byte(sb.String())); err != nil {
				return err
			}
		}
		return nil
	}))
)

func Colorize(partialFileContents string, fileContents string, filename string, themeName string) {
	buff.Reset()

	// attempt to the language from its filename.
	l := lexers.Match(filename)
	if l == nil {
		// otherwise attempt by its partial contents
		l = lexers.Analyse(partialFileContents)
		if l == nil {
			// otherwise attempt by its entire contents
			l = lexers.Analyse(fileContents)
			if l == nil {
				// otherwise fall back to default
				l = lexers.Fallback
			}
		}
	}
	s := styles.Get(themeName)
	c := Color
	it, err := l.Tokenise(nil, fileContents)
	if err != nil {
		buff.Write([]byte(fileContents))
	}
	c.Format(&buff, s, it)
	return
}
