package pagesize

import "testing"

func TestParsePageSize(t *testing.T) {
	for _, raw := range []string{"1", "20", "100"} {
		if _, err := ParsePageSize(raw); err != nil {
			t.Errorf("ParsePageSize(%q) error = %v", raw, err)
		}
	}
	for _, raw := range []string{"", "abc", "0", "-1"} {
		if _, err := ParsePageSize(raw); err == nil {
			t.Errorf("ParsePageSize(%q) expected error", raw)
		}
	}
}
