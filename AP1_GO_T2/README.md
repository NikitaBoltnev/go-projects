# Project 02 — Go_Bootcamp

## Project “Concurrent Execution”

### Task 1. Stopwatch for Asynchronous Tasks

1. The program reads two launch arguments: N and M.
2. Parameters N and M are passed as arguments when launching the program.
3. The program launches N goroutines, each of which sleeps (`time.Sleep`) for a random duration of up to M milliseconds.
4. The program waits for all goroutines to finish.
5. The program prints a list to the console consisting of pairs `<goroutine number, sleep time>`, sorted in descending order of sleep time.
6. The goroutine number is the iteration index of the loop in which the goroutine was launched.
7. The sleep time is the number of milliseconds the goroutine slept.

**The use of channels is not allowed.**

**Hint:** Use the sync package. Wait for all goroutines to complete before starting the output. Use the flag package for parsing arguments.

### Task 2. Square Number Generator

1. The program reads two arguments from the command line: `K` and `N`.
2. The parameters `K` and `N` are passed via command-line arguments.
3. The program launches two functions: a generator and a squaring function.
4. The parameters `N` and `M` are of type `int`.
5. The generator function starts a goroutine and returns channel 1. Inside the goroutine, numbers from `K` to `N` (inclusive) are generated and sent into channel 1.
6. The squaring function starts a goroutine and returns channel 2. Inside the goroutine, numbers are read from channel 1, squared, and the result is sent to channel 2.
7. The main program (`main`) reads numbers from channel 2 and prints them to the console.
8. The squaring function must accept channel 1 as a read-only channel, which is returned by the generator function.
9. Both the squaring and generator functions must run concurrently.
10. Squaring must occur sequentially. After reading a number from channel 1, it must immediately be squared and sent to the next channel, and only then should the next number be processed.

**Hint:**  
Channels must be created inside the functions, and returned with read/write restrictions applied.  
Channels must be closed once the function finishes its work.  
The generator and squaring functions must operate concurrently.

### Task 3. Ticker

1. The program reads the parameter `K` from the command-line arguments.
2. The parameter `K` is passed through the arguments when launching the program.
3. `K` defines the ticker interval in seconds and must be of type `uint`.
4. The program prints to stdout the message Tick `<i>` since `<time>`, where `<i>` is the tick number and `<time>` is the time in seconds since the ticker started.
5. The program runs until the user sends a SIGTERM or SIGINT signal.
6. Upon receiving one of these signals, the program stops the ticker and prints the message Termination.
7. The ticker must operate asynchronously. It is forbidden to use functions from the `time` package such as `time.After` or `time.Ticker`.  
    You are allowed to use constants from the `time` package and the `Sleep` function.

### Bonus Task 4. LRU Cache Using Generics

1. You need to implement a package that contains a Cache structure using the **LRU** (Least Recently Used) algorithm.
2. A constructor function is required, which accepts the element type and cache capacity, and returns a pointer to the Cache structure.
3. The `Cache` should provide the following methods:
    - `Set` — add an item with a given key;
    - `Get` — retrieve an item by key;
    - `Clear` — delete all items in the cache.
4. The time complexity for `Set` and `Get` operations must be **O(1)**.
5. **LRU logic:**
    - When adding an item:
        - If the item is already in the cache, its position moves to the front.
        - If the item is not in the cache and capacity is not exceeded, it is added to the front.
        - If the item is not in the cache and capacity is exceeded, it is added to the front, and the last (least recently used) item is removed.
    - When retrieving an item:
        - If the item exists, it is moved to the front and returned along with an additional `true` flag.
        - If the item does not exist, return the type’s zero value and the `false` flag.
6. The cache must only work with the type specified when the constructor function is called (use **Golang generics**).
7. The cache must be **thread-safe**, meaning concurrent access from different goroutines must not lead to race conditions.
8. You must implement tests for the following scenarios:
    - Rarely used items are removed from the cache.
    - If capacity is exceeded, rarely used items should be evicted.