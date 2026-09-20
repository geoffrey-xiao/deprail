# v0.3.1 Traceability Matrix

| Requirement | Planned evidence | Review |
| --- | --- | --- |
| UX-001/UX-002 shared presentation boundary | Event and renderer contract tests | Architecture review |
| UX-003/UX-004 TTY and non-TTY behavior | Stream capture and manual terminal evidence | CLI review |
| UX-005 truthful outcomes | Complete/partial/failed/cancelled fixtures | Safety review |
| UX-006 flag compatibility | Quiet/verbose CLI cases | Compatibility review |
| UX-007 JSON purity | Parse and stream-separation tests | Contract review |
| UX-008 hostile-label safety | ANSI/control-character adversarial cases | Security review |
| UX-009 narrow/monochrome behavior | Width- and color-controlled captures | UX review |
| UX-010 cross-platform equivalence | Linux/macOS/Windows smoke matrix | Release review |
| UX-011 unchanged scanner/repository behavior | Invocation audit and before/after tree checks | Architecture/security review |
| UX-012 cancellation cleanup | Controlled cancellation integration test | Reliability review |

Implementation issues, when authorized, must map to one or more rows or explain why they are documentation/process-only. No issue files are created by this planning update.
