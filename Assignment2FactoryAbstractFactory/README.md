# Builder Design Pattern - Custom Bicycle (Go Implementation)

## Overview
This project implements the **Builder Design Pattern** in Go (Golang). It demonstrates how to construct complex objects step by step using a fluent API.

## Project Structure
- **Bicycle**: The immutable product struct (`name`, `model`, `year`).
- **Builder**: The interface declaring the construction steps.
- **ObjectBuilder**: A concrete builder that constructs a `Bicycle` struct.
- **SpecBuilder**: A concrete builder that generates a formatted text specification using `strings.Builder`.
- **Director**: Manages the build process to create specific presets (e.g., Mountain Bike, Road Bike).

## How to Run
```bash
go run main.go