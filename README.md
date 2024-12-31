# Probe: Interactive SQL Query Tool for File Analysis 

**Probe** is a lightweight, open-source CLI tool designed to make it simpler to investigate files. 

## Key Features:

- **Run SQL Queries**: Easily filter, aggregate, or analyze your data using SQL syntax.
- **Interactive Interface**: A clean TUI for entering SQL queries and viewing results.
- **Column Insights**: Quickly display all column names and structure of the dataset.

## Usage
Probe runs in your terminal in interactive mode. After installation, use the following command to start analyzing your data files:

```bash
probe <FILE_PATH>
```

### Example:
For example, given a CSV file named `example.csv`:
```csv
id,name,age
1,John,30
2,Jane,25
3,Sam,22
```

Run Probe and start querying:
```bash
probe example.csv
```

### Interactive CLI Commands
1. **Default SQL Query**: By default, Probe executes `SELECT * FROM data LIMIT 10` to display a preview of the data.
2. **Type Your Query**: Enter any valid SQL query in the input field (e.g., `SELECT name, age FROM data WHERE age > 23;`).
3. **Error Handling**: If a query is invalid, Probe will display an error message in red below the results pane.
4. **Keyboard Shortcuts**: Press `Ctrl+C` to exit.


## How It Works
Probe works by leveraging DuckDB's in-memory structured query capabilities. When you provide a file (such as CSV), Probe:
1. Maps the file to a virtual table named `data`.
2. Allows SQL queries to filter, group, and analyze.
3. Displays the query results in a rich, interactive table in the terminal.


