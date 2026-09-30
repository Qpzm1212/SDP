# Design Rationale: Batch Photo Processing and Export System

## 1. Selected Domain
For this project, a photo preparation and export system has been developed. The abstraction (`ExportTask`) represents the final export formats (for example, compressed export for Instagram or full-size export for high-quality printing). The implementer hierarchy (`ColorEngine`) is responsible for color correction algorithms (standard RGB correction, application of cinematic LUT filters, and processing of RAW files).

## 2. Required Complexity Module
The module **Dynamic implementor selection** was selected. The client is not tied to a specific color correction engine at the compilation stage. The code implements the `ResolveEngineFunc' function, which analyzes the input file extension during program execution (runtime) and dynamically substitutes the desired `ColorEngine' into the abstraction.

## 3. Incompatibility justification
The adaptable class is the `LegacyRawProcessor', a conditional closed library for working with RAW files. This class is **really incompatible** with the standard `ColorEngine` interface for the following reasons:
* **Data types and signature:** It accepts and returns data in the form of a `string`, whereas the system uses byte slices of `[]byte'. The intensity is passed as an integer `int' (0-100), not `float64'.
* **The mechanism of failures (Failure mechanism):** It does not return the standard Go `error` type but uses integer return codes (0 for success, -1 or -2 for specific errors). The adapter intercepts these specific codes and translates them into standard system error constants (failure translation)[cite: 1].

## 4. Why one pattern is not enough
*   **Why not just Bridge:** The Bridge pattern does an excellent job of preventing combinatorial explosion of subclasses. If we used only it, we could easily add new formats (WebExport) and new standard algorithms (BlackAndWhiteEngine). However, Bridge requires that all implementations strictly follow the `ColorEngine` interface. Without the Adapter, we wouldn’t be able to connect the `LegacyRawProcessor` without rewriting its source code (which is prohibited).
*   **Why not just Adapter:** If we used only the Adapter, we would solve the problem of RAW engine incompatibility. However, when adding new export options (Instagram, Print), we would have to create classes such as `RawInstagramExport`, `LutInstagramExport`, `RawPrintExport`, etc., mixing the export logic with the color logic. This would lead to a severe violation of the Open/Closed principle and a combinatorial explosion.

## 5. Limitation of the final design
The main limitation of the current architecture lies in the dynamic implementer selection module. At present, the routing function `ResolveEngineFunc` contains hardcoded conditional statements `if strings.HasSuffix(...)` for `.cr2` and `.png`. Adding a new color correction algorithm will require modifying this function, which constitutes a specific violation of the Open/Closed principle at the level of dependency injection. A potential solution in the future could be to use the Registry pattern for dynamically registering handlers.