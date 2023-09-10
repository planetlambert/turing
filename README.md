# Turing

An open source implementation of Alan Turing's famous paper, [On Computable Numbers, with an Application to the Entscheidungsproblem](https://www.cs.virginia.edu/~robins/Turing_Paper_1936.pdf).

## Why?

I wanted to read Turing's paper, and found it too difficult to understand. I couldn't find a complete and easily-readable reference implementation, so I decided to make my own.

## How to use this repository
***Disclaimer**: There are still large swaths of the paper (mainly the logic parts of sections 8 through 11) that are in progress as I don't understand well enough to explain to others. If you have a good resource to help a non-mathematician understand these sections, please reach out!*

For those who intend to read the paper I recommend starting with [The Annotated Turing](https://www.amazon.com/Annotated-Turing-Through-Historic-Computability/dp/0470229055) by Charles Petzold (which explains the paper line-by-line with annotations), alongside this repository's [GUIDE.md](./GUIDE.md) which will guide you through the paper and codebase section-by-section.

For those who want to use the implementation, here is how to get started:

```shell
go get github.com/planetlambert/turing@latest
```

```go
import (
    "fmt"

    "github.com/planetlambert/turing"
)

func main() {
    m := &turing.Machine{
        MConfigurations: turing.MConfigurations{
            {Name: "b", Symbols: []string{" "}, Operations: []string{"P0", "R"}, FinalMConfiguration: "c"},
            {Name: "c", Symbols: []string{" "}, Operations: []string{"R"},       FinalMConfiguration: "e"},
            {Name: "e", Symbols: []string{" "}, Operations: []string{"P1", "R"}, FinalMConfiguration: "k"},
            {Name: "k", Symbols: []string{" "}, Operations: []string{"R"},       FinalMConfiguration: "b"},
        },
    }
    m.MoveN(50)
    fmt.Println(m.TapeString()) // Prints "0 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 0 1 0"
}
```

[Full Documentation here.](https://pkg.go.dev/github.com/planetlambert/turing)

## Progress
- [X] [Introduction](./GUIDE.md#introduction)
- [X] [Section 1 - Computing machines](./GUIDE.md#section-1---computing-machines)
- [X] [Section 2 - Definitions](./GUIDE.md#section-2---definitions)
- [X] [Section 3 - Examples of computing machines](./GUIDE.md#section-3---examples-of-computing-machines)
- [X] [Section 4 - Abbreviated tables](./GUIDE.md#section-4---abbreviated-tables)
- [X] [Section 5 - Enumeration of computable sequences](./GUIDE.md#section-5---enumeration-of-computable-sequences)
- [X] [Section 6 - The universal computing machine](./GUIDE.md#section-6---the-universal-computing-machine)
- [X] [Section 7 - Detailed description of the universal machine](./GUIDE.md#section-7---detailed-description-of-the-universal-machine)
- [ ] Section 8 - Application of the diagonal process
- [ ] Section 9 - The extent of computable numbers
- [ ] Section 10 - Examples of large classes of numbers which are computable
- [ ] Section 11 - Application to the Entscheidungsproblem
- [ ] Appendix - Computability and effective calculability

## FAQ
- Why Go?
  - I like Go, and it is the most readable language for me currently.
- How is the performance?
  - My goal for this repository is to be a learning resource, so when possible I bias towards readability.
