package cmd

// // PayArg hey
// type PayArg struct {
// 	payment int
// }

// // TableNameArg hey
// type TableNameArg struct {
// 	name string
// }

// // ValueArgs hey
// type ValueArgs struct {
// 	values []string
// }

// // ColumnArgs hey
// type ColumnArgs struct {
// 	columnNames []string
// }

// AddDataRequiredArgs hey
type AddDataRequiredArgs struct {
	//tableName   *TableNameArg
	tableName   string
	columnNames []string
	valueNames  []string
	payment     int
}

// ReadDataArgs hey
type ReadDataArgs struct {
	// "readData COLUMNS <column_name(s)>... FROM <table_name> [WHERE] [<condition>] {[AND|OR] [<condition>]}...",
	columnNames []string
	tableName   string
	conditions  []string
}

// ConnectArgs hey
type ConnectArgs struct {
	// "connect <ip address>"
	ip string
}

// CreateTable hey
type CreateTable struct {
	// createTable <table name> {<column_name> <data_type>}... PAY <max payment allowed>
	tableName      string
	columnNames    []string
	columnDataType []string
	payment        int
}
