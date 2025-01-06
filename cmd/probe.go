package cmd

import (
	"probe/cmd/query"
	"probe/cmd/tui"

	"github.com/gdamore/tcell/v2"
	_ "github.com/marcboeker/go-duckdb"
)

func initializeViews(tuiAssets tui.TUI, filename string) {
	results, err := query.ExecuteSQL("select * from data", filename)
	tuiAssets.UpdateTable(results, err)

	columnText := query.GetAllColumns(filename)
	tuiAssets.Columns.SetText(columnText).ScrollToBeginning().SetBorder(true).SetTitle("Columns")
}

func RunProbe(filename string) error {
	tuiAssets := tui.CreateTUIAssets(filename)
	query.InitializeDB()
	initializeViews(*tuiAssets, filename)
	tuiAssets.InputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		key := event.Key()
		if key == tcell.KeyTab {
			if tuiAssets.InputField.HasFocus() {
				tuiAssets.App.SetFocus(tuiAssets.ResultsTable)
			}
			return nil
		}
		if key == tcell.KeyEnter {

			query_text := tuiAssets.InputField.GetText()
			results, err := query.ExecuteSQL(query_text, filename)
			tuiAssets.UpdateTable(results, err)

			return nil
		}
		return event

	})

	if err := tuiAssets.App.Run(); err != nil {
		panic(err)
	}
	return nil
}
