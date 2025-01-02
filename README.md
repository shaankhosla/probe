# Probe: Interactive SQL Query Tool for File Analysis 

**Probe** is a lightweight, open-source CLI tool designed to make it simpler to investigate files. 


https://github.com/user-attachments/assets/02fdcda3-85e3-4878-8efc-cdd5032a06e3


## Key Features:

- **Run SQL Queries**: Easily filter, aggregate, or analyze your data using SQL syntax.
- **Interactive Interface**: A clean TUI for entering SQL queries and viewing results.
- **Column Insights**: Quickly display all column names and structure of the dataset.
- **Flexible filetypes**: Works out of the box on many different filetypes that [DuckDB supports]([url](https://duckdb.org/docs/data/data_sources.html)).

## Usage
Probe runs in your terminal in interactive mode. After installation, use the following command to start analyzing your data files:

```bash
probe <FILE_PATH>
```

### Example 1:
Given a CSV file named `example.csv`:
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

### Example 2:
Given a group of Parquet files in a directory with the same schema:

Run Probe and start querying:
```bash
probe *.parquet
```

### Interactive CLI Commands
1. **Default SQL Query**: By default, Probe executes `SELECT * FROM data LIMIT 10` to display a preview of the data.
2. **Type Your Query**: Enter any valid SQL query in the input field (e.g., `SELECT name, age FROM data WHERE age > 23;`).
3. **Error Handling**: If a query is invalid, Probe will display an error message in red below the results pane.
4. **Keyboard Shortcuts**: Press `Ctrl+C` to exit.


## Installation:
Follow one of the methods below to install `probe` on your system:


### **1. Automated Installation (Recommended)**

We provide an automated installation script for convenience. This script detects your system's operating system and architecture, downloads the appropriate binary from the latest release, and installs it for you.

To install `probe`, simply run:

```bash
curl -sL https://raw.githubusercontent.com/shaankhosla/probe/main/install.sh | bash
```

Alternatively, download the installation script manually and then execute it:

```bash
wget https://raw.githubusercontent.com/shaankhosla/probe/main/install.sh
chmod +x install.sh
./install.sh
```

Once installed, you can verify the installation by running:

```bash
probe --help
```

This will display the help information for `probe`.


### **2. Download Prebuilt Binaries (Manually)**

You can manually download a prebuilt binary from the **[Releases](https://github.com/shaankhosla/probe/releases)** page.

1. Navigate to the latest release.
2. Expand the "Assets" section.
3. Download the binary that matches your system:
   - For Linux (`amd64` or `arm64`), download `probe-linux-<architecture>`.
   - For macOS (`amd64` or `arm64`), download `probe-darwin-<architecture>`.
4. Move the binary to a directory in your system's `PATH` (such as `/usr/local/bin`) and make it executable:
   ```bash
   mv /path/to/downloaded/probe /usr/local/bin/probe
   chmod +x /usr/local/bin/probe
   ```


### **3. Build from Source**

If you'd prefer to build `probe` from source, follow these steps:

#### Prerequisites
- A working Go environment is required. You can install Go by following the instructions [here](https://golang.org/doc/install).

#### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/shaankhosla/probe.git
   cd probe
   ```

2. Build the binary:
   ```bash
   go build -o probe main.go
   ```

3. Move the binary to a directory in your system's `PATH` (such as `/usr/local/bin`) and make it executable:
   ```bash
   mv ./probe /usr/local/bin/probe
   chmod +x /usr/local/bin/probe
   ```

4. Verify the installation:
   ```bash
   probe --help
   ```

## How It Works
Probe works by leveraging DuckDB's in-memory structured query capabilities. When you provide a file (such as CSV), Probe:
1. Maps the file to a virtual table named `data`.
2. Allows SQL queries to filter, group, and analyze.
3. Displays the query results in a rich, interactive table in the terminal.


