package history

// SchemaVersion identifies the persisted projection contract.
const SchemaVersion = "history-v1"

// Projection is the explicit, privacy-safe history representation.
type Projection struct {
	SchemaVersion    string       `json:"schema_version"`
	HistoryEntryID   string       `json:"history_entry_id"`
	RecordedAtUS     int64        `json:"recorded_at_us"`
	OperationOutcome string       `json:"operation_outcome"`
	RepositoryLabel  *string      `json:"repository_label"`
	Workspaces       []Workspace  `json:"workspaces"`
	Report           *Report      `json:"report"`
	Diagnostics      []Diagnostic `json:"diagnostics"`
}
type Workspace struct {
	WorkspaceID           string `json:"workspace_id"`
	Path                  string `json:"path"`
	Ecosystem             string `json:"ecosystem"`
	PackageManager        string `json:"package_manager"`
	DiscoveryCompleteness string `json:"discovery_completeness"`
}
type Report struct {
	SourceSchemaVersion string     `json:"source_schema_version"`
	SourceScanID        string     `json:"source_scan_id"`
	Status              string     `json:"status"`
	RepositoryState     string     `json:"repository_state"`
	Provenance          Provenance `json:"provenance"`
	Findings            []Finding  `json:"findings"`
	ArtifactDigests     []string   `json:"artifact_digests"`
}
type Provenance struct {
	DeprailVersion         *string `json:"deprail_version"`
	ScannerName            *string `json:"scanner_name"`
	ScannerVersion         *string `json:"scanner_version"`
	ScannerDatabaseVersion *string `json:"scanner_database_version"`
}
type Finding struct {
	StableFindingKey string   `json:"stable_finding_key"`
	ComponentName    string   `json:"component_name"`
	ComponentPURL    *string  `json:"component_purl"`
	Version          string   `json:"version"`
	WorkspaceID      string   `json:"workspace_id"`
	RelativePath     string   `json:"relative_path"`
	Ecosystem        string   `json:"ecosystem"`
	VulnerabilityID  string   `json:"vulnerability_id"`
	Aliases          []string `json:"aliases"`
	Severity         *string  `json:"severity"`
	FixedVersion     *string  `json:"fixed_version"`
}
type Diagnostic struct {
	Code        string  `json:"code"`
	Scope       string  `json:"scope"`
	Message     string  `json:"message"`
	WorkspaceID *string `json:"workspace_id"`
}
