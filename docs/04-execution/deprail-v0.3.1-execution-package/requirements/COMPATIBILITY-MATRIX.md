# v0.3.1 Compatibility Matrix

| Surface | Required behavior | Evidence |
| --- | --- | --- |
| Interactive macOS/Linux/Windows TTY | Useful progress, safe styling, final summary | Manual smoke capture |
| Non-TTY redirect | Stable lines, no animation or ANSI | stdout/stderr capture |
| JSON stdout | Parseable machine data only | Contract test |
| `--quiet` | No successful progress/summary; required errors remain | CLI smoke |
| `--verbose` | Additional safe diagnostics only | CLI smoke |
| Complete zero findings | Explicit successful empty summary | Golden/contract test |
| Partial scan | Explicit incomplete summary | Failure fixture |
| Failed scan | Explicit failure and stable code | Failure fixture |
| Narrow terminal | Readable fallback layout | Width-controlled smoke |
| Monochrome terminal | Meaning independent of color | Color-disabled smoke |
| Cancellation | No corrupted terminal state | Controlled process test |
| Hostile labels | No control-sequence injection | Adversarial presenter test |

Serialized domain meaning and exit semantics must remain equivalent across supported platforms.
