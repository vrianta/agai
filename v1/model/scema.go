package model

import (
	"regexp"
	"strconv"
	"strings"
)

// ParseSQLType parses a SQL type like "VARCHAR(20)" or "TEXT"
// Returns: base type (e.g. "VARCHAR"), length (e.g. 20), or 0 if no length
func (sc *schema) parseSQLType() (string, int) {
	sqlType := strings.ToUpper(strings.TrimSpace(sc.fieldType))
	re := regexp.MustCompile(`^([A-Z]+)\((\d+)\)$`)

	matches := re.FindStringSubmatch(sqlType)
	if len(matches) == 3 {
		length, _ := strconv.Atoi(matches[2])
		return matches[1], length
	}

	// No length specified (e.g. TEXT or just VARCHAR)
	return sqlType, 0
}
