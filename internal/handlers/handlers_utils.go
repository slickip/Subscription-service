package handlers

import (
	"errors"
	"strconv"
	"strings"
)

func monthsBetween(
	startMonth, startYear int,
	endMonth, endYear int,
) int {
	return (endYear-startYear)*12 + (endMonth - startMonth) + 1
}

func maxYearMonth(y1, m1, y2, m2 int) (int, int) {
	if y1 > y2 || (y1 == y2 && m1 >= m2) {
		return y1, m1
	}
	return y2, m2
}

func minYearMonth(y1, m1, y2, m2 int) (int, int) {
	if y1 < y2 || (y1 == y2 && m1 <= m2) {
		return y1, m1
	}
	return y2, m2
}

func parseMonthYear(value string) (int, int, error) {
	parts := strings.Split(value, "-")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid date format")
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return 0, 0, errors.New("invalid month")
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, errors.New("invalid year")
	}

	return month, year, nil
}
