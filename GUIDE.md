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
- You may find that I ocassionaly use `halt` as a final-m-configuration, and I never define the actual m-configuration. This is because in Turing's machines, the only way for the machine to stop (or halt) is if we are unable to find the next m-configuration after a "move". By convention, whenever I want to configure a machine to stop at some point, I will use the undefined m-configuration `halt`.

## Section 4 - Abbreviated tables

I got super tripped up by this section. Turing explains his "abbreviated tables" briefly and then piles them on hard. It is only with Petzold's help that I was able to figure out some of the nuances here. I'll try to start with simple example so we can work our way up (Turing starts with a complex one.)

The full implementation for this section can be found in [abbreviated.go](./abbreviated.go) and [abbreviated_tests.go](./abbreviated_test.go). It works like this:

```go
m := NewMachine(NewAbbreviatedTable(AbbreviatedTableInput{
    MConfigurations: { ... }
}))
```

### Example 1 - Substituting symbols in `Operations`
```go
// This table prints `a` and repeats
MConfigurations{
    {
        Name: "b",
        Symbols: []string{"*"},
        Operations: []string{},
        FinalMConfiguration: "f(0)",
    },
    {
        Name: "f(a)",
        Symbols: []string{"*"},
        // We are printing whatever is passed to `f`
        Operations: []string{"Pa"}, 
        FinalMConfiguration: "b",
    },
}

// compiles to

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
         // Note that we substituted `a` with `0` here
        Operations: []string{"P0"},
        FinalMConfiguration: "b",
    },
}

// We have taken the liberty of choosing `c` as the m-configuration name
// for our compiled m-function. Turing will just use `q1`, `q2`, `q3`, ...
// during compilation, so our m-configuration should really look like

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
MConfigurations{
    {
        Name: "b",
        Symbols: []string{"*"},
        Operations: []string{},
        FinalMConfiguration: "f(c)",
    },
    {
        Name: "f(A)",
        Symbols: []string{"*"},
        Operations: []string{"P0"},
        // We move to whatever m-configuration was passed as a parameter to `f`
        FinalMConfiguration: "A",
    },
}

// compiles to

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
        // The parameter (`c`) was compiled to q3
        FinalMConfiguration: "q3",
    },
}
```
What is `c` in this example? We never defined it so the machine will just halt when it gets to `q3`.

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
TODO
```

```

### Example 4 - Symbol parameters
TODO
```

```

## Section 5 - Enumeration of computable sequences

TODO

## Section 6 - The universal computing machine

TODO

## Section 7 - Detailed description of the universal machine

TODO

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