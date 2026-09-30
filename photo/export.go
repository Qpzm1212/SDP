package photo

import (
	"fmt"
	"strings"
)

// ExportTask - br abst
type ExportTask interface {
	Export(filename string, imageData []byte) (string, error)
}

// ResolveEngineFunc - динамический выбор реализатора
var ResolveEngineFunc = func(filename string) ColorEngine {
	lowerName := strings.ToLower(filename)

	if strings.HasSuffix(lowerName, ".cr2") {
		return NewRawEngineAdapter()
	}

	if strings.HasSuffix(lowerName, ".png") {
		return &CinematicLUTEngine{}
	}

	return &StandardRGBEngine{}
}

// InstagramExport - ref abst 1
type InstagramExport struct{}

func (i *InstagramExport) Export(filename string, imageData []byte) (string, error) {
	engine := ResolveEngineFunc(filename)

	// Динамический вызов реализатора
	processedData, err := engine.ApplyCorrection(imageData, 0.8)
	if err != nil {
		return "", err
	}

	// Форматируем строку с обработанными данными, а не с исходными
	return fmt.Sprintf("Instagram Export (%s)", string(processedData)), nil
}

// PrintExport - ref abst 2
type PrintExport struct{}

func (p *PrintExport) Export(filename string, imageData []byte) (string, error) {
	engine := ResolveEngineFunc(filename)

	processedData, err := engine.ApplyCorrection(imageData, 0.1)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Print Export (%s)", string(processedData)), nil
}
