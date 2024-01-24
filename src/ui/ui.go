package ui

import (
	"bufio"
	"fmt"

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
	App                    *tview.Application
	Label                  string
	EndOfOptionsSeparator  bool
	CommandText            *tview.TextView
	OptionsInput           *tview.InputField
	EndOptionsText         *tview.TextView
	OpeningQuoteText       *tview.TextView
	ArgumentsInput         *tview.InputField
	ArgumentsInputWide     *tview.TextArea
	ArgumentsInputWideFlex *tview.Flex
	ClosingQuoteText       *tview.TextView
	EndArgumentsText       *tview.TextView
	FileOptionsStdin       string
	FileOptionsText        *tview.TextView
	FileOptionsTreeNode    *tview.TreeNode
	FileOptionsTreeView    *tview.TreeView
	FileOptionsInputMap    map[string]bool
	FileOptionsInputSlice  []string
	OutputView             *tview.TextView
	FileView               *tview.TextView
	ChildFlex              *tview.Flex
	Flex                   *tview.Flex
	ActiveInput            **tview.InputField
	ActiveFlex             **tview.Flex
}

type nodeReference struct {
	path                string
	partialFileContents string
}

func getNodePath(node *tview.TreeNode) string {
	ref := node.GetReference()
	return ref.(nodeReference).path
}

func getNodePartialFileContents(node *tview.TreeNode) string {
	ref := node.GetReference()
	return ref.(nodeReference).partialFileContents
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

// Returns the TextView for the options separator
func endOptionsText() *tview.TextView {
	return tview.NewTextView().
		SetTextColor(titleColor)
}

// Returns the TextView for the arguments separator
func endArgumentsText() *tview.TextView {
	return tview.NewTextView().
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

// Returns the Flex used for wide positional arguments
func argumentsInputWideFlex() *tview.Flex {
	return tview.NewFlex()
}

// Returns the TextView for the closing single quote
func closingQuoteText() *tview.TextView {
	return tview.NewTextView().
		SetText("'")
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
func childFlex() *tview.Flex {
	return tview.NewFlex()
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
func NewUI(program string, respectsEndOfOptions bool, stdin string) *UI {
	ui := &UI{
		App:                    tview.NewApplication(),
		Label:                  program,
		EndOfOptionsSeparator:  respectsEndOfOptions,
		CommandText:            commandText(program),
		OptionsInput:           optionsInput(),
		EndOptionsText:         endOptionsText(),
		OpeningQuoteText:       openingQuoteText(),
		ArgumentsInput:         argumentsInput(),
		ArgumentsInputWide:     argumentsInputWide(),
		ArgumentsInputWideFlex: argumentsInputWideFlex(),
		ClosingQuoteText:       closingQuoteText(),
		EndArgumentsText:       endArgumentsText(),
		FileOptionsStdin:       stdin,
		FileOptionsText:        fileOptionsText(),
		FileOptionsTreeNode:    fileOptionsTreeNode(),
		FileOptionsTreeView:    fileOptionsTreeView(),
		FileOptionsInputMap:    make(map[string]bool),
		FileOptionsInputSlice:  []string{},
		OutputView:             outputView(),
		FileView:               fileView(),
		ChildFlex:              childFlex(),
		Flex:                   flex(),
		ActiveInput:            nil,
		ActiveFlex:             nil,
	}
	return ui
}

// Helper function for getting active input text
func (ui *UI) getActiveInputText() string {
	if ui.ActiveFlex == &ui.Flex {
		return ui.ArgumentsInput.GetText()
	} else {
		return ui.ArgumentsInputWide.GetText()
	}
}

// Helper function for evaluating expressions
func (ui *UI) evaluateExpression() func() {
	return func() {
		var sb strings.Builder
		sb.WriteString(ui.Label)
		sb.WriteString(" ")
		sb.WriteString(ui.OptionsInput.GetText())
		if ui.EndOfOptionsSeparator {
			sb.WriteString(" -- ")
		} else {
			sb.WriteString(" ")
		}
		sb.WriteString(ui.OpeningQuoteText.GetText(false))
		sb.WriteString(ui.getActiveInputText())
		sb.WriteString(ui.ClosingQuoteText.GetText(false))
		sb.WriteString(" ")
		if !ui.EndOfOptionsSeparator {
			sb.WriteString(" -- ")
		}
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

// Helper function to resize flex based on argument input size
func (ui *UI) resizeChildFlexIfNeeded() {
	argumentsInputLength := len(ui.ArgumentsInput.GetText())
	if argumentsInputLength >= 40 {
		ui.ChildFlex.ResizeItem(ui.ArgumentsInput, 0, 3)
	} else if argumentsInputLength > 19 && argumentsInputLength < 40 {
		ui.ChildFlex.ResizeItem(ui.ArgumentsInput, 0, 1)
	} else if argumentsInputLength <= 19 {
		ui.ChildFlex.ResizeItem(ui.ArgumentsInput, 22, 1)
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

// Helper function to exclude certain file types from file picker
func isExtensionInvalid(fileExtension string) bool {
	invalidExtensions := [...]string{".gif", ".png"}
	for _, invalidExtension := range invalidExtensions {
		if invalidExtension == fileExtension {
			return true
		}
	}
	return false
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
		nodeRef := nodeReference{filepath.Join(path, file.Name()), ""}
		node.SetReference(nodeRef)

		if file.IsDir() {
			node.SetColor(tcell.ColorGreen)
			target.AddChild(node)
		} else {
			// only add text files https://stackoverflow.com/a/75940070/3390419
			f, _ := os.Open(file.Name())
			scanner := bufio.NewScanner(f)
			scanner.Split(bufio.ScanLines)
			scanner.Scan()
			text := string(scanner.Text())
			fileExtension := filepath.Ext(file.Name())
			if utf8.ValidString(text) && !isExtensionInvalid(fileExtension) {
				target.AddChild(node)
				nodeRef := nodeReference{filepath.Join(path, file.Name()), text}
				node.SetReference(nodeRef)
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
			ui.resizeChildFlexIfNeeded()
		case tcell.KeyDelete:
			ui.OutputView.ScrollToBeginning()
			ui.resizeChildFlexIfNeeded()
		case tcell.KeyBackspace2:
			ui.OutputView.ScrollToBeginning()
			ui.resizeChildFlexIfNeeded()
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
			ui.ArgumentsInputWide.SetText(ui.ArgumentsInput.GetText(), true)
			ui.ActiveFlex = &ui.ArgumentsInputWideFlex
			ui.App.SetRoot(ui.ArgumentsInputWideFlex, true).
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
		// on ctrl+enter return a new enter event
		if event.Modifiers() == 2 && event.Rune() == 10 {
			return tcell.NewEventKey(tcell.KeyEnter, 10, 0)
		}
		switch key {
		case tcell.KeyCtrlO:
			ui.ActiveFlex = &ui.Flex
			ui.App.SetRoot(ui.Flex, true).
				SetFocus(ui.ArgumentsInput)
			ui.ArgumentsInput.SetText(ui.ArgumentsInputWide.GetText())
			ui.resizeChildFlexIfNeeded()
		case tcell.KeyEsc:
			ui.ActiveFlex = &ui.Flex
			ui.App.SetRoot(ui.Flex, true).
				SetFocus(ui.ArgumentsInput)
			ui.ArgumentsInput.SetText(ui.ArgumentsInputWide.GetText())
			ui.resizeChildFlexIfNeeded()
		case tcell.KeyEnter:
			ui.App.SetFocus(ui.OutputView)
			return nil
		}
		return event
	})
	ui.ArgumentsInputWide.SetBorder(true)
	ui.ArgumentsInputWide.SetBorderColor(playBorderColor)
}

// Function for configuring ArgumentsInputWideFlex Flex
func (ui *UI) configArgumentsInputWideFlex() {
	ui.ArgumentsInputWideFlex.SetDirection(tview.FlexRow).
		AddItem(ui.ArgumentsInputWide, 0, 1, false).
		AddItem(ui.OutputView, 0, 1, false)
	ui.ArgumentsInputWideFlex.SetBorder(true)
	ui.ArgumentsInputWideFlex.SetTitle(" play ")
	ui.ArgumentsInputWideFlex.SetTitleColor(playTitleColor)
	ui.ArgumentsInputWideFlex.SetBorderColor(playBorderColor)
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
	nodeRef := nodeReference{"", ""}
	ui.FileOptionsTreeNode.SetReference(nodeRef)
	add(ui.FileOptionsTreeNode, rootDir, ui)
}

// Function for configuring FileOptionsTreeView TreeView
func (ui *UI) configFileOptionsTreeView() {
	ui.FileOptionsTreeView.SetRoot(ui.FileOptionsTreeNode).SetCurrentNode(ui.FileOptionsTreeNode)
	defaultColor := ui.FileOptionsTreeNode.GetColor()

	ui.FileOptionsTreeView.SetSelectedFunc(func(node *tview.TreeNode) {

		nodePath := getNodePath(node)
		if nodePath == "" {
			return
		}
		if stat, _ := os.Stat(nodePath); !stat.IsDir() {
			if node.GetColor() == tcell.ColorRed {
				node.SetColor(defaultColor)
			} else {
				node.SetColor(tcell.ColorRed)
			}
			// when a file is selected, update the sorted, unique list of files
			updateFileOptionsInput(ui.FileOptionsInputMap, &ui.FileOptionsInputSlice, nodePath)
			ui.FileOptionsText.SetText(getFileOptionsText(&ui.FileOptionsInputSlice))
			ui.OutputView.ScrollToBeginning()
			return
		}
		children := node.GetChildren()
		if len(children) == 0 {
			add(node, nodePath, ui)
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
			filename := getNodePath(ui.FileOptionsTreeView.GetCurrentNode())
			file, err := os.ReadFile(filename)
			if err == nil {
				fileContents := string(file)
				Colorize(getNodePartialFileContents(ui.FileOptionsTreeView.GetCurrentNode()), fileContents, filename)
				ui.FileView.SetText(buff.String())
				ui.FileView.SetBackgroundColor(tcell.GetColor(backGroundColor))
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
			if ui.ActiveFlex == &ui.Flex {
				ui.App.SetRoot(ui.Flex, true)
				ui.App.SetFocus(*ui.ActiveInput)
			} else {
				ui.App.SetRoot(ui.ArgumentsInputWideFlex, true)
				ui.App.SetFocus(ui.ArgumentsInputWide)
			}
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
		case tcell.KeyEsc:
			ui.App.SetRoot(ui.Flex, true).
				SetFocus(ui.FileOptionsTreeView)
		}

		return event
	})
}

// Helper function
func (ui *UI) endOptionsSeparator() (*tview.TextView, int, int, bool) {
	if ui.EndOfOptionsSeparator {
		return ui.EndOptionsText.SetText(" -- "), 4, 1, false
	} else {
		return ui.EndOptionsText.SetText(" "), 1, 1, false
	}
}

// Helper function
func (ui *UI) endArgumentsSeparator() (*tview.TextView, int, int, bool) {
	if ui.EndOfOptionsSeparator {
		return ui.EndArgumentsText.SetText(" "), 1, 1, false
	} else {
		return ui.EndArgumentsText.SetText(" -- "), 4, 1, false
	}
}

// Function for configuring ChildFlex Flex
func (ui *UI) configChildFlex() {
	ui.ChildFlex.SetDirection(tview.FlexColumn).
		AddItem(ui.CommandText, len(ui.Label)+4, 1, false).
		AddItem(ui.OptionsInput, 17, 1, false).
		AddItem(ui.endOptionsSeparator()).
		AddItem(ui.OpeningQuoteText, 1, 1, false).
		AddItem(ui.ArgumentsInput, 22, 1, false).
		AddItem(ui.ClosingQuoteText, 1, 1, false).
		AddItem(ui.endArgumentsSeparator()).
		AddItem(ui.FileOptionsText, 0, 1, false).
		AddItem(tview.NewBox(), 2, 1, false)
}

// Function for configuring Flex Flex
func (ui *UI) configFlex() {

	ui.Flex.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 2, 1, false).
		AddItem(ui.ChildFlex, 3, 1, false).
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
	ui.configArgumentsInputWideFlex()
	ui.configFileOptionsInput()
	ui.configFileOptionsTreeNode()
	ui.configFileOptionsTreeView()
	ui.configOutputView()
	ui.configFileView()
	ui.configChildFlex()
	ui.configFlex()

	// on Ctrl+S shut down the application and print the expression to stdout
	ui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		switch key {
		case tcell.KeyCtrlS:
			endOptionsSeparator, _, _, _ := ui.endOptionsSeparator()
			endArgumentsSeparator, _, _, _ := ui.endArgumentsSeparator()

			var sb strings.Builder
			sb.WriteString(ui.Label)
			sb.WriteString(" ")
			sb.WriteString(ui.OptionsInput.GetText())
			sb.WriteString(endOptionsSeparator.GetText(false))
			sb.WriteString(ui.OpeningQuoteText.GetText(false))
			sb.WriteString(ui.getActiveInputText())
			sb.WriteString(ui.ClosingQuoteText.GetText(false))
			sb.WriteString(endArgumentsSeparator.GetText(false))
			sb.WriteString(strings.Join(ui.FileOptionsInputSlice, " "))

			ui.App.Stop()
			fmt.Println(sb.String())
		}
		return event
	})

	ui.ActiveFlex = &ui.Flex
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
