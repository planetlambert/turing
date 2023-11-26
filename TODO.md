# TODO

1. Finish Guide
   1. Section 11 Lemma 1
   1. Appendix
1. Finish implementation of `H` and `G` machines in [diagonal.go](./diagonal.go) and [diagonal_test.go](./diagonal_test.go)
   1. Requires refactoring the Universal Machine to use custom sentinel values (`:` and `::`)
   1. Thought: Chain machines together using Go rather than using Machine Tape
1. Finish implementation of `K` and `Ka`machines in [hilbert.go](./hilbert.go) and [hilbert_test.go](./hilbert_test.go)
1. Finish implementation of the impossible undecidable machine in [decision.go](./decision.go) and [decision_test.go](./decision_test.go).
1. Go back to section 10 some day