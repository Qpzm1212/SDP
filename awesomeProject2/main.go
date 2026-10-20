package main

import (
	"fmt"
	"os"
)

type Hero interface {
	Attack() string
}

type Villain interface {
	Scheme() string
}

type Vehicle interface {
	Drive() string
}

type UniverseFactory interface {
	CreateHero() Hero
	CreateVillain() Villain
	CreateVehicle() Vehicle
}

type SteampunkHero struct{}

func (h *SteampunkHero) Attack() string { return "Hero fires steam revolver!" }

type SteampunkVillain struct{}

func (v *SteampunkVillain) Scheme() string { return "Villain builds giant clockwork automaton!" }

type SteampunkVehicle struct{}

func (v *SteampunkVehicle) Drive() string { return "Vehicle launches airship!" }

type SteampunkUniverseFactory struct{}

func (f *SteampunkUniverseFactory) CreateHero() Hero       { return &SteampunkHero{} }
func (f *SteampunkUniverseFactory) CreateVillain() Villain { return &SteampunkVillain{} }
func (f *SteampunkUniverseFactory) CreateVehicle() Vehicle { return &SteampunkVehicle{} }

type CyberpunkHero struct{}

func (h *CyberpunkHero) Attack() string { return "Hero hacks neural link!" }

type CyberpunkVillain struct{}

func (v *CyberpunkVillain) Scheme() string { return "Villain releases rogue AI!" }

type CyberpunkVehicle struct{}

func (v *CyberpunkVehicle) Drive() string { return "Vehicle speeds on hoverbike!" }

type CyberpunkUniverseFactory struct{}

func (f *CyberpunkUniverseFactory) CreateHero() Hero       { return &CyberpunkHero{} }
func (f *CyberpunkUniverseFactory) CreateVillain() Villain { return &CyberpunkVillain{} }
func (f *CyberpunkUniverseFactory) CreateVehicle() Vehicle { return &CyberpunkVehicle{} }

type MedievalHero struct{}

func (h *MedievalHero) Attack() string { return "Hero swings broadsword!" }

type MedievalVillain struct{}

func (v *MedievalVillain) Scheme() string { return "Villain summons dark curse!" }

type MedievalVehicle struct{}

func (v *MedievalVehicle) Drive() string { return "Vehicle rides warhorse!" }

type MedievalUniverseFactory struct{}

func (f *MedievalUniverseFactory) CreateHero() Hero       { return &MedievalHero{} }
func (f *MedievalUniverseFactory) CreateVillain() Villain { return &MedievalVillain{} }
func (f *MedievalUniverseFactory) CreateVehicle() Vehicle { return &MedievalVehicle{} }

type NoirHero struct{}

func (h *NoirHero) Attack() string { return "Hero shoots snub-nosed revolver!" }

type NoirVillain struct{}

func (v *NoirVillain) Scheme() string { return "Villain runs corrupt syndicate!" }

type NoirVehicle struct{}

func (v *NoirVehicle) Drive() string { return "Vehicle cruises in 1940s sedan!" }

type NoirUniverseFactory struct{}

func (f *NoirUniverseFactory) CreateHero() Hero       { return &NoirHero{} }
func (f *NoirUniverseFactory) CreateVillain() Villain { return &NoirVillain{} }
func (f *NoirUniverseFactory) CreateVehicle() Vehicle { return &NoirVehicle{} }

var registry = map[string]UniverseFactory{
	"steampunk": &SteampunkUniverseFactory{},
	"cyberpunk": &CyberpunkUniverseFactory{},
	"medieval":  &MedievalUniverseFactory{},
	"noir":      &NoirUniverseFactory{},
}

func RegisterUniverse(name string, factory UniverseFactory) {
	registry[name] = factory
}

type UniverseClient struct {
	factory UniverseFactory
}

func NewUniverseClient(factory UniverseFactory) *UniverseClient {
	return &UniverseClient{factory: factory}
}

func (c *UniverseClient) Simulate() {
	hero := c.factory.CreateHero()
	villain := c.factory.CreateVillain()
	vehicle := c.factory.CreateVehicle()

	fmt.Println("  [Hero]:   ", hero.Attack())
	fmt.Println("  [Villain]:", villain.Scheme())
	fmt.Println("  [Vehicle]:", vehicle.Drive())
}

func GetFactory(flag string) (UniverseFactory, error) {
	factory, exists := registry[flag]
	if !exists {
		return nil, fmt.Errorf("unknown universe: %s", flag)
	}
	return factory, nil
}

func main() {
	universes := []string{"steampunk", "cyberpunk", "medieval", "noir"}

	for _, flag := range universes {
		fmt.Printf("=== Universe: [%s] ===\n", flag)
		factory, err := GetFactory(flag)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		client := NewUniverseClient(factory)
		client.Simulate()
		fmt.Println()
	}
}
