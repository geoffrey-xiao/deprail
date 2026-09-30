// Package process provides bounded and cancellable external-process execution.
//
// Captured stdout and stderr each retain at most Request.OutputCap bytes.
// Reaching either cap is reported as SCANNER_OUTPUT_LIMIT.
package process
