package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	_ "github.com/marcboeker/go-duckdb"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

func reformatSQL(query string, filename string) string {
	fromStatement := fmt.Sprintf("from '%s'", filename)
	if !strings.Contains(query, "limit") {
		fromStatement += " limit 10"
	}
	query = strings.Replace(query, "from data", fromStatement, 1)
	return query

}

func getAllColumns(filename string) string {
	query := "select * from data"
	query = reformatSQL(query, filename)
	db, err := sql.Open("duckdb", "")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Failed to execute query: %s", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("Failed to fetch column names: %s", err)
	}

	columnText := ""
	for _, col := range columns {
		columnText += col + "\n"
	}
	return columnText

}
func initializeViews(table *tview.Table, columns *tview.TextView, filename string) {
	results, err := executeSQL("select * from data", filename)
	if err == nil {
		updateTable(table, results)
	}

	columnText := getAllColumns(filename)
	columns.SetText(columnText).ScrollToBeginning()
}

func executeSQL(query string, filename string) ([][]string, error) {
	query = reformatSQL(query, filename)
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DuckDB: %w", err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch column names: %w", err)
	}

	var result [][]string
	result = append(result, columns)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePointers := make([]interface{}, len(columns))
		for i := range values {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make([]string, len(columns))
		for i, v := range values {
			if v != nil {
				row[i] = fmt.Sprintf("%v", v)
			} else {
				row[i] = ""
			}
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return result, nil
}

func runProbe(filename string) error {
	app := tview.NewApplication()

	instructions := tview.NewTextView().
		SetText("Type SQL Query below and hit Enter to execute. Press Ctrl+C to exit.").
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)

	resultsTable := tview.NewTable().
		SetBorders(true)

	columnsTextView := tview.NewTextView().
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)

	initializeViews(resultsTable, columnsTextView, filename)
	columnsBoxWithContent := tview.NewFlex().
		AddItem(columnsTextView, 0, 1, false)

	var inputField *tview.InputField
	var errorTextView *tview.TextView
	errorTextView = tview.NewTextView().
		SetWrap(true).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	inputField = tview.NewInputField().
		SetLabel("SQL Query: ").
		SetText("select * from data").
		SetFieldWidth(40).
		SetFieldTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetDoneFunc(func(key tcell.Key) {
			query := inputField.GetText()
			results, err := executeSQL(query, filename)
			if err != nil {
				errorTextView.SetText(fmt.Sprintf("Error: %s", err.Error()))
				resultsTable.Clear()
			} else {
				updateTable(resultsTable, results)
				errorTextView.SetText("")
			}
		})
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			app.Stop()
			return nil
		}
		return event
	})
	flex := tview.NewFlex().
		AddItem(columnsBoxWithContent, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(instructions, 0, 1, false).
			AddItem(inputField, 3, 1, true).
			AddItem(errorTextView, 3, 1, false).
			AddItem(resultsTable, 0, 10, false), 0, 5, true)

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
	return nil
}

func updateTable(table *tview.Table, rows [][]string) {
	table.Clear()

	for rowIndex, row := range rows {
		for colIndex, cell := range row {
			table.SetCell(rowIndex, colIndex,
				tview.NewTableCell(cell).
					SetAlign(tview.AlignCenter).
					SetSelectable(rowIndex != 0),
			)
		}
	}
	table.ScrollToBeginning()
}

func main() {
	app := &cli.App{
		Name:  "probe",
		Usage: "Interactive SQL query tool for file analysis.",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("error: you must provide a filename as an argument")
			}
			filename := c.Args().Get(0)

			if _, err := os.Stat(filename); err != nil {
				return fmt.Errorf("file does not exist: %s", filename)
			}

			return runProbe(filename)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
