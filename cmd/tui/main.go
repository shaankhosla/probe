package tui

import (
	_ "github.com/marcboeker/go-duckdb"

	"fmt"
	"github.com/gdamore/tcell/v2"
	_ "github.com/marcboeker/go-duckdb"
	"github.com/rivo/tview"
)

const (
	maxColSize  = 50
	defaultRows = 100
)

type TUI struct {
	App          *tview.Application
	ResultsTable *tview.Table
	Columns      *tview.TextView
	Error        *tview.TextView
	InputField   *tview.TextArea
	InnerFlex    *tview.Flex
	OuterFlex    *tview.Flex
}

func ShowError(TUI TUI, err error) {
	TUI.ResultsTable.Clear()
	TUI.InnerFlex.ResizeItem(TUI.ResultsTable, 0, 0)
	TUI.InnerFlex.ResizeItem(TUI.Error, 0, 10)
	TUI.Error.SetText(fmt.Sprintf("Error: %s", err.Error())).ScrollToBeginning()
}
func UpdateTable(TUI TUI, rows [][]string) {
	TUI.InnerFlex.ResizeItem(TUI.Error, 0, 0)
	TUI.InnerFlex.ResizeItem(TUI.ResultsTable, 0, 10)
	TUI.Error.Clear()

	TUI.ResultsTable.Clear()

	for rowIndex, row := range rows {
		for colIndex, cell := range row {
			TUI.ResultsTable.SetCell(rowIndex, colIndex,
				tview.NewTableCell(cell).
					SetAlign(tview.AlignCenter).
					SetSelectable(rowIndex != 0).SetMaxWidth(maxColSize),
			)
		}
	}
	TUI.ResultsTable.ScrollToBeginning()
	TUI.ResultsTable.SetFixed(1, 0)
}

func CreateTUIAssets(filename string) *TUI {
	app := tview.NewApplication()

	resultsTable := tview.NewTable().
		SetBorders(true)

	resultsTable.SetTitle(filename)
	columnsTextView := tview.NewTextView().
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)

	columnsBoxWithContent := tview.NewFlex().
		AddItem(columnsTextView, 0, 1, false)

	errorTextView := tview.NewTextView().
		SetWrap(true).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	var inputField *tview.TextArea
	var innerFlex *tview.Flex

	inputField = tview.NewTextArea().SetLabel("Query: ").SetOffset(1, 1)
	inputField.SetText("select * from data", true)
	inputField.SetTitle("Query").SetBorder(true)

	resultsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		if key == tcell.KeyTab {
			if resultsTable.HasFocus() {
				app.SetFocus(inputField)
			}
			return nil
		}
		return event

	})
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			app.Stop()
			return nil
		}
		return event
	})
	innerFlex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(inputField, 0, 3, true).
		AddItem(errorTextView, 0, 0, false).
		AddItem(resultsTable, 0, 10, false)
	flex := tview.NewFlex().
		AddItem(columnsBoxWithContent, 0, 1, false).
		AddItem(innerFlex, 0, 5, true)

	app.SetRoot(flex, true).SetFocus(inputField)
	return &TUI{
		App:          app,
		ResultsTable: resultsTable,
		Columns:      columnsTextView,
		Error:        errorTextView,
		InputField:   inputField,
		InnerFlex:    innerFlex,
		OuterFlex:    flex,
	}
}
