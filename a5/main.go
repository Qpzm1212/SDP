package main

import "fmt"

type LightSensor interface {
GetLux() int
}

type SimulatedSensor struct {
Lux int
}

func (s SimulatedSensor) GetLux() int {
return s.Lux
}

func clamp(val int) int {
if val > 255 {
return 255
}
if val < 0 {
return 0
}
return val
}

type ToneStrategy interface {
ApplyTone(pixel int) int
}

type BrightTone struct{}
func (BrightTone) ApplyTone(pixel int) int { return clamp(pixel + 40) }

type ContrastTone struct{}
func (ContrastTone) ApplyTone(pixel int) int {
if pixel > 128 { return clamp(pixel + 30) }
return clamp(pixel - 30)
}

type SensorTone struct{ Sensor LightSensor }
func (s SensorTone) ApplyTone(pixel int) int {
if s.Sensor.GetLux() < 50 { return clamp(pixel * 2) }
return clamp(pixel + 10)
}

// Bad way: ToneContext checks concrete types with switch instead of polymorphism
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
switch s := c.strategy.(type) {
case BrightTone:
return clamp(pixel + 40)
case ContrastTone:
if pixel > 128 {
return clamp(pixel + 30)
}
return clamp(pixel - 30)
case SensorTone:
if s.Sensor.GetLux() < 50 {
return clamp(pixel * 2)
}
return clamp(pixel + 10)
default:
return pixel
}
}

type ResizeStrategy interface {
CalcSize(width, height int) int
}

type HalfScale struct{}
func (HalfScale) CalcSize(w, h int) int { return (w / 2) * (h / 2) }

type SquareCrop struct{}
func (SquareCrop) CalcSize(w, h int) int {
if w < h { return w * w }
return h * h
}

type Thumbnail struct{}
func (Thumbnail) CalcSize(w, h int) int {
if w == 0 { return 0 }
newH := (h * 100) / w
return 100 * newH
}

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
fmt.Println("Legacy version before refactoring ToneContext")
}
