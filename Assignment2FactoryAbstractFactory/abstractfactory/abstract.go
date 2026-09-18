package abstractfactory

import "fmt"

type Frame interface {
	GetGeometry() string
}

type Wheel interface {
	GetTireType() string
}

// +++++++++++++++++++++++++++++++++++

type BikeFactory interface {
	CreateFrames() Frame
	CreateWheels() Wheel
}

// Down
type DownhillFrame struct{}

func (f *DownhillFrame) GetGeometry() string {
	return "Downhill frame: Strong construction"
}

type DownhillWheels struct{}

func (f *DownhillWheels) GetTireType() string {
	return "Downhill wheels 29 inch with wiiiiide protectors"
}

type DownhillFactory struct{}

func (f *DownhillFactory) CreateFrames() Frame {
	return &DownhillFrame{}
}

func (f *DownhillFactory) CreateWheels() Wheel {
	return &DownhillWheels{}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type RoadbikeFrame struct{}

func (r *RoadbikeFrame) GetGeometry() string {
	return "Road Bike, lightweight and high profile "
}

type RoadbikeWheel struct{}

func (r *RoadbikeWheel) GetTireType() string {
	return "Road Bike wheels, lightweight without camera, alumni, small width"
}

type RoadbikeFactory struct{}

func (r *RoadbikeFactory) CreateFrames() Frame {
	return &RoadbikeFrame{}
}
func (f *RoadbikeFactory) CreateWheels() Wheel {
	return &RoadbikeWheel{}
}

func AssembleBike(factory BikeFactory) {
	frame := factory.CreateFrames()
	wheel := factory.CreateWheels()

	fmt.Println("Assembling///")
	fmt.Println(frame.GetGeometry())
	fmt.Println(wheel.GetTireType())
	fmt.Println("Bicycle is ready\n")
}
