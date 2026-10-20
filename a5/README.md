# Assignment 5 - Strategy Pattern

Domain: **Image Processing**

## 1. Class Responsibilities (SRP)

| Class / Interface | Responsibility |
| :--- | :--- |
| `LightSensor` | Defines an interface for getting light levels from an external source. |
| `SimulatedSensor` | Provides simulated light values during normal program runs. |
| `FakeSensor` | Returns fixed light values to keep unit tests isolated. |
| `ToneStrategy` | Declares a common interface for pixel brightness algorithms. |
| `BrightTone` | Increases pixel brightness by adding a constant offset value. |
| `ContrastTone` | Pushes pixel brightness away from the middle threshold value. |
| `SensorTone` | Adjusts pixel brightness based on external light sensor readings. |
| `ToneContext` | Executes the active tone strategy without knowing its concrete type. |
| `ResizeStrategy` | Declares a common interface for image size calculation algorithms. |
| `HalfScale` | Calculates image area after reducing both dimensions by half. |
| `SquareCrop` | Computes the area of a square crop using the shortest side. |
| `Thumbnail` | Calculates image area scaled proportionally to a fixed width. |
| `ResizeContext` | Delegates size calculation to the currently set resize strategy. |
| `Selector` | Chooses the matching strategy instance based on current input conditions. |

## 2. Why Strategy Instead of State or Template Method?

This code uses the **Strategy** pattern because `ToneStrategy` and `ResizeStrategy` represent interchangeable algorithms that solve the same problem in different ways. The client (`main` or `Selector`) picks which algorithm to pass into `ToneContext` or `ResizeContext` at runtime. Strategies do not know about each other, and we can add new algorithms without modifying the contexts (Open/Closed Principle).

It is **not the State pattern** because the contexts do not go through internal lifecycle transitions (like an order changing from `Created` to `Paid`). In the State pattern, states usually switch the context to the next state themselves. Here, `BrightTone` never switches `ToneContext` into `ContrastTone` — switching is controlled from the outside.

It is **not the Template Method pattern** because there is no shared base class with a fixed step-by-step skeleton algorithm. Each strategy calculates the result in one independent step using its own rule, so interface composition is much cleaner than inheritance.

## 3. Refactoring Commits
* **Before refactoring:** [Commit 1 link]
* **After refactoring:** [Commit 2 link]
