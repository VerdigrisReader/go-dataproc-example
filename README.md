# Go Data Processing example

At my current company, we started to use the below coding homework exercise question to get an idea of candidate's coding ability. I decided to tackle this question in Go, a language I'm relatively new to, because it seemed unfair to assess candidates on a question that I hadn't completed myself. Using Go also meant that I couldn't lean on my data processing experience in python, and had to rethink many things that I had taken for granted before in a less strict and more feature-y language.

## Coding Question:
Write a python command line utility to read data from standard input and print the 'Average Order Value' for each day between a given start & end date.

Imagine that your expected usage is the unix command (or windows equivalent):

`cat orders.csv | python my-script.py --start=2015-02-03 --end=2015-03-23`

## Usage
Test:
```bash
go test
```

Build:
```bash
go build -o aov aov.go
```

Run:
```bash
cat orders.csv | ./aov
```

## Comments

* Doing financial calculations with float64 is a big red flag here, a future improvement should involve using `math/big . Rat()` for reliable Decimal handling
* Golang's very-clever non-standard date formatting rule is meant to make it easier by using sequential numbers insted of `%Y-%m-%d`, but isn't so clear for European date formats: `2006-01-02 15:04:05` is the formatter for the MySQL date format, not so memorable.
* Types everywhere. An Orders struct isn't so controversial, but defining a struct for a 'Thread', which contains an input & output channel might be overkill/a risky design pattern for larger scripts. I've gone for readability in this case: `thread.input` & `thread.output` are far more readable than `thread[0]`.
* Channels everywhere. A lot of functions/goroutines take 2 channels as input, which feels very reliable in some ways, (you can guarantee that the channel recieves every single piece of input, and only outputs once the input is closed). But as a Go newbie, i was wary of overuse of this pattern. It probably gets out of hand quite quickly when each function accepts a very particular set of channels to recieve a specific type. Maybe these configurations of channel input/outputs should be implemented as interfaces.