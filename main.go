package main

import (
	"errors"
	"fmt"
	"strings"
)

type Bicycle struct {
	name  string
	model string
	year  int
}

func (b Bicycle) String() string {
	return fmt.Sprintf("Bicycle [Name=%s, Model=%s, Year=%d]", b.name, b.model, b.year)
}
func (b Bicycle) ToJSON() (string, error) {
	jsonData, err := json.MarshalIndent(struct {
		Name  string `json:"name"`
		Model string `json:"model"`
		Year  int    `json:"year"`
	}{
		Name:  b.name,
		Model: b.model,
		Year:  b.year,
	}, "", "  ")

	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

type Builder interface {
	SetName(name string) Builder
	SetModel(model string) Builder
	SetYear(year int) Builder
}

type ObjectBuilder struct {
	name  string
	model string
	year  int
}

func (b *ObjectBuilder) SetName(name string) Builder {
	b.name = name
	return b
}

func (b *ObjectBuilder) SetModel(model string) Builder {
	b.model = model
	return b
}

func (b *ObjectBuilder) SetYear(year int) Builder {
	b.year = year
	return b
}

func (b *ObjectBuilder) Build() (Bicycle, error) {
	if b.name == "" {
		return Bicycle{}, errors.New("name cannot be empty")
	}
	if b.model == "" {
		return Bicycle{}, errors.New("model cannot be empty")
	}
	if b.year <= 0 {
		return Bicycle{}, errors.New("year must be valid")
	}
	return Bicycle{
		name:  b.name,
		model: b.model,
		year:  b.year,
	}, nil
}

type SpecBuilder struct {
	sb strings.Builder
}

func (b *SpecBuilder) SetName(name string) Builder {
	b.sb.WriteString(fmt.Sprintf("Brand/Name: %s\n", name))
	return b
}

func (b *SpecBuilder) SetModel(model string) Builder {
	b.sb.WriteString(fmt.Sprintf("Model:      %s\n", model))
	return b
}

func (b *SpecBuilder) SetYear(year int) Builder {
	b.sb.WriteString(fmt.Sprintf("Year:       %d\n", year))
	return b
}

func (b *SpecBuilder) Build() (string, error) {
	return b.sb.String(), nil
}

type Director struct{}

func (d *Director) MakeMountainBike(b Builder) {
	b.SetName("Trinx").
		SetModel("M600").
		SetYear(2023)
}

func (d *Director) MakeRoadBike(b Builder) {
	b.SetName("Specialized").
		SetModel("Tarmac").
		SetYear(2024)
}

func main() {
	director := &Director{}

	objBuilder := &ObjectBuilder{}
	director.MakeRoadBike(objBuilder)

	roadBike, err := objBuilder.Build()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Representation 1 (Struct)")
		fmt.Println(roadBike)
	}

	fmt.Println()

	specBuilder := &SpecBuilder{}
	director.MakeMountainBike(specBuilder)

	specSheet, _ := specBuilder.Build()
	fmt.Println("Representation 2 (Spec)")
	fmt.Println(specSheet)
}
