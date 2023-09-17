# Guide

This guide and codebase is a companion resource for those reading [On Computable Numbers, with an Application to the Entscheidungsproblem](https://www.cs.virginia.edu/~robins/Turing_Paper_1936.pdf) by Alan Turing.

I can't recommend Charles Petzold's [The Annotated Turing](https://www.amazon.com/Annotated-Turing-Through-Historic-Computability/dp/0470229055) enough, this project would not have been possible without it.

The guide annotates the paper section-by-section. I only quote directly from the paper when there is something to call out explicitly. When possible I bias towards a fully working implementation in code that the reader can use themselves.

## Introduction

### What are computable numbers?

Turing explains that his paper will be about [computable numbers](https://en.wikipedia.org/wiki/Computable_number). He briefly discusses [computable functions](https://en.wikipedia.org/wiki/Computable_function) saying they share the same fundamental problems, however he chooses computable numbers as the subject of the paper as they involve "the least cumbersome technique".

> The "computable" numbers may be described briefly as the real numbers whose expressions as a decimal are calculable by finite means.

#### Number Theory 101

- $\mathbb{N}$ ([naturals](https://en.wikipedia.org/wiki/Natural_number)): $0$, $1$, $2$, ...
- $\mathbb{Z}$ ([integers](https://en.wikipedia.org/wiki/Integer)): $-2$, $-1$, $0$, $1$, $2$, ...
- $\mathbb{Q}$ ([rationals](https://en.wikipedia.org/wiki/Rational_number)): $-1$, $1$, $\tfrac{1}{2}$, $\tfrac{7}{44}$, ...
- $\mathbb{R}$ ([reals](https://en.wikipedia.org/wiki/Real_number)): $-1$, $1$, ${\sqrt {2}}$, $\pi$, ...
- $\mathbb{C}$ ([complex](https://en.wikipedia.org/wiki/Complex_number)): $-1$, $1$, ${\sqrt {2}}$, $\pi$, $i$, $2i+3$, ...

The real numbers are all numbers which are not imaginary, and our "computable" numbers are a subset of the reals. When he says "expressions as a decimal", he is simply saying that he wants to deal with strings of digits ($0.333...$ rather than $\tfrac{1}{3}$) and in fact he will further limit his numbers to just binary digits ($0.010101...$).

By "finite means" Turing means that there must be some rule to arrive at the number without just infintely listing every digit. For example, we can describe $0.0101010101...$ finitely by saying you can repeat $01$ infinitely.

Of course, there are infinitely *random* real numbers. Whether or not we can calculate these numbers by finite means is a major area of the paper.

> According to my definition, a number is computable
if its decimal can be written down by a machine.

Turing is giving a sneak peak to the reader, telling you he will define a "machine" which will embody the calculability by finite means.

He goes on to say he will prove that large classes of the reals are computable. He then gives another sneak peak:

> The computable numbers do not, however, include
all definable numbers, and an example is given of a definable number
which is not computable.

Turing is saying he will give an example of a real number his machine cannot compute.

### Enumerabililty

> Although the class of computable numbers is so great, and in many
ways similar to the class of real numbers, it is nevertheless enumerable.
In §8 I examine certain arguments which would seem to prove the contrary.

Turing here says that computable numbers are enumerable, and then says there are arguments that prove they are not enumerable. Which is it?

By enumerable Turing means [countable](https://en.wikipedia.org/wiki/Countable_set) (I will not be providing a TL;DR on Set Theory here).

> By the correct application of one of these arguments, conclusions are
reached which are superficially similar to those of Gödel.

Turing is referring to the [diagonalization](https://en.wikipedia.org/wiki/Diagonal_lemma) used in [Gödel's incompleteness theorems](https://en.wikipedia.org/wiki/G%C3%B6del%27s_incompleteness_theorems) which we will talk about in [section 8](./GUIDE.md#section-8---application-of-the-diagonal-process).

### The Entscheidungsproblem

> These results have valuable applications. In particular, it can be shown (§11) that the
Hilbertian Entscheidungsproblem can have no solution.

Here Turing reveals **the point of the paper**. All of this work is to serve the purpose of giving a result to [the Entscheidungsproblem](https://en.wikipedia.org/wiki/Entscheidungsproblem) (in English "the *decision* problem"). We will talk in depth about the Entscheidungsproblem later in [section 11](./GUIDE.md#section-11---application-to-the-entscheidungsproblem), but here is a short description for now:

The decision problem asks if it is possible for there to be an algorithm that decides if a logic statement is **provable** from a set of axioms, for every possible statement. Again, a shortened description of the algorithm:

- Input: A logic statement, and a set of axioms
- Output: A proof of the statement's truth (or falsity) based on the set of axioms

Turing's answer to whether such an algorithm exists: *No*.

### Effective Calculability vs Computability

Apparently within the same couple of months, [Alonzo Chuch](https://en.wikipedia.org/wiki/Alonzo_Church) also came to the conclusion that there is no solution to the decision problem, and published his paper first. Turing added an Appendix that explains how Church's paper compares to his own, and gives the reader a heads up here in the introduction. Read on to the [Appendix](./GUIDE.md#appendix---computability-and-effective-calculability) if this interests you.

## Section 1 - Computing machines

### Finite Means

Turing begins this section with a preamble that says he won't attempt to justify the given definition of a computable number (one whose digits are calculable by finite means) until [section 9](./GUIDE.md#section-9---the-extent-of-computable-numbers).

He says this though:

> For the present I shall only say that the justification
lies in the fact that the human memory is necessarily limited.

A "computer" during Turing's time was an actual human performing calculations on pen and paper.

This is quite philosophical but Turing is essentially just saying that the human mind is limited to finiteness in terms of the *means* of arriving at a number, for example:

| | Finite Means | Infinite Means |
|-| ------------ | -------------- |
| Finite Number | $0.1$ | *Not possible* |
| Infinite Number | $0.010101...$ | *Infinitely random* |

### The "Machine"

Turing spends the remainder of the section giving a textual description of his "machines". These are of course the famous [Turing Machines](https://en.wikipedia.org/wiki/Turing_machine). I think his description is quite readable, so I won't try to explain it here. I will instead provide a simplified version of the type structure of the machine below (it can also be found at the top of [machine.go](./machine.go)).

```go
type (
    Machine struct {
        mConfigurations []MConfiguration
        tape            Tape
    }

    MConfiguration struct {
        Name                string
        Symbols             []string
        Operations          []string
        FinalMConfiguration string
    }

    Tape []string
)
```

## Section 2 - Definitions

Turing provides a list of definitions he will rely on later. This section is also very readable, so I'll just provide a quick reference and not belabor the details:

- **Automatic machines (a-machines)** - Machines where humans are not in the loop.
- **Computing machines** - Machines that print $0$ and $1$ as "figures", and also print other characters that help with computation.
- **Sequence computed by the machine** - The sequence of "figures" (only $0$ and $1$) computed (i.e. $010101...$).
- **Number computed by the machine** - The real number obtained by prepending a decimal place to the sequence (i.e. $0.010101...$).
- **Complete configuration** - These three things which describe the full state of the machine:
  - The full tape up to this point
  - The number of the currently scanned square
  - The name of the current m-configuration
- **Moves** - Changes of the machine and tape between successive complete
configurations.
- **Circular machine** - Machines that halt or do not print "figures" ($0$'s or $1$'s) infinitely.
- **Circle-free machine** - Must print "figures" ($0$'s or $1$'s) infinitely.
- **Computable sequence** - A sequence of "figures" that can be computed by a circle-free machine.
- **Computable number** - A number that can be derived from a computable sequence.

## Section 3 - Examples of computing machines

In this section Turing gives simple and concrete examples for his machines.

A full implementation can be found in [machine.go](./machine.go) and [machine_test.go](./machine_test.go).

**Note**: In this guide and throughout the codebase I use English letters only. Turing makes use of lowercase Greek, upper and lowercase German letters (`ə`, for example), and I will always use the English version, as its easier for me to type them. At some point I may replace all letters in the repository with Turing's original Greek/German characters. 

[machine.go](./machine.go) contains the machine's functionality:

```go
// Creates a machine
m := NewMachine(MachineInput{ ... })

// Moves the machine once
m.Move()

// Moves the machine 10 times
m.MoveN(10)

// Prints the Tape (as a string)
fmt.Println(m.TapeString())

// Print's the machine's complete configuration
fmt.Println(m.CompleteConfiguration())
```

[machine_test.go](./machine_test.go) contains Turing's three machine examples, which we test to ensure our implementation works as expected.

The final paragraph documents some conventions that Turing will always use with his machines, I'll repeat them here (along with others that he doesn't explicitly call out):
1. Begin the tape with `ee` so the machine can always find the start of the tape.
2. The square directly after `ee` is the first `F`-square (for figures). Only our sequence figures (`0` or `1`) will be printed on `F`-squares. `F`-squares will by-convention never be erased.
3. After every `F` square is an `E`-square (for erasable). In these squares Turing prints temporary characters to help with keeping track of things during computation. Otherwise they are always blank. After each `E`-square is another `F`-square.
4. We "mark" an `F`-square by placing a character directly to the right (that is, the `E`-square on the right). Turing uses "marking" extensively.

Examples:

```go
// Visually how to think about E-squares and F-squares
"eeFEFEFEFE ..."
// How it will look in practice
"ee0 1 0 1x1x0 ..."
```

Here are some implementation details to note for our `Machine`:

- Our m-configuration rows are stored in one giant list (rather than grouping each m-configuration of the same name in some structure). I found the independent m-configuration rows easier to implement.
- For the `Symbols` m-configuration field, I require that ` ` (None) be provided in addition to `*` (Any), or `!x` (Not `x`) if the m-configuration should match the blank square. This is because exactly what is meant by None, Any, Not in Turing's machine is dependent on the other m-configurations of the same name. Our implementation depends on all m-configuration rows being independent from one another, as stated above.
- Some optional fields are provided (and used by our implementation) to make things cleaner. These are `StartingMConfiguration`, `PossibleSymbols`, `NoneSymbol`, and `Debug`. These should be self-explanatory, and it should be clear by [section 7](./GUIDE.md#section-7---detailed-description-of-the-universal-machine) why they are necessary.
- You may find that I ocassionaly use `halt` as a final-m-configuration, and I never define the actual m-configuration. This is because in Turing's machines, the only way for the machine to stop (or halt) is if we are unable to find the next m-configuration after a "move". By convention, whenever I want to configure a machine to stop at some point, I will use the undefined m-configuration `halt` (however our machine, like Turing's, halts when it encounters an m-configuration that was never defined).

## Section 4 - Abbreviated tables

I got super tripped up by this section. Turing explains his "abbreviated tables" briefly and then piles them on hard. It is only with Petzold's help that I was able to figure out some of the nuances here. I'll try to start with simple example so we can work our way up (Turing starts with a complex one.) Turing's abbreviated table examples also build upon eachother (most are dependant on others) to accomplish something complex. Our four examples will be independent from one another to keep things as simple as possible.

The full implementation for this section can be found in [abbreviated.go](./abbreviated.go) and [abbreviated_tests.go](./abbreviated_test.go). It works like this:

```go
// NewAbbreviatedTable compiles our abbreviated table to normal MachineInput
m := NewMachine(NewAbbreviatedTable(AbbreviatedTableInput{
    MConfigurations: { ... }
}))

// Machines "compiled" from abbreviated tables can do anything normal machines can do
m.Move()
```

### Example 1 - Substituting symbols in `Operations`
```go
// This table moves to the right, prints `0` and repeats:
MConfigurations{
    {
        Name: "b",
        Symbols: []string{"*"},
        Operations: []string{"R"},
        FinalMConfiguration: "f(0)",
    },
    {
        Name: "f(a)",
        Symbols: []string{"*"},
        Operations: []string{"Pa"}, // We are printing whatever is passed to `f`
        FinalMConfiguration: "b",
    },
}

// and compiles to:

MConfigurations{
    {
        Name: "b",
        Symbols: []string{"*"},
        Operations: []string{},
        FinalMConfiguration: "c",
    },
    {
        Name: "c",
        Symbols: []string{"*"},
        Operations: []string{"P0"}, // Note that we substituted `a` with `0` here
        FinalMConfiguration: "b",
    },
}

// Above we have taken the liberty of choosing `c` as the "compiled" m-configuration name
// for our m-function. Turing will just use `q1`, `q2`, `q3`, ... during compilation, so
// our m-configuration should really look like:
MConfigurations{
    {
        Name: "q1",
        Symbols: []string{"*"},
        Operations: []string{},
        FinalMConfiguration: "q2",
    },
    {
        Name: "q2",
        Symbols: []string{"*"},
        Operations: []string{"P0"},
        FinalMConfiguration: "q1",
    },
}
```

So our m-functions have names and parameters. They will be called by other m-configurations in the final-m-configuration column.

### Example 2 - Substituting a `FinalMConfiguration`:

```go
// This table also moves to the right, prints `0` and repeats.
MConfigurations{
    {
        Name: "b",
        Symbols: []string{"*"},
        Operations: []string{"R"},
        FinalMConfiguration: "f(b)", // Pass the name of an m-configuration to `f`
    },
    {
        Name: "f(A)", // Receive the m-configuration name as parameter `A`
        Symbols: []string{"*"},
        Operations: []string{"P0"},
        FinalMConfiguration: "A", // Move to `A` (which is really `b`)
    },
}

// compiles to:

MConfigurations{
    {
        Name: "q1",
        Symbols: []string{"*"},
        Operations: []string{"R"},
        FinalMConfiguration: "q2",
    },
    {
        Name: "q2",
        Symbols: []string{"*"},
        Operations: []string{"P0"},
        FinalMConfiguration: "q1", // Move back to `b` (the first m-configuration)
    },
}
```
Turing's convention with parameters is that symbols are lowercase (his are Greek, ours English) and m-configurations are uppercase (his are German, ours English).

Note that in Turing's first example (`f`), there are a bunch of `f` rows, and then `f1` rows and `f2` rows. When he does this, he is saying that `f` is the m-function others will call, and anything with a number after it is just a helper for the main bit. He groups all of these together under one letter to show that they work together to offer some functionality. Its sort of like saying:

```go
// Exposed functionality
func F(a, b, c) {
    f1(a, b, c)
    f2(a, b, c)
}

// private helper
func f1(a, b, c) {
    // ...
}

// private helper
func f2(a, b, c) {
    // ...
}
```

In the paper `f` specifically stands for "find". It will go to the first `e` in the tape (the beginning of the tape), and then begin to move rightward. If it finds the desired character (`a`), it moves to the m-configuration `C`. If it cannot find `a` before it hits two blank squares in a row it will move to m-configuration `B`. Our implementation of `f` can be found at the top of [abbreviated.go](./abbreviated.go).

### Example 3 - Functions within functions
```go
// This table moves to the right, printing 0 twice in a row, and continuing on infinitely.
// Example: " 00 00 00"
MConfigurations{
    {
        Name: "b",
        Symbols: []string{"*"},
        Operations: []string{"R"},
        // Take note here, the first move to `f` will have params `f(b, 0)`, and `0`.
        // The second move to `f` will have params `b` and `0`.
        FinalMConfiguration: "f(f(b, 0), 0)",
    },
    {
        Name: "f(A, a)",
        Symbols: []string{"*"},
        Operations: []string{"Pa", "R"}, // Prints `a`, and moves to the right
        FinalMConfiguration: "A", // Moves to the m-configuration provided as parameter `A`
    },
}

// compiles to:

MConfigurations{
    {
        Name: "q1",
        Symbols: []string{"*"},
        Operations: []string{"R"},
        FinalMConfiguration: "q2", // First move to `f`
    },
    {
        Name: "q2",
        Symbols: []string{"*"},
        Operations: []string{"P0", "R"},
        FinalMConfiguration: "q3", // Second move to `f` 
    },
        {
        Name: "q3",
        Symbols: []string{"*"},
        Operations: []string{"P0", "R"},
        FinalMConfiguration: "q1", // Moves back to `b`
    },
}
```

### Example 4 - Symbol parameters
This one was fun. If there is a symbol (in the Symbol column) of our abbreviated table, but it is not a symbol the machine prints and also not a parameter of our m-function, it is a "symbol parameter". We basically "read" the symbol, and pass whatever it was as a parameter. By convention I will always prepend symbol parameters with an underscore (`_`) for clarity.

```go
// This table moves to the right, copying the symbol it was just looking at (infinitely).
// Assume the table is only capable of printing `0`` and `1`.
MConfigurations{
    {
        Name: "b",
        Symbols: []string{"_y"}, // Our "symbol parameter"
        Operations: []string{"R"},
        FinalMConfiguration: "f(_y)", // Pass the symbol parameter to `f`
    },
    {
        Name: "f(a)",
        Symbols: []string{"*"},
        Operations: []string{"Pa", "R"}, // where it is simply printed
        FinalMConfiguration: "b",
    },
}

// compiles to:

MConfigurations{
    {
        Name: "q1",
        Symbols: []string{"0"}, // We must enumerate this m-function once for each possible symbol
        Operations: []string{"R"},
        FinalMConfiguration: "q2",
    },
    {
        Name: "q1",
        Symbols: []string{"1"}, // Here is the other possible symbol
        Operations: []string{"R"},
        FinalMConfiguration: "q2",
    },
    {
        Name: "q2",
        Symbols: []string{"*"},
        Operations: []string{"P0"}, // The version of `f` that prints `0`
        FinalMConfiguration: "q1",
    },
        {
        Name: "q2",
        Symbols: []string{"*"},
        Operations: []string{"P1"}, // The version of `f` that prints `0`
        FinalMConfiguration: "q1",
    },
}
```

Throughout the rest of the section, Turing gives a bunch of m-functions which he will use later in [section 7](./GUIDE.md#section-7---detailed-description-of-the-universal-machine). These are enumerated with comments at the top of [abbreviated.go](./abbreviated.go) and are tested in [abbreviated_tests.go](./abbreviated_test.go). Most of them are helpers for copying, erasing, printing, etc.

## Section 5 - Enumeration of computable sequences

In this section Turing tells us how to "standardize" his tables. This standardization occurs in three steps:

- **Standard Table**: Modify the m-configurations (which could include adding new ones) such that there is only one Symbol, one Print operation, and one Move operation.
- **Standard Description**: Take these m-configurations, and convert them to one long string.
- **Description Number**: Convert our Standard Description string into a number.

Performing these standardizations allow us to do some interesting things. The first is to put the Standard Description on a Tape and compute on it (the Standard Description is just a string of symbols), which is explored in [section 6](./GUIDE.md#section-6---the-universal-computing-machine). Another is to treat a Description Number just like any other number (which allows to treat Machines like numbers) - this is explored in [section 8](./GUIDE.md#section-8---application-of-the-diagonal-process).

### Implementation Details

This entire section is quite readable, and its implementation can be found in [standard.go](./standard.go), and tested in [standard_test](./standard_test.go). Here is how it works:
```go
standardTable := NewStandardTable(MachineInput{
    MConfigurations: []MConfiguration{
        {"b", []string{" "}, []string{"P0", "R"}, "c"},
        {"c", []string{" "}, []string{"R"}, "e"},
        {"e", []string{" "}, []string{"P1", "R"}, "k"},
        {"k", []string{" "}, []string{"R"}, "b"},
    },
    PossibleSymbols: []string{"0", "1"}, // Required so we can convert `*`, `!x`, etc.
})

// Can be used just like any other MachineInput
machineInput := standardTable.MachineInput
machine := NewMachine(machineInput)
machine.Move(50)

// Capable for converting the resulting Tape from running this Machine back to the original symbols
symbolMap := standardTable.SymbolMap
fmt.Println(symbolMap.TranslateTape(machine.Tape()))

// Turing's Standard Description
fmt.Println(standardTable.StandardDescription)

// Turing's Description Number
fmt.Println(standardTable.DescriptionNumber)
```

## Section 6 - The universal computing machine

This section, and the [next section](./GUIDE.md#section-7---detailed-description-of-the-universal-machine) are dedicated to explaining and implementing Turing's "universal computing machine". Turing explains that `U` is a Machine that is capable of computing the same sequence as any other Machine `M`, provided the Standard Description of `M` is supplied as the Tape of `U`.

### The first programmable machine
It is worth stopping to think about this for a minute. Turing is saying we can create a table of m-configurations that take a Machine as input and reproduce that same Machine's output. As far as I am aware, this is the first conceptualization of software in the modern sense. With `U`, we no longer have to configure Machines directly via m-configurations any more (analogous to hardward). We simply print out our desired Machine on a Tape and feed it into `U` (in other words, we just "program" `U`).

### Outline of `U`

Turing now gives an overview of how he can achieve `U` by breaking the problem into subparts. The first is to build a machine `M'` that will print the complete configuration of `M` on the `F`-squares. The complete configuration is a full history of the moves of a machine. He suggests that we encode the complete configuration in the same standard form as the Standard Description itself.

Even if we were to accomplish `M'`, we would need to somehow locate `M`'s Tape output somewhere. Turing suggests that this output be interweaved between complete configurations between colons, like this:
```go
" ... : 0 : ... : 1 : ... "
``` 

So we don't actually get the *exact* output of `M` (there will be giant complete configurations between the actual computed sequence), but `U` would indeed "compute" the same sequence as `M`.

What follows in the next section is fascinating - all on paper Turing builds `U`, arguably the first computer.

## Section 7 - Detailed description of the universal machine

Here Turing gives the full table of m-configurations and m-functions for `U`, relying on the helper m-functions he created in [section 4](./GUIDE.md#section-4---abbreviated-tables). Here is an outline of the actual working of `U`:

1. `U` takes a Tape starting with `ee`, and then `M`'s Standard Description on the `F`-squares, and finally a double-colon (`::`).
2. ... TODO

Note that there are at least 4 small bugs in Turing's original paper in this section. Fortunately, Petzold compiled a list of fixes for them (originally spotted by [Davies](https://en.wikipedia.org/wiki/Donald_Davies) and [Post](https://en.wikipedia.org/wiki/Emil_Leon_Post)).

### Implementation Details

Our implementation is in [universal.go](./universal.go) and tested in [universal_test.go](./universal_test.go). Here is the interface:

```go
machineInput := MachineInput{
    MConfigurations: []MConfiguration{
        {"b", []string{" "}, []string{"P0", "R"}, "c"},
        {"c", []string{" "}, []string{"R"}, "e"},
        {"e", []string{" "}, []string{"P1", "R"}, "k"},
        {"k", []string{" "}, []string{"R"}, "b"},
    },
}

standardTable := NewStandardTable(machineInput)

univeralMachineInput := NewUniversalMachine(UniversalMachineInput{
    StandardDescription: st.StandardDescription,
    SymbolMap:           st.SymbolMap,
})
universalMachine := NewMachine(univeralMachineInput)
universalMachine.MoveN(500000)

// Should be the same as if we just created a Machine from machineInput
fmt.Println(universalMachine.TapeStringFromUniversalMachine())
```

## Section 8 - Application of the diagonal process

TODO

## Section 9 - The extent of computable numbers

TODO

## Section 10 - Examples of large classes of numbers which are computable

TODO

## Section 11 - Application to the Entscheidungsproblem

TODO

## Appendix - Computability and effective calculability

TODO

## Overall Thoughts

- Turing essentially invents (in theory) the computer and the concept of software/programming, all in service to solve a math problem.
- It felt as if Turing piled three or four genius insights on-top of one another, so while reading I was never able to get comfortable before he took things to the next level.
- There are a lot of bugs (which makes sense as he was not able to run the machine himself).
- Attempting to learn this paper as a non-mathematician makes me want to get a deeper understanding of the history of math/logic (probably starting with Hilbert, Gödel, etc.)
- The appendix is a great segue to Church's Lambda Calculus, which will probably be my next project.