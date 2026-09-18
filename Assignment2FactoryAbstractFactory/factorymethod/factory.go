package factorymethod

import "fmt"

type BikeComponent interface {
	Install() string
	Service() string
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type HydraulicBrakes struct {
	Brand string
	Fluid string
}

func (h *HydraulicBrakes) Install() string {
	return fmt.Sprintf("Installing the hydraulic brakes %s.", h.Brand)
}
func (h *HydraulicBrakes) Service() string {
	return fmt.Sprintf("Changing the fluid %s.", h.Fluid)
}

type Suspension struct {
	ForkModel    string
	SuspensionMM string
}

func (s *Suspension) Install() string {
	return fmt.Sprintf("Installing the fork %s.", s.ForkModel)
}
func (s *Suspension) Service() string {
	return fmt.Sprintf("Cleaning and lubricating the seals  %s.", s.SuspensionMM)
}

type Drivetrain struct {
	Brand string
	Gears int
}

func (d *Drivetrain) Install() string {
	return fmt.Sprintf("Installing the drivetrain %s with %d gears.", d.Brand, d.Gears)
}
func (d *Drivetrain) Service() string {
	return fmt.Sprintf("Regulating the shifter  %s.", d.Brand)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++
// ===================================================

type ComponentCreator interface {
	CreateComponent() BikeComponent
}

func PerformMaintenance(creator ComponentCreator) {
	component := creator.CreateComponent()
	fmt.Println("Starting the work")
	fmt.Println(component.Install())
	fmt.Println(component.Service())
	fmt.Println("Finished the work")
}

type BrakeCreator struct{}

func (b *BrakeCreator) CreateComponent() BikeComponent {
	return &HydraulicBrakes{Brand: "SHIMANO MT-200", Fluid: "Mineral Oil"}
}

type SuspensionCreator struct{}

func (s *SuspensionCreator) CreateComponent() BikeComponent {
	return &Suspension{ForkModel: "Air", SuspensionMM: "180"}
}

type DrivetrainCreator struct{}

func (d *DrivetrainCreator) CreateComponent() BikeComponent {
	return &Drivetrain{Brand: "Shimano DEORE", Gears: 24}
}
