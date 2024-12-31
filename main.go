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

func executeSQL(query string, filename string) [][]string {
	fromStatement := fmt.Sprintf("from '%s'", filename)
	if !strings.Contains(query, "limit") {
		fromStatement += " limit 10"
	}
	query = strings.Replace(query, "from data", fromStatement, 1)
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

	var result [][]string
	result = append(result, columns)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePointers := make([]interface{}, len(columns))
		for i := range values {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			log.Fatalf("Failed to scan row: %s", err)
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
		log.Fatalf("Row iteration error: %s", err)
	}

	// for _, row := range result {
	// 	fmt.Println(row)
	// }

	return result
}

func runProbe(filename string) error {
	app := tview.NewApplication()

	instructions := tview.NewTextView().
		SetText("Type SQL Query below and hit Enter to execute. Press Ctrl+C to exit.").
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)

	resultsTable := tview.NewTable().
		SetBorders(true)

	var inputField *tview.InputField

	inputField = tview.NewInputField().
		SetLabel("SQL Query: ").
		SetFieldWidth(40).
		SetFieldTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetDoneFunc(func(key tcell.Key) {
			query := inputField.GetText()
			results := executeSQL(query, filename)
			updateTable(resultsTable, results)
		})
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC { // Ctrl+C to quit
			app.Stop()
			return nil
		}
		return event
	})
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(instructions, 0, 1, false).
		AddItem(inputField, 3, 1, true).
		AddItem(resultsTable, 0, 10, false)
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
