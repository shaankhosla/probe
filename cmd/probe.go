package cmd

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	_ "github.com/marcboeker/go-duckdb"
	"github.com/rivo/tview"
	"probe/cmd/query"
	"probe/cmd/tui"
)

const ()

func initializeViews(table *tview.Table, columns *tview.TextView, filename string) {
	results, err := query.ExecuteSQL("select * from data", filename)
	if err == nil {
		tui.UpdateTable(table, results)
	}

	columnText := query.GetAllColumns(filename)
	columns.SetText(columnText).ScrollToBeginning().SetBorder(true).SetTitle("Columns")
}

func RunProbe(filename string) error {
	app := tview.NewApplication()

	resultsTable := tview.NewTable().
		SetBorders(true)

	resultsTable.SetTitle(filename)
	columnsTextView := tview.NewTextView().
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)

	initializeViews(resultsTable, columnsTextView, filename)
	columnsBoxWithContent := tview.NewFlex().
		AddItem(columnsTextView, 0, 1, false)

	errorTextView := tview.NewTextView().
		SetWrap(true).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	var inputField *tview.TextArea
	var inner_flex *tview.Flex

	inputField = tview.NewTextArea().SetLabel("Query: ").SetOffset(1, 1)
	inputField.SetText("select * from data", true)
	inputField.SetTitle("Query").SetBorder(true).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		key := event.Key()
		if key == tcell.KeyEnter {

			query_text := inputField.GetText()
			results, err := query.ExecuteSQL(query_text, filename)
			if err == nil {
				inner_flex.ResizeItem(errorTextView, 0, 0)
				inner_flex.ResizeItem(resultsTable, 0, 10)
				tui.UpdateTable(resultsTable, results)
				errorTextView.Clear()
			} else {
				resultsTable.Clear()
				inner_flex.ResizeItem(resultsTable, 0, 0)
				inner_flex.ResizeItem(errorTextView, 0, 10)
				errorTextView.SetText(fmt.Sprintf("Error: %s", err.Error())).ScrollToBeginning()
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
	inner_flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(inputField, 0, 3, true).
		AddItem(errorTextView, 0, 0, false).
		AddItem(resultsTable, 0, 10, false)
	flex := tview.NewFlex().
		AddItem(columnsBoxWithContent, 0, 1, false).
		AddItem(inner_flex, 0, 5, true)

	if err := app.SetRoot(flex, true).SetFocus(inputField).Run(); err != nil {
		panic(err)
	}
	return nil
}
