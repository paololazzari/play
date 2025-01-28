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
	"base16-snazzy": {
		BackGroundColor: tcell.GetColor("#282a36"),
		TitleColor:      tcell.GetColor("#aa0000"),
		KeywordColor:    tcell.GetColor("#ff6ac1"),
		TextColor:       tcell.GetColor("#5af78e"),
	},
	"dracula": {
		BackGroundColor: tcell.GetColor("#282a36"),
		TitleColor:      tcell.GetColor("#8be9fd"),
		KeywordColor:    tcell.GetColor("#ff79c6"),
		TextColor:       tcell.GetColor("#f1fa8c"),
		BorderColor:     tcell.GetColor("#f8f8f2"),
	},
	"fruity": {
		BackGroundColor: tcell.GetColor("#111111"),
		TitleColor:      tcell.GetColor("#ff0086"),
		KeywordColor:    tcell.GetColor("#fb660a"),
		TextColor:       tcell.GetColor("#0086d2"),
	},
	"monokai": {
		BackGroundColor: tcell.GetColor("#272822"),
		TitleColor:      tcell.GetColor("#66d9ef"),
		KeywordColor:    tcell.GetColor("#f92672"),
		TextColor:       tcell.GetColor("#e6db74"),
		BorderColor:     tcell.GetColor("#f8f8f2"),
	},
	"vim": {
		BackGroundColor: tcell.GetColor("#000000"),
		TitleColor:      tcell.GetColor("#56d364"),
		KeywordColor:    tcell.GetColor("#cd00cd"),
		TextColor:       tcell.GetColor("#cd0000"),
	},
	"witchhazel": {
		BackGroundColor: tcell.GetColor("#433e56"),
		TitleColor:      tcell.GetColor("#c2ffdf"),
		KeywordColor:    tcell.GetColor("#ffb8d1"),
		TextColor:       tcell.GetColor("#1bc5e0"),
		BorderColor:     tcell.GetColor("#f8f8f2"),
	},
}
