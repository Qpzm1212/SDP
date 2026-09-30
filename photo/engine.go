package photo

import (
	"errors"
)

// Errors
var (
	ErrInvalidFormat = errors.New("invalid format")
	ErrProcessing    = errors.New("processing failed")
)

// Implementor interf ++++++++++++
type ColorEngine interface {
	ApplyCorrection(imageData []byte, intensity float64) ([]byte, error)
}

// conc impl +++
type StandardRGBEngine struct{}

func (e *StandardRGBEngine) ApplyCorrection(imageData []byte, intensity float64) ([]byte, error) {
	if len(imageData) == 0 {
		return nil, ErrInvalidFormat
	}
	result := append([]byte("RGB_Processed"), imageData...)
	return result, nil
}

// conc impl2 ++++
type CinematicLUTEngine struct{}

func (e *CinematicLUTEngine) ApplyCorrection(imageData []byte, intensity float64) ([]byte, error) {
	if len(imageData) == 0 {
		return nil, ErrInvalidFormat
	}
	result := append([]byte("LUT_Processed"), imageData...)
	return result, nil
}
