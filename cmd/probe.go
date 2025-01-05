package cmd

import (
	"github.com/gdamore/tcell/v2"
	_ "github.com/marcboeker/go-duckdb"
	"probe/cmd/query"
	"probe/cmd/tui"
)

const ()

func initializeViews(TUI tui.TUI, filename string) {
	results, err := query.ExecuteSQL("select * from data", filename)
	if err == nil {
		tui.UpdateTable(TUI, results)
	}

	columnText := query.GetAllColumns(filename)
	TUI.Columns.SetText(columnText).ScrollToBeginning().SetBorder(true).SetTitle("Columns")
}

func RunProbe(filename string) error {
	TUI := tui.CreateTUIAssets(filename)
	initializeViews(*TUI, filename)
	TUI.InputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		key := event.Key()
		if key == tcell.KeyTab {
			if TUI.InputField.HasFocus() {
				TUI.App.SetFocus(TUI.ResultsTable)
			}
			return nil
		}
		if key == tcell.KeyEnter {

			query_text := TUI.InputField.GetText()
			results, err := query.ExecuteSQL(query_text, filename)
			if err == nil {
				tui.UpdateTable(*TUI, results)
			} else {
				tui.ShowError(*TUI, err)
			}

			return nil
		}
		return event

	})

	if err := TUI.App.Run(); err != nil {
		panic(err)
	}
	return nil
}
