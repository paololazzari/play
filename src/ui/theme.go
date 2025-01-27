package ui

import (
	"github.com/gdamore/tcell/v2"
)

type Theme struct {
	BackGroundColor tcell.Color
	TitleColor      tcell.Color
	KeywordColor    tcell.Color
	TextColor       tcell.Color
	BorderColor     tcell.Color
}

var Themes = map[string]Theme{
	"monokai": {
		BackGroundColor: tcell.GetColor("#272822"),
		TitleColor:      tcell.GetColor("#66d9ef"),
		KeywordColor:    tcell.GetColor("#f92672"),
		TextColor:       tcell.GetColor("#e6db74"),
		BorderColor:     tcell.GetColor("#f8f8f2"),
	},
}
