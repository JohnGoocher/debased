package cmd

func pos(s []string, value string) int {
	for i, v := range s {
		if value == v {
			return i
		}
	}
	return -1
}

func contains(s []string, value string) bool {
	for _, v := range s {
		if value == v {
			return true
		}
	}
	return false
}

func hasValidAddDataKeywords(s []string) bool {
	// "addData INTO <table_name> COLUMNS <column_name(s)>... VALUES <value(s)>... PAY <max_payment_allowed>",
	if !contains(s, "INTO") {
		return false
	}
	if !contains(s, "COLUMNS") {
		return false
	}
	if !contains(s, "VALUES") {
		return false
	}
	if !contains(s, "PAY") {
		return false
	}
	return true
}

func hasValidReadDataKeywords(s []string) bool {
	// "readData COLUMNS <column_name(s)>... FROM <table_name> [WHERE] [<condition>] {[AND|OR] [<condition>]}...",
	if !contains(s, "COLUMNS") {
		return false
	}
	if !contains(s, "FROM") {
		return false
	}
	return true
}

func hasValidCreateTableKeywords(s []string) bool {
	// "createTable <table name> COLUMNS {<column_name> <data_type>}... PAY <max payment allowed>",
	if !contains(s, "COLUMNS") {
		return false
	}
	if !contains(s, "PAY") {
		return false
	}
	return true
}

func areColumnsValuesSameSize(s []string) bool {
	columnSize := pos(s, "VALUES") - pos(s, "COLUMNS") + 1
	valueSize := len(s) - 2 - pos(s, "VALUES") + 1
	return columnSize == valueSize
}
