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
2. `b`: `U` starts by finding `::`, printing `: D A`, and moving to m-configuration `anf`. You can think of `b` as initialization/setup, and `anf` as the start of `U`'s main loop.
3. `anf`: Find the last complete configuration, and mark out the following:
    1. The m-configuration name (somewhere within the complete configuration)
    2. The symbol directly to the right (which is the scanned square).
    3. This marking is done by printing `y` to the right of these squares.
4. `kom`: Moving leftwards, find the m-configuration within the Standard Description that matches that marked by `y`. This is done by:
    1. Finding a semi-colon (representing the start of an m-configuration)
    2. Marking the first two sections (m-configuration name and symbol) with `x`.
    3. Comparing all squares marked with `y` with those marked with `x` (using `kmp`).
    4. If they don't match, mark the semi-colon with `z`, and start again (this time skipping all semi-colons marked `z`).
    5. If they *do* match, we found our m-configuration, move on to `sim`.
    6. In both cases all `x`'s and `y`'s are erased.
5. `sim`: This step does three things:
    1. Mark the symbol within the Print operation with `u`.
    2. Mark the final m-configuration with `y`.
    3. Move on to `mk`.
6. `mk`: Returns to the complete configuration, and divides into four sections via markers:
    1. `x`: The symbol directly to the left of the m-configuration name.
    2. `v`: Every symbol before `x`.
    3. ` `: Skip over the m-configuration name, and the scanned symbol.
    4. `w`: Everything to the right of the scanned symbol.
    5. Finally, we move onto `sh`
7. `sh`: This is the section that prints the actual symbol from `M` (if applicable). We do this by:
    1. Check if the last square scanned (`u`), is blank.
    2. If it is indeed blank, we are writing a new character in the sequence for `M`. Print that character after the complete configuration between two colons (as described in [section 6](./GUIDE.md#section-6---the-universal-computing-machine)).
    3. Move on to `inst`
8. `inst`: Finally, write the next complete configuration.
    1. We already have all of the relevant sections marked out with `v`, `y`, `x`, `u`, and `w`. We just need to stitch them together.
    2. This is done depending on if our Move operation is `L`, `R`, or `N`.
        1. All three options below print characters marked by `v` first, and `w` last.
        2. `L`: `y`, then `x`, then `u`.
        3. `R`: `x`, then `u`, then `y`.
        4. `N`: `x`, then `y`, then `u`.
    3. After we print the complete configuration in the correct order, start again (go back to `anf`, which is out #3.)

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

This section returns to the question of whether the computable numbers are enumerable, and the nuances that determine the difference between computable numbers and real numbers.

### Cantor's diagonal argument

This entire section relies on the reader understanding the general idea of [Cantor's diagonal argument](https://en.wikipedia.org/wiki/Cantor%27s_diagonal_argument). The gist of the argument is quite intuitive, and I think after watching a quick video (I recommend [Numberphile's](https://www.youtube.com/watch?v=elvOZm0d4H0)) you should be up-to-speed.

### An incorrect application of the diagonal process

In the first paragraph Turing notes that one may think the computable numbers are not enumerable for the same reason (Cantor's diagonal argument) as the real numbers. The next sentence tripped me up:

> It might, for instance, be thought that the limit of a sequence of computable numbers must be computable. This is clearly only true if the sequence of computable numbers is defined by some rule.

Turing is simply saying that computable numbers have an extra property (they must actually be calculable by finite means), so you wouldn't be able to do necessarily do things that you would do to the real numbers, and expect another computable number to result. You might be able to, but you might need some additional rules on top to ensure you arrive at a computable number.

To show this, Turing now fallaciously applies the diagonal argument to computable numbers. He first assumes the computable sequences are enumerable, and calculates the "sequence on the diagonal" which he calls `b` (he uses the lowercase Greek letter beta). Given that `b` is computable, he uses some math (the actual math doesn't matter) to prove that `b` cannot exist (and therefore computable sequences are not enumerable).

Turing points out that the reason this argument doesn't hold is because of the assumption that `b` is computable. The reason we thought `b` was computable in the first place is because we "computed" it from the enumeration of computable sequences. But referring back to Turing's definiton of computable (calculable by finite means), this would necessarily mean that the process of enumeration would have to be finite. The next few paragraphs are dedicated to proving that this process of enumerating computable sequences cannot be done in a finite number of steps.

### Does `D` exist?

To prove that we can't enumerate computable sequences by finite means, Turing uses a proof by contradition. He assumes we *can* enumerate computable sequences by finite means, and finds a paradox.

I will make a small note here on the phrase `circle-free`. The wording is unclear which is unfortunate, but its worth digging at what Turing is getting at. By `circle`, Turing is referring to some "bug" that prevents a machine from making progress in printing its sequence. A `circular` machine is one that will never make progress (or may even halt) when attempting to print its sequence. `circle-free` machines *eventually* make progress.

Back to the proof. Turing supposes the existence of two machines `D` and `H`:

- `D` - When supplied with an S.D., it will print `s` (for satisfactory) if the S.D. represents a machine that is circle-free, or `u` (for unsatisfactory) if the S.D. represents a machine is circular.
  - It is assumed that any scratch work of `D` will be erased.
  - Note that in order to be satisfactory, the machine must also be well-defined (the S.D. actually represents a functioning machine).
  - `D` is a part of `H`.
  - The inner workings of `D` are not provided by Turing.
- `H` - Relying on `D` and `U`, `H` will print `b'` (the "sequence on the diagonal", without the modification by $1$ to make it unique). `H` does this by:
  - Start `N` at the Description Number $1$, and at each step increment ($2$, $3$, ...).
  - Keep a running counter `R` for what place we are at in the diagonal.
  - At each step, do the following:
    -  Provide the S.D. of `N` to `D`. If it is satisfactory (`s`), increment `R`. We say that `R` at the `N`th step is `R(N)`.
       - For example, if `N` is the first circle-free machine found, then `R(N)` is $1$.
    - `H` then must find the `R(N)`th digit of the sequence for the machine described by `N`, which is the `R(N)`th digit of `b'`.

Here is a table outlining the steps of `H` (with thanks to Charles Petzold):

| N             | well-defined? | circle-free? | R(N) | `b'` so far | `b` so far | note      |
| ------------- | ------------- | ------------ | ---- | ----------- | ---------- | --------- |
| 1             | No            | N/A          | 0    | N/A         | N/A        |           |
| 2             | No            | N/A          | 0    | N/A         | N/A        |           |
| ...           |               |              |      |             |            |           |
| 313,325,317   | Yes           | Yes          | 1    | 1           | 0          | prints 1s |
| ...           |               |              |      |             |            |           |
| 3,133,255,317 | Yes           | Yes          | 2    | 10          | 01         | prints 0s |
| ...           |               |              |      |             |            |           |


If you found the description of `H` is complicated, here is a more intuitive way to think about it: `H` just loops over all natural numbers, converts the number to an S.D., tests if the S.D. is circle-free, and then finds the appropriate digit to add to the sequence on the diagonal. If it is still too complicated, I have implemented `H` (the parts that are possible) in [diagonal.go](./diagonal.go) and [diagonal_test.go](./diagonal_test.go). I think its worth understanding `H` in full before we get to the proof by contradition.

Given this description, if `H` in its entirety (including `D`) is circle-free, then the computable numbers are obviously enumerable by finite means. For his proof, let us assume that `H` is indeed circle-free.

Before he gives the proof, Turing explains that the non-`D` parts of `H` are clearly circle-free (and in fact in [diagonal_test.go](./diagonal_test.go) we proved this). This makes sense because it basically just loops over the natural numbers, outsources the hard part to `D`, keeps track of the count of circle-free machines (`R`), grabs the `R`th digit of a sequence and prints it. Turing does this because he is gearing up to accuse some part of `H` of not being circle-free, and he doesn't want the blame to fall on these simple parts of `H`.

Now comes the proof. Turing supposes that over the course of `H` looping over the Description Numbers of all machines, it will have to (at some point) arrive at the Description Number of `H` itself (which he says is `K`). Its fascinating to think of what happens at this step. `D` will have to determine if `H` is satisfactory (circle-free) or unsatisfactory. `H` can't come back unsatisfactory, because we just assumed `H` is circle-free. At the same time, it cannot be satisfactory - when trying to calculate the `R(K)`th digit in our diagonal sequence, we must go through the entirety of `H`'s steps once more, and in fact we have to do this recursion infinitely. There is an infinite loop, or in Turing's terms, `H` is circular (but we just assumed it was satisfactory, or circle free). `H` paradoxically cannot be satisfactory or unsatisfactory, so `D` cannot exist.

With this proof, Turing shows:
1. No machine exists that can determine whether another machine is circle-free
2. Diagonalization cannot be applied to computable numbers
3. Computable numbers *are* theoretically enumerable in the sense that we can:
   1. Enumerate the natural numbers
   2. Use them to represent machines
   3. Possibly compute a sequence using that machine, **however**...
4. ...you cannot actually perform this enumeration by finite means!

To answer our original question on enumerability of computable numbers: it depends on our definition of enumerability. My opinion is that to be enumerable, someone must be able to actually perform the enumeration, and therefore they are **not** enumerable.

### Does `E` exist?

Turing now employs a similar proof by contradiction, this time for a machine `E` which he will prove useful later. He explains the goal of the proof:

> there can be no machine E which, when supplied with the S.D of an arbitrary machine M, will determine whether M ever prints a given symbol (0 say).

He first supposes a machine `M1` which prints the same sequence as `M` except for it replaces the first printed 0 with a 0̄. Similarly `M2` replaces the first two printed 0's with 0̄'s. These machines are quite easy to implement and are contained in [diagonal.go](./diagonal.go) and [diagonal_test.go](./diagonal_test.go).

Next comes the machine `F` which prints the S.D. of `M1`, `M2`, ..., etc. successively. Turing now supposes the machine `G` which combines `F` and `E`. It uses `F` to loop over `M1`, `M2`, etc., and then uses `E` to test if 0 is ever printed. For each step in the loop, if it is found that a 0 is never printed, then `G` itself will print a 0.

Now Turing turns this in on itself and has `E` test `G`. Because `G` only prints 0 when there is no 0 printed by `M`, we can tell if `M` prints 0 infinitely often by checking if `G` never prints a 0. Now we have a way to determine if 0 is printed infinitely often by `M`. It should be clear that we can similarly determine if 1 is printed infinitely often using the same tactic.

Here is the crux of the proof - the problem is that if we have a way of determining if 0's and 1's are printed infinitely often, then we can tell if a machine is circle-free. We already know that this is impossible from our result in the sub-section above. Therefore, `E` cannot exist.

### There is a general process for determining...

The final paragraph in this section is a small aside that leads us into the next section, and makes a connection between Turing's machine and mathematical logic. The first part explains that we can't assume that "there is a general process for determining ..." means that "there is a machine for determining ...". We have to prove that our machines can truly do anything that human computers can do first.

He gives us a preview of how he will do this: he says that each of these "general process" problems can be broken down into a mathematical logic problem. The sequence computed by a machine can be the answer to the logic problem "is `G(n)` true or false" where `n` is any integer and `G` is any logical function. For example:

| Computed Sequence | Logical Function          |
| ----------------- | ------------------------- |
| `10000000000000...` | `IsZero(n)`             |
| `01000000000000...` | `IsOne(n)`              |
| `10101010101010...` | `IsEven(n)`             |
| `01010101010101...` | `IsOdd(n)`              |
| `10010010010010...` | `IsDivisibleByThree(n)` |
| `10010010010010...` | `IsDivisibleByThree(n)` |

Turing uses this connection between his machines and logic for the rest of the paper. With this insight he can now attempt to create a process (as in *decision process*) that computes any logical function.

## Section 9 - The extent of computable numbers

Until now, Turing has not explained the reasoning behind his definition of "computable". Here is an excerpt from section 1:

> We have said that the computable numbers are those whose decimals are calculable by finite means. This requires rather more explicit definition. No real attempt will be made to justify the definitions given until we reach § 9. For the present I shall only say that the justification lies in the fact that the human memory is necessarily limited.

Turing now will now give some arguments that describe the "extent" of the computable numbers. By being more philosophically rigorous he will be able to build on top of these arguments, which gives him an angle from which to attack the decision problem later.

The secret to understanding section 9 is this: Turing wants to prove that his machine is capable of "computing" anything that a human (or anything else) is capable of "computing". If Turing can convince the reader of this, he can use his machines in proofs related to mathematical logic without the reader worrying that the proofs don't apply universally. Without this section, a reader might object to Turing's proof in the following way:

"Sure, Turing proved *his machines* are incapable of accomplishing X, but that doesn't necessary prove that X is unaccomplishable..." - Someone who skipped section 9.

Turing attempts to convince the reader using three separate arguments. Before diving into the arguments, he gives us a preview of how he will build upon this foundation:

> Once it is granted that computable numbers are all "computable", several other propositions of the same character follow. In particular, it follows that, if there is a general process for determining whether a formula of the Hilbert function calculus is provable, then the determination can be carried out by a machine.

Here he is simply saying that if the answer to Hilbert's decision problem is `yes` (that there is an algorithm that can decide if a logic statement is **provable** from a set of axioms, for every possible statement), then he will be able to construct a machine that carries out this decision process.

### Argument `a` - A direct appeal to intuition

The first of Turing's arguments should be easy for the modern reader to understand. Turing is essentially convincing the reader that his machines are "[Turing Complete](https://en.wikipedia.org/wiki/Turing_completeness)" (a term Turing did not have access to at the time). For us, it is a long-established fact that there is nothing "special" happening in the human brain (there is nothing in the brain that is not replicatable by a machine that is powerful enough).

In Turing's time this was not the case, and he went into painstaking detail about why this intuitively true. This argument in its entirety is quite readable, so I won't add anything else here except that when Turing mentions "computer" he is referring to a human performing computations.

### Argument `b` - A proof of the equivalence of two definitions

Turing's second argument is more complex, and will require an understanding of [first-order logic](https://en.wikipedia.org/wiki/First-order_logic), and specifically [Hilbert's calculus](https://en.wikipedia.org/wiki/Hilbert_system). I won't explain these in detail, but here is a crash course, specifically in Hilbert's (and Turing's) notation:

#### First-order logic 101

- Propositions are sentences that can be true or false.
  - Ex: It will snow today
  - Repesented by capital letters ($X$, $Y$, etc.)
- $\vee$ represents OR
  - Ex: $X \vee Y$ is true if $X$ or $Y$ is true
- $\\&$ represents AND
  - Ex: $X \\& Y$ is true if both $X$ and $Y$ are true
- $-$ represents NOT
  - Ex: $-X$ means $X$ is not true
- $→$ represents implication
  - $X → Y$ is equal to $-X \vee Y$
- $\sim$ represents equality
  - $X \sim Y$ is equal to $(X → Y) \\& (Y → X)$
- Parentheses $()$ denote evaluation order
- Predicates are functions over natural numbers that evaluate to a truth value
  - Ex: $\text{IsPrime}(x)$ where $\text{IsPrime}(4)$ is false and $\text{IsPrime}(5)$ is true
- $\exists$ represents existential quantification
  - Ex: $(\exists x)\text{IsPrime}(x)$ means there exists a natural number that is prime.
  - $(\exists x)\text{B}(x)$ is equal to $\text{B}(0) \vee \text{B}(1) \vee \text{B}(2) ...$ if $x$ represents natural numbers
- $(x)$ represents universal quantification
  - Ex: $(x)\text{IsPrime}(x)$ means that all natural numbers are prime (which is of course false)
  - $(x)\text{B}(x)$ is equal to $B(0) \\& B(1) \\& B(2) ...$ if $x$ represents natural numbers

Turing wants to use a version of Hilbert's calculus that is modified slightly (he wants to use a finite set of natural numbers, etc.). This is just so he can prove that his machine that simulates first-order logic is finite, and therefore whatever is calculated is a "computable number" 

#### Equivalence

Back to argument `b`. This argument has the following outline:
1. Turing's machines are capable of representing and simulating Hilbert's calculus (the two definitions are "equivalent").
2. Numbers defined by Hilbert's calculus in the way Turing describes include *all* computable numbers (just simulate the correct first-order logic formula).
3. Therefore, Turing's machine can be directly applied to problems relating to Hilbert's calculus (like the decision problem)

Turing shows (1) and (2), while (3) is implied. Lets get into the details of the argument.

#### `K` exists

Turing first presupposes that a machine `K` exists (it actually does this time...) that can find all provable formulae of Hilbert's calculus. We know that `K` exists because we can recursively enumerate all formulae from a set of axioms, a consequence of [Gödel's completeness theorum](https://en.wikipedia.org/wiki/G%C3%B6del%27s_completeness_theorem).

How does this actually work though? `K` is implemented in [hilbert.go](./hilbert.go) and [hilbert_test.go](./hilbert_test.go) to show you, but the secret is to just brute force every combination of operators and axioms using prefix notation.

#### Peano Axioms

Turing chooses these axioms to be [Peano axioms](https://en.wikipedia.org/wiki/Peano_axioms), which Turing defines as $P$, or:

$$(\exists u)N(u) \\;\\; \\& \\;\\; (x)(N(x) → (\exists y)F(x, y)) \\;\\; \\& \\;\\; (F(x, y) → N(y))$$

where:

- $N(x)$ - $x$ is a natural number
- $F(x, y)$ - The successor function ($y$ is one greater than $x$)

*Aside: Petzold explains Turing is missing some axioms to ensure the uniqueness of zero, uniqueness of the successor, etc., but lets just assume $P$ correctly enumerates the Peano axioms.*

You can think of Peano's axioms as a way to bootstrap mathematics within first-order logic. Using Peano's axioms you can do things like actually defining $\text{IsPrime}(x)$ using only first-order logic.

#### Computing a sequence

Turing now explains that we will be attempting to compute a sequence $a$, and provides some predicates that will help later:

- $G_a(x)$ - The $x$'th figure of $a$ is $1$
- $-G_a(x)$ - The $x$'th figure of $a$ is $0$

Using $P$ and any other combination of propositions built from the axioms, we can define a formula $U$ (Turing uses the [Fraktur](https://en.wikipedia.org/wiki/Fraktur) letter A which is not available in GitHub's LaTeX) which give the foundation for computing $a$ using first-order logic. Note that these "other combination of propositions" are what makes $a$ unique.

Finally we have two formulas are $A_n$ and $B_n$, which represent the following:

- $A_n$ is the formula built up from our axioms that implies that the $n$'th figure of $a$ is $1$.
- $B_n$ is the formula built up from our axioms that implies that the $n$'th figure of $a$ is $0$.

It should be clear that only $A_n$ or $B_n$ can be true. Turing uses more successor functions in $A_n$/$B_n$ to ensure that we are keeping things finite.

Now we have everything we need to describe a machine $K_a$ that can compute $a$:

- For each motion $n$ of $K_a$ the $n$'th figure of $a$ is computed.
- In motion $n$:
  - Write down $A_n$ and $B_n$ on the tape
  - Enumerate all theorums deduced from the set of axioms using `K`
  - Eventually we will find either $A_n$ or $B_n$ in this enumeration. If it is $A_n$, print $1$. If it is $B_n$, print $0$.

To recap what is happening here: We have created a machine ($K_a$) that computes the sequence $a$ digit-by-digit. It does this by checking if our custom first-order logic formula ($U$) is true or false for a given digit (via brute force) and printing $1$ or $0$ accordingly. Our implementation of $K_a$ can be found in [hilbert.go](./hilbert.go) and [hilbert_test.go](./hilbert_test.go).

Turing sews the argument up by explaining that all computable numbers can be derived this way (we would just need the correct formula $U$). He then gives a gentle reminder that these computable numbers are not all definable numbers (which we learned in section 8).

In the final subsection (III), Turing relates the axiom system described above with the human argument from subsection (I). It's a bit philosophical, but I believe Turing is dispelling doubts readers may have about the human argument from subsection (I). Instead of having to rely on the muddy "state of mind" concept, we can replace it with a more formulaic one as described in subsection (II).

Argument `c` is left to [section 10](./GUIDE.md#section-10---examples-of-large-classes-of-numbers-which-are-computable).

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
- Attempting to learn this paper as a non-mathematician makes me want to get a deeper understanding of the history of math/logic (probably via [Frege to Gödel](https://www.amazon.com/Frege-Godel-Mathematical-1879-1931-Sciences/dp/0674324498))
- The appendix is a great segue to Church's [Lambda Calculus](https://en.wikipedia.org/wiki/Lambda_calculus), which will probably be my next project.