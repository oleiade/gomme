# TODOS

## Dos

- [ ] Create `Bytes` interface type for bytes file's content
- [ ] Create `String` or `Characters` interface type for characters file's content
- [ ] Sort Out Fatal/Non-Fatal errors (distinguish whether a parser failed in an expected manner, or if the whole parsing should be interrupted)
- [ ] Reduce Int8/Int64 allocations (their parsers could be somewhat simplified?)
- [ ] Add combinator to parse whitespace (+ helper for multispace0/1?)
- [ ] Refactor TakeWhileOneOf to be "just" TakeWhile
- [ ] Refactor space to be of the form space0 and space1
- [ ] Rename `LF` to `Newline`
- [X] Document Recognize as explicitly as possible
- [X] Add Examples
- [x] Add Benchmarks
- [x] Make sure the Failure messages are properly cased
- [x] Rename `p` parser arguments to `parse` for clearer code
- [x] Add `Many0` and `Many1` parsers

## Maybes

- [ ] Rename project to `crayon`?
- [ ] Rename `Preceded` to `Prefixed`
- [ ] Rename `Terminated` to `Suffixed`
- [ ] Rename `Sequence` to `List`?
- [ ] Rename `Satisfy` to `Satisfies`?
- [X] Introduce `SeparatedList` as a result of previous?
- [X] Create `bytes.go` file to distinguish from characters

## Track

- [ ] Chase allocations, document them, and reduce their amount as much as possible

## NoNos
- [X] Add an `ErrInfiniteLoop` (`Many0`)