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

func hasValidKeywords(s []string) bool {
	if !contains(s, "INTO") {
		return false
	}
	if !contains(s, "COLUMNS") {
		return false
	}
	if !contains(s, "VALUES") {
		return false
	}
	return true
}

func areColumnsValuesSameSize(s []string) bool {
	columnSize := pos(s, "VALUES") - pos(s, "COLUMNS") + 1
	valueSize := len(s) - 2 - pos(s, "VALUES") + 1
	return columnSize == valueSize
}
