package pagesize

import (
	"fmt"
	"strconv"
)

func ParsePageSize(raw string) (int, error) {
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("page size must be an integer: %w", err)
	}
	if n < 1 {
		return 0, fmt.Errorf("page size must be positive")
	}
	return n, nil
}
