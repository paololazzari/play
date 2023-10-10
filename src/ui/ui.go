package ui

import (
	"bufio"

	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
	program "github.com/paololazzari/play/src/util"
	"github.com/rivo/tview"
	"golang.org/x/exp/slices"
)

// Colors for UI
const (
	playTitleColor  = tcell.ColorSteelBlue
	playBorderColor = tcell.ColorSteelBlue
	titleColor      = tcell.ColorCadetBlue
	borderColor     = tcell.ColorCadetBlue
)

// User interface
type UI struct {
	App                   *tview.Application
	Label                 string
	CommandText           *tview.TextView
	OptionsInput          *tview.InputField
	EndOptionsText        *tview.TextView
	OpeningQuoteText      *tview.TextView
	ArgumentsInput        *tview.InputField
	ArgumentsInputWide    *tview.TextArea
	ClosingQuoteText      *tview.TextView
	PaddingText           *tview.TextView
	FileOptionsStdin      string
	FileOptionsText       *tview.TextView
	FileOptionsTreeNode   *tview.TreeNode
	FileOptionsTreeView   *tview.TreeView
	FileOptionsInputMap   map[string]bool
	FileOptionsInputSlice []string
	OutputView            *tview.TextView
	FileView              *tview.TextView
	Flex                  *tview.Flex
	ActiveInput           **tview.InputField
}

// Returns the TextView with the command itself
func commandText(label string) *tview.TextView {
	return tview.NewTextView().
		SetText(" > " + label).
		SetTextColor(titleColor)
}

// Returns the InputField for the command options
func optionsInput() *tview.InputField {
	return tview.NewInputField().
		SetLabel("").
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetPlaceholder("<command options>").
		SetPlaceholderTextColor(tcell.ColorDefault).
		SetPlaceholderStyle(tcell.StyleDefault)
}

// Returns the TextView with the separator
func endOptionsText() *tview.TextView {
	return tview.NewTextView().
		SetText(" -- ").
		SetTextColor(titleColor)
}

// Returns the TextView with the opening single quote
func openingQuoteText() *tview.TextView {
	return tview.NewTextView().
		SetText("'")
}

// Returns the InputField for the positional arguments
func argumentsInput() *tview.InputField {
	return tview.NewInputField().
		SetPlaceholder("<positional arguments>").
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetPlaceholderTextColor(tcell.ColorDefault).
		SetPlaceholderStyle(tcell.StyleDefault)
}

// Returns the TextArea for the positional arguments
func argumentsInputWide() *tview.TextArea {
	t := tview.NewTextArea()
	t.SetBorder(true)
	t.SetTitle(" Positional arguments ")
	t.SetTitleColor(titleColor)
	t.SetBorderColor(borderColor)
	return t
}

// Returns the TextView for the closing single quote
func closingQuoteText() *tview.TextView {
	return tview.NewTextView().
		SetText("'")
}

// Returns the TextView used for padding
func paddingText() *tview.TextView {
	return tview.NewTextView()
}

// Returns the TextView for the input files
func fileOptionsText() *tview.TextView {
	return tview.NewTextView().
		SetText("<input files>")
}

// Returns the TreeNode used for file picker
func fileOptionsTreeNode() *tview.TreeNode {
	return tview.NewTreeNode(".")
}

// Returns the TreeView used for file picker
func fileOptionsTreeView() *tview.TreeView {
	t := tview.NewTreeView()
	t.SetBorder(true)
	t.SetTitle(" File picker ")
	t.SetTitleColor(titleColor)
	t.SetBorderColor(borderColor)
	return t
}

// Returns the TextView used for output
func outputView() *tview.TextView {
	t := tview.NewTextView().
		SetDynamicColors(true)
	t.SetBorder(true)
	t.SetTitle(" Output ")
	t.SetTitleColor(titleColor)
	t.SetBorderColor(borderColor)
	return t
}

// Returns the TextView used for file view
func fileView() *tview.TextView {
	t := tview.NewTextView().
		SetDynamicColors(true)
	t.SetBorder(true)
	t.SetTitle(" File view ")
	t.SetTitleColor(titleColor)
	t.SetBorderColor(borderColor)
	return t
}

// Returns the Flex used for layout
func flex() *tview.Flex {
	return tview.NewFlex()
}

// Helper function to keep unique, sorted list of selected files
func updateFileOptionsInput(m map[string]bool, a *[]string, file string) {
	if m[file] {
		delete(m, file)
		for i, v := range *a {
			if file == v {
				*a = slices.Delete(*a, i, i+1)
			}
		}
	} else {
		*a = append(*a, file)
		m[file] = true
	}
}

// Helper function to return the string used as file inputs
func getFileOptionsText(a *[]string) string {
	var sb strings.Builder

	for _, v := range *a {
		sb.WriteString(v)
		sb.WriteString(" ")
	}
	return strings.TrimSpace(sb.String())
}

// UI constructor
func NewUI(label string, stdin string) *UI {
	ui := &UI{
		App:                   tview.NewApplication(),
		Label:                 label,
		CommandText:           commandText(label),
		OptionsInput:          optionsInput(),
		EndOptionsText:        endOptionsText(),
		OpeningQuoteText:      openingQuoteText(),
		ArgumentsInput:        argumentsInput(),
		ArgumentsInputWide:    argumentsInputWide(),
		ClosingQuoteText:      closingQuoteText(),
		PaddingText:           paddingText(),
		FileOptionsStdin:      stdin,
		FileOptionsText:       fileOptionsText(),
		FileOptionsTreeNode:   fileOptionsTreeNode(),
		FileOptionsTreeView:   fileOptionsTreeView(),
		FileOptionsInputMap:   make(map[string]bool),
		FileOptionsInputSlice: []string{},
		OutputView:            outputView(),
		FileView:              fileView(),
		Flex:                  flex(),
		ActiveInput:           nil,
	}
	return ui
}

// Helper function for evaluating expressions
func (ui *UI) evaluateExpression() func() {
	return func() {
		var sb strings.Builder
		sb.WriteString(ui.Label)
		sb.WriteString(" ")
		sb.WriteString(ui.OptionsInput.GetText())
		sb.WriteString(" -- ")
		sb.WriteString(ui.OpeningQuoteText.GetText(false))
		sb.WriteString(ui.ArgumentsInput.GetText())
		sb.WriteString(ui.ClosingQuoteText.GetText(false))
		sb.WriteString(" ")
		if t := ui.FileOptionsText.GetText(false); t != "<input files>" {
			if len(ui.FileOptionsStdin) > 0 && len(ui.FileOptionsInputSlice) == 0 {
				sb.WriteString("<<<")
				sb.WriteString("'")
				sb.WriteString(ui.FileOptionsStdin)
				sb.WriteString("'")
			} else {
				sb.WriteString(t)
			}
		}
		out, _ := program.Run(sb.String())
		ui.OutputView.SetText(out)
	}
}

// Callback function for InputField
func (ui *UI) changedInputField() func(string) {
	return func(text string) {
		go ui.App.QueueUpdateDraw(ui.evaluateExpression())
	}
}

// Callback function for TextView
func (ui *UI) changedText() func() {
	return func() {
		go ui.App.QueueUpdateDraw(ui.evaluateExpression())
	}
}

// Helper function for populating nodes of TreeNode
func add(target *tview.TreeNode, path string, ui *UI) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		node := tview.NewTreeNode(file.Name())

		// set custom reference with full path
		node.SetReference(filepath.Join(path, file.Name()))

		if file.IsDir() {
			node.SetColor(tcell.ColorGreen)
			target.AddChild(node)
		} else {
			// only add text files https://stackoverflow.com/a/75940070/3390419
			f, _ := os.Open(file.Name())
			scanner := bufio.NewScanner(f)
			scanner.Split(bufio.ScanLines)
			scanner.Scan()
			if utf8.ValidString(string(scanner.Text())) {
				target.AddChild(node)
			}
		}
	}
}

// Function for configuring OptionsInput InputField
func (ui *UI) configOptionsInput() {
	ui.OptionsInput.SetChangedFunc(ui.changedInputField())
	ui.OptionsInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		switch key {
		case tcell.KeyTab:
			ui.App.SetFocus(ui.ArgumentsInput)
		case tcell.KeyBacktab:
			ui.App.SetFocus(ui.FileOptionsTreeView)
		case tcell.KeyEnter:
			ui.ActiveInput = &ui.OptionsInput
			ui.App.SetFocus(ui.OutputView)
		case tcell.KeyRune:
			ui.OutputView.ScrollToBeginning()
		}
		return event
	})
}

// Function for configuring ArgumentsInput InputField
func (ui *UI) configArgumentsInput() {
	ui.ArgumentsInput.SetChangedFunc(ui.changedInputField())
	ui.ArgumentsInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		switch key {
		case tcell.KeyRune:
			ui.OutputView.ScrollToBeginning()
		case tcell.KeyTab:
			ui.App.SetFocus(ui.FileOptionsTreeView)
		case tcell.KeyBacktab:
			ui.App.SetFocus(ui.OptionsInput)
		case tcell.KeyEnter:
			ui.ActiveInput = &ui.ArgumentsInput
			ui.App.SetFocus(ui.OutputView)
		case tcell.KeyCtrlSpace:
			if ui.OpeningQuoteText.GetText(false) == "'" {
				ui.OpeningQuoteText.SetText("\"")
				ui.ClosingQuoteText.SetText("\"")
			} else {
				ui.OpeningQuoteText.SetText("'")
				ui.ClosingQuoteText.SetText("'")
			}
		case tcell.KeyCtrlO:
			ui.ArgumentsInputWide.SetText(ui.ArgumentsInput.GetText(), false)
			ui.App.SetRoot(ui.ArgumentsInputWide, true).
				SetFocus(ui.ArgumentsInputWide)
		}
		return event
	})

}

// Function for configuring ArgumentsInputWide InputField
func (ui *UI) configArgumentsInputWide() {
	ui.ArgumentsInputWide.SetChangedFunc(ui.changedText())

	ui.ArgumentsInputWide.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		switch key {
		case tcell.KeyCtrlO:
			ui.App.SetRoot(ui.Flex, true).
				SetFocus(ui.ArgumentsInput)
			ui.ArgumentsInput.SetText(ui.ArgumentsInputWide.GetText())
		}
		return event
	})
	ui.ArgumentsInputWide.SetBorder(true)
	ui.ArgumentsInputWide.SetBorderColor(playBorderColor)
}

// Function for configuring FileOptionsText TextView
func (ui *UI) configFileOptionsInput() {
	ui.FileOptionsText.SetChangedFunc(ui.changedText())

	if len(ui.FileOptionsStdin) > 0 {
		ui.FileOptionsText.SetText(ui.FileOptionsStdin)
	}
	ui.FileOptionsText.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		switch key {
		case tcell.KeyTab:
			ui.App.SetFocus(ui.OptionsInput)
		case tcell.KeyBacktab:
			ui.App.SetFocus(ui.ArgumentsInput)
		case tcell.KeyEnter:
			ui.App.SetFocus(ui.OutputView)
		case tcell.KeyDown:
			ui.App.SetFocus(ui.FileOptionsTreeView)
		case tcell.KeyRune:
			return nil
		}
		return event
	})
}

// Function for configuring FileOptionsTreeNode TreeNode
func (ui *UI) configFileOptionsTreeNode() {
	rootDir := "."
	ui.FileOptionsTreeNode = tview.NewTreeNode(rootDir)
	add(ui.FileOptionsTreeNode, rootDir, ui)
}

// Function for configuring FileOptionsTreeView TreeView
func (ui *UI) configFileOptionsTreeView() {
	ui.FileOptionsTreeView.SetRoot(ui.FileOptionsTreeNode).SetCurrentNode(ui.FileOptionsTreeNode)
	defaultColor := ui.FileOptionsTreeNode.GetColor()

	ui.FileOptionsTreeView.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()

		if reference == nil {
			return
		}
		if stat, _ := os.Stat(reference.(string)); !stat.IsDir() {
			if node.GetColor() == tcell.ColorRed {
				node.SetColor(defaultColor)
			} else {
				node.SetColor(tcell.ColorRed)
			}
			// when a file is selected, update the sorted, unique list of files
			updateFileOptionsInput(ui.FileOptionsInputMap, &ui.FileOptionsInputSlice, reference.(string))
			ui.FileOptionsText.SetText(getFileOptionsText(&ui.FileOptionsInputSlice))
			ui.OutputView.ScrollToBeginning()
			return
		}
		children := node.GetChildren()
		if len(children) == 0 {
			path := reference.(string)
			add(node, path, ui)
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})

	ui.FileOptionsTreeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		switch key {
		case tcell.KeyEsc:
			ui.App.SetFocus(ui.ArgumentsInput)
		case tcell.KeyBacktab:
			ui.App.SetFocus(ui.ArgumentsInput)
		case tcell.KeyTab:
			ui.App.SetFocus(ui.OptionsInput)

		case tcell.KeyCtrlO:
			if ui.FileOptionsTreeView.GetCurrentNode() == ui.FileOptionsTreeView.GetRoot() {
				return event
			}
			file, err := os.ReadFile(ui.FileOptionsTreeView.GetCurrentNode().GetReference().(string))
			if err == nil {
				ui.FileView.SetText(string(file))
				ui.App.SetRoot(ui.FileView, true).
					SetFocus(ui.FileView)
			}
		}

		return event
	})
}

// Function for configuring OutputView TextView
func (ui *UI) configOutputView() {
	ui.OutputView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			ui.App.SetFocus(*ui.ActiveInput)
		}
		return event
	})
}

// Function for configuring FileView TextView
func (ui *UI) configFileView() {
	ui.FileView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		switch key {
		case tcell.KeyCtrlO:
			ui.App.SetRoot(ui.Flex, true).
				SetFocus(ui.FileOptionsTreeView)
		}
		return event
	})
}

// Function for configuring Flex Flex
func (ui *UI) configFlex() {
	ui.Flex.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 2, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(ui.CommandText, len(ui.Label)+4, 1, false).
			AddItem(ui.OptionsInput, 17, 1, false).
			AddItem(ui.EndOptionsText, 4, 1, false).
			AddItem(ui.OpeningQuoteText, 1, 1, false).
			AddItem(ui.ArgumentsInput, 22, 1, false).
			AddItem(ui.ClosingQuoteText, 1, 1, false).
			AddItem(ui.PaddingText, 1, 1, false).
			AddItem(ui.FileOptionsText, 0, 1, false).
			AddItem(tview.NewBox(), 2, 1, false), 3, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(ui.OutputView, 0, 10, false).
			AddItem(ui.FileOptionsTreeView, 0, 2, false), 0, 1, false), 0, 1, false)
	ui.Flex.SetBorder(true)
	ui.Flex.SetTitle(" play ")
	ui.Flex.SetTitleColor(playTitleColor)
	ui.Flex.SetBorderColor(playBorderColor)

}

// Initialize UI
func (ui *UI) InitUI() error {

	ui.configOptionsInput()
	ui.configArgumentsInput()
	ui.configArgumentsInputWide()
	ui.configFileOptionsInput()
	ui.configFileOptionsTreeNode()
	ui.configFileOptionsTreeView()
	ui.configOutputView()
	ui.configFileView()
	ui.configFlex()

	ui.App.SetRoot(ui.Flex, true).
		SetFocus(ui.OptionsInput)

	return nil
}

// Run the application
func (ui *UI) Run() error {
	if err := ui.App.Run(); err != nil {
		panic(err)
	}
	return nil
}
