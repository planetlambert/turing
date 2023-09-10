# Guide

This guide and codebase can be thought of as a companion resource for those reading [On Computable Numbers, with an Application to the Entscheidungsproblem](https://www.cs.virginia.edu/~robins/Turing_Paper_1936.pdf) by Alan Turing.

I can't recommend [The Annotated Turing](https://www.amazon.com/Annotated-Turing-Through-Historic-Computability/dp/0470229055) enough, this project would not have been possible without it.

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
- $\mathbb{C}$ ([complex](https://en.wikipedia.org/wiki/Real_number)): $-1$, $1$, ${\sqrt {2}}$, $\pi$, $i$, $2i+3$, ...

The real numbers are all numbers which are not imaginary. When he says "expressions as a decimal", he is simply saying that he wants to deal with strings of digits ($0.333...$ rather than $\tfrac{1}{3}$) and in fact he will further limit his numbers to just binary digits ($0.010101...$).

By "finite means" he means that there must be some rule to arrive at the number without just infintely listing every digit. For example, we can describe $0.0101010101...$ finitely by saying you can repeat $01$ infinitely.

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

**The point of the paper**. All of this work is to serve the purpose of giving a result to [the Entscheidungsproblem](https://en.wikipedia.org/wiki/Entscheidungsproblem) (in English "the *decision* problem"). We will talk in depth about the Entscheidungsproblem later in [section 11](./GUIDE.md#section-11---application-to-the-entscheidungsproblem), but here is a short description for now:

The decision problem asks if it is possible for there to be an algorithm that decides if a logic statement is **provable** from a set of axioms, for every possible statement. Again, a description of the algorithm:

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
        MConfigurations
        Tape
    }

    MConfigurations []MConfiguration

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

TODO

## Section 3 - Examples of computing machines

TODO

## Section 4 - Abbreviated tables

TODO

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