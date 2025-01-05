package tui

import (
	_ "github.com/marcboeker/go-duckdb"
	"github.com/rivo/tview"
)

const (
	maxColSize  = 50
	defaultRows = 100
)

func UpdateTable(table *tview.Table, rows [][]string) {
	table.Clear()

	for rowIndex, row := range rows {
		for colIndex, cell := range row {
			table.SetCell(rowIndex, colIndex,
				tview.NewTableCell(cell).
					SetAlign(tview.AlignCenter).
					SetSelectable(rowIndex != 0).SetMaxWidth(maxColSize),
			)
		}
	}
	table.ScrollToBeginning()
}
