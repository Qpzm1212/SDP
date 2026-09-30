package photo

import "fmt"

// wraps the legacy RAW lib so it fits the ColorEngine interface
type RawEngineAdapter struct {
	old *LegacyRawProcessor
}

func NewRawEngineAdapter() *RawEngineAdapter {
	return &RawEngineAdapter{old: &LegacyRawProcessor{}}
}

func (a *RawEngineAdapter) ApplyCorrection(img []byte, k float64) ([]byte, error) {
	res, code := a.old.ProcessRAW(string(img), int(k*100))

	switch code {
	case 0:
		return []byte(res), nil
	case -1:
		return nil, fmt.Errorf("%w: legacy engine got empty data", ErrInvalidFormat)
	case -2:
		return nil, fmt.Errorf("%w: contrast too high", ErrProcessing)
	default:
		return nil, fmt.Errorf("%w: unknown legacy code %d", ErrProcessing, code)
	}
}
