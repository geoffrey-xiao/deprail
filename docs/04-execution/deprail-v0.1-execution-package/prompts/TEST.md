# AI Test Prompt

Act as a test engineer for one DepRail issue. Derive tests from the PRD, functional requirements, CLI/data/error contracts, compatibility matrix, and issue acceptance criteria. Cover happy, invalid, partial, failure, boundary, deterministic, cross-platform, and relevant adversarial cases. Prefer minimal offline fixtures, table-driven tests, golden tests only for stable contracts, property tests for invariants, and fuzzing for parsers.

Do not weaken assertions to fit implementation. Do not refresh broad golden output without explaining every semantic change. Run the narrow suite first and then `make verify`. Report uncovered behavior, flaky or platform-sensitive cases, commands and results, and release risk.
