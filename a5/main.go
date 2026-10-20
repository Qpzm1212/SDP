package main

import "fmt"

// External sensor interface
type LightSensor interface {
GetLux() int
}

type SimulatedSensor struct {
Lux int
}

func (s SimulatedSensor) GetLux() int {
return s.Lux
}

// Helper to keep pixel in 0..255 range
func clamp(val int) int {
if val > 255 {
return 255
}
if val < 0 {
return 0
}
return val
}

// --- Family 1: Tone Strategies ---

type ToneStrategy interface {
ApplyTone(pixel int) int
}

// 1. Simple linear shift
type BrightTone struct{}

func (BrightTone) ApplyTone(pixel int) int {
return clamp(pixel + 40)
}

// 2. Threshold logic: pushes darks down and lights up
type ContrastTone struct{}

func (ContrastTone) ApplyTone(pixel int) int {
if pixel > 128 {
return clamp(pixel + 30)
}
return clamp(pixel - 30)
}

// 3. Uses external sensor data
type SensorTone struct {
Sensor LightSensor
}

func (s SensorTone) ApplyTone(pixel int) int {
if s.Sensor.GetLux() < 50 {
return clamp(pixel * 2) // dark room -> double brightness
}
return clamp(pixel + 10)
}

// Context 1 (no if/switch inside)
type ToneContext struct {
strategy ToneStrategy
}

func NewToneContext(s ToneStrategy) *ToneContext {
return &ToneContext{strategy: s}
}

func (c *ToneContext) SetStrategy(s ToneStrategy) {
c.strategy = s
}

func (c *ToneContext) Execute(pixel int) int {
return c.strategy.ApplyTone(pixel)
}

// --- Family 2: Resize Strategies ---

type ResizeStrategy interface {
CalcSize(width, height int) int
}

// 1. Keeps half of total area
type HalfScale struct{}

func (HalfScale) CalcSize(w, h int) int {
return (w / 2) * (h / 2)
}

// 2. Crops to a square using the smaller side
type SquareCrop struct{}

func (SquareCrop) CalcSize(w, h int) int {
if w < h {
return w * w
}
return h * h
}

// 3. Fixed width thumbnail (width = 100)
type Thumbnail struct{}

func (Thumbnail) CalcSize(w, h int) int {
if w == 0 {
return 0
}
newH := (h * 100) / w
return 100 * newH
}

// Context 2 (no if/switch inside)
type ResizeContext struct {
strategy ResizeStrategy
}

func NewResizeContext(s ResizeStrategy) *ResizeContext {
return &ResizeContext{strategy: s}
}

func (c *ResizeContext) SetStrategy(s ResizeStrategy) {
c.strategy = s
}

func (c *ResizeContext) Execute(w, h int) int {
return c.strategy.CalcSize(w, h)
}

// --- Selector (only place with switch) ---

type Selector struct{}

func (Selector) PickTone(env string, sensor LightSensor) ToneStrategy {
switch env {
case "contrast":
return ContrastTone{}
case "auto":
return SensorTone{Sensor: sensor}
default:
return BrightTone{}
}
}

func (Selector) PickResize(target string) ResizeStrategy {
switch target {
case "avatar":
return SquareCrop{}
case "preview":
return Thumbnail{}
default:
return HalfScale{}
}
}

func main() {
sensor := SimulatedSensor{Lux: 30}
pixel := 100

fmt.Println("1. Same input on all Tone strategies:")
tones := []ToneStrategy{
BrightTone{},
ContrastTone{},
SensorTone{Sensor: sensor},
}

toneCtx := NewToneContext(BrightTone{})
for _, s := range tones {
toneCtx.SetStrategy(s)
fmt.Printf("%T -> %d\n", s, toneCtx.Execute(pixel))
}

fmt.Println("\n2. Manual choice vs Selector:")
sel := Selector{}

toneCtx.SetStrategy(BrightTone{})
fmt.Println("Manual (BrightTone):", toneCtx.Execute(pixel))

toneCtx.SetStrategy(sel.PickTone("auto", sensor))
fmt.Println("Selector ('auto'):", toneCtx.Execute(pixel))

fmt.Println("\n3. Switching Resize strategy at runtime:")
resizeCtx := NewResizeContext(HalfScale{})
fmt.Println("HalfScale (200x100):", resizeCtx.Execute(200, 100))

resizeCtx.SetStrategy(sel.PickResize("avatar"))
fmt.Println("Switched to avatar (SquareCrop):", resizeCtx.Execute(200, 100))

resizeCtx.SetStrategy(Thumbnail{})
fmt.Println("Switched to Thumbnail:", resizeCtx.Execute(200, 100))
}
