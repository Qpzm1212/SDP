package photo

import (
	"errors"
	"strings"
	"testing"
)

// fake engine so tests don't depend on the real one
type stub struct{ fail bool }

func (s *stub) ApplyCorrection(img []byte, _ float64) ([]byte, error) {
	if s.fail {
		return nil, ErrProcessing
	}
	return append([]byte("STUB_"), img...), nil
}

func TestExportUsesEngine(t *testing.T) {
	old := ResolveEngineFunc
	t.Cleanup(func() { ResolveEngineFunc = old }) // restore the real resolver

	ResolveEngineFunc = func(string) ColorEngine { return &stub{} }

	got, err := (&InstagramExport{}).Export("test.jpg", []byte("data"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "STUB_data") {
		t.Errorf("stub wasn't called, got: %s", got)
	}
}

func checkErr(t *testing.T, img []byte, k float64, want error) {
	t.Helper()
	_, err := NewRawEngineAdapter().ApplyCorrection(img, k)
	if !errors.Is(err, want) {
		t.Errorf("want %v, got %v", want, err)
	}
}

func TestAdapterEmptyData(t *testing.T) {
	checkErr(t, []byte(""), 0.5, ErrInvalidFormat) // legacy code -1
}

func TestAdapterHighIntensity(t *testing.T) {
	checkErr(t, []byte("data"), 1.5, ErrProcessing) // legacy code -2, 1.5 = 150%
}
