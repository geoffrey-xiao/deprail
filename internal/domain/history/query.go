package history

// Entry is one validated history occurrence and its exact artifact references.
type Entry struct {
	Projection   Projection
	ArtifactRefs []string
}

// Cursor is a stable keyset boundary for history entries ordered newest first.
type Cursor struct {
	RecordedAtUS   int64
	HistoryEntryID string
}

// Page bounds a history list read. A zero limit selects the default page size.
type Page struct {
	Limit  int
	Before *Cursor
}

// Summary contains only validated fields used to list history entries.
type Summary struct {
	HistoryEntryID            string
	RecordedAtUS              int64
	OperationOutcome          string
	RepositoryLabel           *string
	WorkspaceCount            *int
	ReportStatus              *string
	FindingCount              *int
	SourceScanID              *string
	SourceReportSchemaVersion *string
}

// PageResult reports whether the store contains an additional row after Entries.
type PageResult struct {
	Entries []Summary
	HasMore bool
}
