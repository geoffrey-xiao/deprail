# v0.3.1 Error Model

Existing error codes and exit meanings remain authoritative. Presentation adds clarity; it does not reinterpret errors.

| State | Human requirement | Machine requirement |
| --- | --- | --- |
| Complete | Concise result summary; zero findings may be called no known vulnerabilities. | Preserve complete status and findings. |
| Partial | Explicitly say incomplete and identify safe affected scope. | Preserve diagnostics and partial status. |
| Failed | Explicitly say failed and show stable code/guidance. | Preserve failure code and exit behavior. |
| Cancelled | Say cancelled and restore terminal state. | Preserve cancellation/error semantics. |

Missing tools, timeouts, malformed output, output limits, artifact failures, and path failures must remain errors. They must never render as successful empty scans.

Untrusted error details are escaped and redacted before terminal output.
