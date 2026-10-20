package main

import "testing"

// Fake sensor for tests
type FakeSensor struct {
	Lux int
}

func (f FakeSensor) GetLux() int {
	return f.Lux
}

func TestToneStrategies(t *testing.T) {
	if got := (BrightTone{}).ApplyTone(100); got != 140 {
		t.Errorf("want 140, got %d", got)
	}

	if got := (ContrastTone{}).ApplyTone(100); got != 70 {
		t.Errorf("want 70, got %d", got)
	}

	sensor := FakeSensor{Lux: 20} // < 50, so pixel * 2
	if got := (SensorTone{Sensor: sensor}).ApplyTone(100); got != 200 {
		t.Errorf("want 200, got %d", got)
	}
}

func TestResizeStrategies(t *testing.T) {
	if got := (HalfScale{}).CalcSize(200, 100); got != 5000 {
		t.Errorf("want 5000, got %d", got)
	}

	if got := (SquareCrop{}).CalcSize(200, 100); got != 10000 {
		t.Errorf("want 10000, got %d", got)
	}

	if got := (Thumbnail{}).CalcSize(200, 100); got != 5000 {
		t.Errorf("want 5000, got %d", got)
	}
}

func TestSelector(t *testing.T) {
	sel := Selector{}
	sensor := FakeSensor{Lux: 20}

	if _, ok := sel.PickTone("contrast", sensor).(ContrastTone); !ok {
		t.Error("expected ContrastTone")
	}
	if _, ok := sel.PickTone("auto", sensor).(SensorTone); !ok {
		t.Error("expected SensorTone")
	}
	if _, ok := sel.PickTone("other", sensor).(BrightTone); !ok {
		t.Error("expected BrightTone")
	}

	if _, ok := sel.PickResize("avatar").(SquareCrop); !ok {
		t.Error("expected SquareCrop")
	}
	if _, ok := sel.PickResize("preview").(Thumbnail); !ok {
		t.Error("expected Thumbnail")
	}
	if _, ok := sel.PickResize("default").(HalfScale); !ok {
		t.Error("expected HalfScale")
	}
}

// Parameterized test: checks that output pixel is never negative and <= 255
func TestToneBoundsParameterized(t *testing.T) {
	sensor := FakeSensor{Lux: 10}
	strategies := []struct {
		name string
		s    ToneStrategy
	}{
		{"Bright", BrightTone{}},
		{"Contrast", ContrastTone{}},
		{"Sensor", SensorTone{Sensor: sensor}},
	}

	for _, tc := range strategies {
		t.Run(tc.name, func(t *testing.T) {
			ctx := NewToneContext(tc.s)
			res := ctx.Execute(150)
			if res < 0 || res > 255 {
				t.Errorf("%s out of bounds: %d", tc.name, res)
			}
		})
	}
}
