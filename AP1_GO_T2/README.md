# Project 02 — Go_Bootcamp

## Project "Concurrent Execution"

### Implementation 1. Stopwatch for Asynchronous Tasks

Implemented a program that measures execution time of concurrent goroutines:

- Reads parameters N and M from command line arguments
- Launches N goroutines, each sleeping for random time up to M milliseconds
- Uses sync.WaitGroup to wait for all goroutines completion
- Collects and sorts results by sleep time in descending order
- Outputs pairs: <goroutine number, sleep time>

**Key features:**
- No channels used (as required)
- Proper synchronization with WaitGroup
- Random sleep duration generation
- Result sorting algorithm

### Implementation 2. Square Number Generator

Built a concurrent pipeline for number processing:

- Generator produces numbers from K to N and sends to channel 1
- Squaring function reads from channel 1, squares numbers, sends to channel 2
- Main program outputs results from channel 2
- Both processes run concurrently

**Key features:**
- Proper channel closing after completion
- Read/write restrictions on channels
- Sequential number processing
- Concurrent execution of generator and processor

### Implementation 3. Ticker with Signal Handling

Created asynchronous ticker with graceful shutdown support:

- Parameter K sets ticker interval in seconds
- Outputs "Tick <i> since <time>" messages with tick number and elapsed time
- Handles SIGTERM and SIGINT signals for graceful termination
- Outputs "Termination" message on shutdown

**Key features:**
- Uses signal.Notify for signal handling
- Asynchronous ticker operation
- Proper shutdown on OS signals
- No use of time.After or time.Ticker

### Implementation 4. LRU Cache Using Generics (Bonus)

Implemented thread-safe LRU cache with generics:

- Constructor accepts element type and cache capacity
- Methods Set, Get, Clear with O(1) complexity
- LRU algorithm with eviction of least recently used items
- Thread-safe for concurrent access from multiple goroutines

**Key features:**
- Generics for type safety
- Double-linked list + hash map for O(1) efficiency
- Mutex for thread safety
- Comprehensive test coverage for all scenarios

All implementations use only Go standard library.