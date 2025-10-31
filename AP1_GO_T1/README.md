# Project 01 — Go_Bootcamp

## Project «Smart Utilities»

### Implementation 1. Console Calculator

Implemented an interactive console calculator supporting basic arithmetic operations:

- Supported operations: addition, subtraction, multiplication, division
- Operands processed as floating-point numbers (float64)
- Division precision — 3 decimal places
- Input validation with handling of invalid values
- On input errors, program displays "Invalid input" and prompts for data again

**Example usage:**

```
Input left operand: 
> 10.5
Input operation: 
> *
Input right operand: 
> 2
Result: 21.000
```

### Implementation 2. Word Frequency Analyzer

Program for text analysis and identification of most frequently occurring words:

- Accepts input string and number K
- Counts frequency of each word in the string
- Returns K most frequent words sorted by descending frequency
- Words with equal frequency are sorted lexicographically
- Handles edge cases: empty strings, K larger than unique word count
- Word separator — space character

**Example usage:**

```
Input: "aa bb cc aa cc cc cc aa ab ac bb"
K: 3
Result: cc aa bb
```

### Implementation 3. Collection Intersection Finder

Utility for finding intersection of two number sets:

- Reads two lists of integers separated by spaces
- Finds common elements in both lists
- Preserves element order from the first list
- Handles input errors (non-numeric values)
- Returns "Empty intersection" when no common elements exist

**Example usage:**

```
First list: 5 3 4 2 1 6
Second list: 6 4 2 4
Result: 4 2 6
```

### Implementation 4. Visit Tracking System

Interactive system for maintaining medical visit records:

**Supported commands:**

- **Save** — stores visit information (patient name, doctor specialization, date)
- **GetHistory** — retrieves complete visit history for a patient
- **GetLastVisit** — gets date of last visit to specific specialist

**Implementation features:**

- In-memory data storage using map structures
- Date format: YYYY-MM-DD
- Error handling (PatientNotFoundError for missing patients)
- Interactive operation without program restart

**Example usage:**

```
Save
Ivanov Ivan Ivanovich
orthopedist
2024-04-13

GetHistory
Ivanov Ivan Ivanovich
Result: orthopedist 2024-04-13
```

All implementations use only Go standard library.