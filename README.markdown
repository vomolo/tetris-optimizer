# Tetris Solver

## Overview
This Go program solves a Tetris puzzle by taking a text file containing tetromino definitions and arranging them on a square board of minimal size without overlap. Each tetromino is represented by a 4x4 grid of `#` (blocks) and `.` (empty spaces), and the program outputs a string representation of the board with tetrominoes labeled by unique letters (A, B, C, ...).

The solver uses two strategies:
- An optimized approach for cases where all tetrominoes are identical (at least 5 pieces).
- A general backtracking algorithm for arbitrary tetromino sets.

The program includes robust input validation, error handling, and a test suite to ensure correctness.

## Author
- **Name**: Valentine Odhiambo Omolo
- **GitHub**: [vomolo](https://github.com/vomolo)
- **Email**: valentineomolo1@gmail.com

## Features
- Validates input files for correct format and tetromino connectivity.
- Handles both repetitive and unique tetromino sets.
- Uses a backtracking algorithm to find a valid board configuration.
- Outputs a clean string representation of the solution.
- Includes a test suite for verifying functionality.

## Prerequisites
- Go 1.16 or later
- A text editor or IDE (e.g., VS Code, GoLand)
- Basic familiarity with command-line interfaces

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/vomolo/tetris-solver.git
   cd tetris-solver
   ```
2. Ensure Go is installed. Verify with:
   ```bash
   go version
   ```
3. No external dependencies are required, as the program uses only the Go standard library.

## Usage
1. Create an input file in the `testfiles` directory with tetromino definitions. Each tetromino is a 4x4 grid of `#` and `.`, separated by an empty line. Example (`testfiles/input.txt`):
   ```
   ##..
   ##..
   ....
   ....
   ```
2. Run the program, specifying the input file:
   ```bash
   go run main.go testfiles/input.txt
   ```
3. The program outputs the solution to stdout, e.g.:
   ```
   AA
   AA
   ```
   If an error occurs (e.g., invalid file or unsolvable puzzle), it prints `ERROR` to stderr and exits.

## File Structure
- `main.go`: Entry point, handles command-line arguments and initiates solving.
- `main_test.go`: Test suite for the main function.
- `board.go`: Defines the `Board` struct and methods for placing/removing tetrominoes.
- `errors.go`: Custom error type for validation errors.
- `solver.go`: Core solving logic, including optimized and general solvers.
- `tetromino.go`: Defines and validates tetromino structures.
- `validator.go`: Handles file reading and input validation.
- `testfiles/`: Directory for input files (created automatically during tests).

## Input File Format
- The file must have a `.txt` extension and reside in the `testfiles` directory.
- Each tetromino is defined in a 4x4 grid using `#` for blocks and `.` for empty spaces.
- Tetrominoes are separated by a single empty line.
- Each tetromino must have exactly 4 `#` characters, forming a connected shape.
- Example:
  ```
  ##..
  ##..
  ....
  ....
  ```

## Running Tests
The program includes a test suite in `main_test.go`. To run the tests:
```bash
go test -v
```
The tests verify:
- Correct output for a valid input file.
- Proper error handling for missing command-line arguments.

## How It Works
1. **Input Validation**: The program checks the file path, extension, and content format.
2. **Tetromino Creation**: Parses the input file into tetrominoes, ensuring each has 4 connected blocks.
3. **Solving**:
   - If all tetrominoes are identical (and there are at least 5), an optimized grid-based placement is attempted.
   - Otherwise, a backtracking algorithm tries all possible placements on increasing board sizes.
4. **Output**: The solution is a string where each tetromino is represented by a unique letter (A, B, C, ...), with `.` for empty spaces.

## Limitations
- The program assumes tetrominoes are provided in a fixed orientation (no rotations or reflections).
- The solver may be slow for large numbers of tetrominoes due to the exponential nature of backtracking.
- The optimized solver only applies to identical tetrominoes with at least 5 pieces.

## Contributing
Contributions are welcome! Please:
1. Fork the repository.
2. Create a feature branch (`git checkout -b feature-name`).
3. Commit your changes (`git commit -m "Add feature"`).
4. Push to the branch (`git push origin feature-name`).
5. Open a pull request.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.