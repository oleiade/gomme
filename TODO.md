# TODOS

## Dos

- [ ] Create `Bytes` interface type for bytes file's content
- [ ] Create `String` or `Characters` interface type for characters file's content
- [ ] Add Examples
- [ ] Document Recognize as explicitly as possible
- [ ] Add an `ErrInfiniteLoop` (`Many0`)
- [ ] Sort Out Fatal/Non-Fatal errors (distinguish whether a parser failed in an expected manner, or if the whole parsing should be interrupted)
- [ ] Reduce Int8/Int64 allocations (their parsers could be somewhat simplified?)
- [x] Add Benchmarks
- [x] Make sure the Failure messages are properly cased
- [x] Rename `p` parser arguments to `parse` for clearer code
- [x] Add `Many0` and `Many1` parsers

## Maybes

- [ ] Rename project to `crayon`?
- [ ] Rename `Preceded` to `Prefixed`
- [ ] Rename `Terminated` to `Suffixed`
- [ ] Rename `Sequence` to `List`?
- [ ] Introduce `SeparatedList` as a result of previous?
- [ ] Rename `Satisfy` to `Satisfies`?

## Track

- [ ] Chase allocations, document them, and reduce their amount as much as possible
