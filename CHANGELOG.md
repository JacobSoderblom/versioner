# versioner

## 0.1.0

### New features

Conventional changesets

### Bug fixes

Fixing the order of titles in changelog generation and duplicates should not appear anymore.

### Documentation

Added installation through go in README

### Refactoring

Change how project values are used, instead of refetching values as repository and working dir, we now store them in a context struct

Change how project values are used, instead of refetching values as repository and working dir, we now store them in a context struct

Moving most of the cli logic to the command package, so the logic is not spread across the repository that much.
