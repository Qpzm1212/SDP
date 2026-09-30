package photo

import "fmt"

// old RAW lib we can't touch, hence the adapter
type LegacyRawProcessor struct{}

// takes string instead of []byte, returns an int code instead of error:
// 0 = ok, -1 = empty data, -2 = contrast too high
func (l *LegacyRawProcessor) ProcessRAW(raw string, contrast int) (string, int) {
	if raw == "" {
		return "", -1
	}
	if contrast > 100 {
		return "", -2
	}
	return fmt.Sprintf("RAW_Processed(contrast:%d): %s", contrast, raw), 0
}
