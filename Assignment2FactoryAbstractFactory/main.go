package main

import (
	"awesomeProject/abstractfactory"
	"awesomeProject/factorymethod"
)

func main() {
	// A
	brakeFactory := &factorymethod.BrakeCreator{}
	suspensionFactory := &factorymethod.SuspensionCreator{}
	drivetrainFactory := &factorymethod.DrivetrainCreator{}

	factorymethod.PerformMaintenance(brakeFactory)
	factorymethod.PerformMaintenance(suspensionFactory)
	factorymethod.PerformMaintenance(drivetrainFactory)

	// B
	downhillFactory := &abstractfactory.DownhillFactory{}
	roadbikeFactory := &abstractfactory.RoadbikeFactory{}

	abstractfactory.AssembleBike(downhillFactory)
	abstractfactory.AssembleBike(roadbikeFactory)
}
