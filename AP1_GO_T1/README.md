# Project 01 — Go_Bootcamp

### Project «Smart Utilities»

### Task 1. Console Calculator

You need to implement a console calculator:

1. The user enters the left operand and presses Enter.
2. The user enters the operation and presses Enter.
3. The user enters the right operand and presses Enter.
4. The program performs the selected operation and outputs the result to the console.
5. The calculator must support addition, subtraction, multiplication, and division.
6. Division accuracy — 3 decimal places.
7. Before entering each operand or operation, the program should display a prompt to the user.
8. Operands and operation are entered via the console by pressing Enter after each input.
9. Operands must be of type Float64.
10. If invalid input is provided, the program should print an error message (Invalid input) to stdout and prompt the user to enter the data again.

**Functional requirements:**

- Example:

  Input left operand: 
  > `10`

  Input operation 
  > `+`

  Input right operand
  > `15`

Only the standard library may be used.

_Hint:_ Parameters are read from the console as strings and must be converted to the required type.

### Task 2. Most Frequent Words

You need to implement a program that receives a list of words and a number K as input. The output must be a list of the K most frequently occurring words, sorted in descending order of frequency.

1. The program reads a string of words.
2. The program reads the number K, which is the limit on the number of output words.
3. The program determines how many times each word appears in the input string.
4. The program sorts the resulting array in descending order and returns the first K words.
5. The delimiter between words is the space character.
6. If the number of unique words is less than K, the program returns a list of all unique words, sorted by frequency.
7. If an empty string is given as input, the output must also be an empty string.
8. A word is any sequence of characters separated by spaces.
9. If multiple words have the same frequency, they must be sorted lexicographically.
10. You must implement tests that cover the following cases:  
    a. Normal behavior when K is less than the number of unique words.  
    b. Input is an empty list of words.  
    c. Input list contains fewer unique words than K.

**Functional requirements:**

- Example:

Input:

```
aa bb cc aa cc cc cc aa ab ac bb
3
```

Output:

```
cc aa bb
```

Because:  
`cc` appears 4 times,  
`aa` — 3 times,  
`bb` — 2 times.

Only the standard library is allowed.  
We recommend using the testing package for writing your tests.

### Task 3. Slice Intersection

1. The program reads the first list of numbers from the console, separated by spaces; input is finalized by pressing Enter.
2. The program reads the second list of numbers from the console, separated by spaces; input is finalized by pressing Enter.
3. The input lists are unordered.
4. The program works only with numbers of type int.
5. The program finds the intersection of the two lists and returns it in the order the values appear in the first list.

**Input:**

```
5 3 4 2 1 6
6 4 2 4
```

**Output:**

```
4 2 6
```

1. If the intersection is empty, the program outputs: Empty intersection.
2. If the list contains values that are not of type int, the program outputs the error message Invalid input and terminates.

### Task 4. Visit Log

1. The program implements 3 commands in interactive mode (commands are entered in the console without stopping the main execution).
2. The **Save** operation allows saving the visitor’s full name, the doctor's specialization they visited, and the visit date into the log.  
    a. The doctor’s specialization is provided as a string.  

    b. The visit date is in the format YYYY-MM-DD. 

    c. `Example: Save /n Ivanov Ivan Ivanovich /n orthopedist /n 2024-04-13`

    d. `/n` means a line break, i.e. Enter.

1. The GetHistory operation allows viewing the patient's visit history. It takes the patient's full name and returns a list of specialization–date pairs.  
    a. Example: `GetHistory /n Ivanov Ivan Ivanovich`

    b. Output: `orthopedist 2024-04-13 /n neurologist 2024-05-24 /n`

    c. `/n` means a line break, i.e. Enter.

1. The GetLastVisit operation returns the last visit to a specific specialist. It takes the patient’s full name and the doctor’s specialization and returns the last visit date.  
    a. Example: `GetLastVisit /n Ivanov Ivan Ivanovich /n orthopedist`

    b. Output: `2024-04-13`
    
    c. `/n` means a line break, i.e. Enter.

1. If the patient is not found, the program returns an error of type `PatientNotFoundError` with the message `patient not found` (this must be implemented manually).

Data is stored in memory only; long-term data persistence is not required.

**Hint:** you may use the map type as a storage structure.  
To store visit data, define an additional structure with the fields: specialization and visit date.