export type Outcome = 'completed' | 'failed' | 'cancelled';
export type ReportStatus = 'complete' | 'partial' | 'failed';
export interface Summary {
  historyEntryID: string; recordedAt: string; sourceScanID: string | null;
  operationOutcome: Outcome; reportStatus: ReportStatus | null; repositoryLabel: string | null;
  workspaceCount: number | null; findingCount: number | null; sourceReportSchemaVersion: string | null;
}
export interface Page<T> { items: T[]; nextCursor: string | null }
export interface HistoryDetail {
  summary: Summary;
  report: null | { sourceSchemaVersion: string; repositoryState: string; provenance: { deprailVersion: string | null; scannerName: string | null; scannerVersion: string | null; scannerDatabaseVersion: string | null } };
  diagnostics: Diagnostic[];
  artifactReferences: { digest: string; integrity: 'verified' | 'missing' | 'digest_mismatch' | 'unavailable' }[];
}
export interface Diagnostic { code: string; message: string; scope: 'repository' | 'workspace'; workspaceID?: string | null }
export interface Workspace { workspaceID: string; relativePath: string; ecosystem: string; packageManager: string; discoveryCompleteness: 'complete' | 'partial' | 'failed' }
export interface Finding { stableFindingKey: string; componentName: string; componentPURL: string | null; version: string; workspaceID: string; relativePath: string; ecosystem: string; vulnerabilityID: string; aliases: string[]; severity: 'unknown' | 'low' | 'moderate' | 'high' | 'critical' | null; fixedVersion: string | null }
export interface Collection<T> { collectionState: 'available' | 'unavailable'; items: T[] | null; nextCursor: string | null; reason?: string; reportStatus?: ReportStatus | null }
export class ApiFailure extends Error {
  constructor(readonly code: string, readonly status: number, message: string) { super(message); this.name = 'ApiFailure'; }
}
let bearer: string | null = null;
export function setToken(value: string | null) { bearer = value; }
export function hasToken() { return bearer !== null; }
export async function get<T>(path: string, cursor?: string): Promise<T> {
  if (!bearer) throw new ApiFailure('API_AUTH_UNAUTHORIZED', 401, 'Authentication is unavailable. Reopen the console from the active local DepRail session.');
  const url = new URL(path, window.location.origin);
  if (cursor) url.searchParams.set('cursor', cursor);
  try {
    const response = await fetch(url, { headers: { Authorization: `Bearer ${bearer}`, Accept: 'application/json' }, cache: 'no-store', credentials: 'omit', referrerPolicy: 'no-referrer', redirect: 'error' });
    let body: unknown;
    try { body = await response.json(); } catch { body = null; }
    if (!response.ok) {
      const err = body as { error?: { code?: string; message?: string } } | null;
      const code = err?.error?.code ?? (response.status === 414 ? 'API_REQUEST_TOO_LONG' : response.status === 431 ? 'API_HEADERS_TOO_LARGE' : 'API_INTERNAL_ERROR');
      if (response.status === 401) bearer = null;
      throw new ApiFailure(code, response.status, err?.error?.message ?? httpMessage(response.status));
    }
    if (!validResponse(path, body)) throw new ApiFailure('API_RESPONSE_INVALID', response.status, 'The local API returned an invalid response. No data from this read is treated as a successful result.');
    return body as T;
  } catch (error) {
    if (error instanceof ApiFailure) throw error;
    throw new ApiFailure('API_CONNECTION_FAILED', 0, 'The local API could not be reached. Check that the DepRail console is still running, then retry.');
  }
}
function httpMessage(status: number): string {
  if (status === 401) return 'Your local session is no longer authenticated. Reopen the console from the active local DepRail session.';
  if (status === 404) return 'This history entry is no longer available.';
  if (status === 414 || status === 431) return 'The request exceeded a local transport limit.';
  if (status === 503) return 'The local history service is unavailable or busy. You can retry this read.';
  if (status === 504) return 'The local read timed out. You can retry it.';
  return 'The local API could not complete this read. Other loaded information is unchanged.';
}

type RecordValue = Record<string, unknown>;
function record(value: unknown): value is RecordValue { return value !== null && typeof value === 'object' && !Array.isArray(value); }
function nullableString(value: unknown): boolean { return value === null || typeof value === 'string'; }
function nullableNumber(value: unknown): boolean { return value === null || typeof value === 'number' && Number.isFinite(value) && value >= 0; }
function status(value: unknown): boolean { return value === 'complete' || value === 'partial' || value === 'failed'; }
function summary(value: unknown): boolean {
  if (!record(value)) return false;
  if (typeof value.historyEntryID !== 'string' || typeof value.recordedAt !== 'string' ||
    !['completed', 'failed', 'cancelled'].includes(String(value.operationOutcome)) ||
    !(value.reportStatus === null || status(value.reportStatus)) ||
    !nullableString(value.repositoryLabel) || !nullableString(value.sourceScanID) ||
    !nullableString(value.sourceReportSchemaVersion) ||
    !nullableNumber(value.workspaceCount) || !nullableNumber(value.findingCount)) return false;
  return value.sourceScanID === null
    ? value.reportStatus === null && value.findingCount === null && value.sourceReportSchemaVersion === null
    : value.reportStatus !== null && value.findingCount !== null && value.sourceReportSchemaVersion !== null;
}
function workspace(value: unknown): boolean {
  return record(value) && typeof value.workspaceID === 'string' && typeof value.relativePath === 'string' &&
    typeof value.ecosystem === 'string' && typeof value.packageManager === 'string' &&
    status(value.discoveryCompleteness);
}
function finding(value: unknown): boolean {
  return record(value) && ['stableFindingKey', 'componentName', 'version', 'workspaceID', 'relativePath', 'ecosystem', 'vulnerabilityID'].every(key => typeof value[key] === 'string') &&
    nullableString(value.componentPURL) && nullableString(value.fixedVersion) &&
    (value.severity === null || ['unknown', 'low', 'moderate', 'high', 'critical'].includes(String(value.severity))) &&
    Array.isArray(value.aliases) && value.aliases.every((alias: unknown) => typeof alias === 'string');
}
function collection(value: unknown, item: (value: unknown) => boolean): boolean {
  if (!record(value) || !nullableString(value.nextCursor)) return false;
  if (value.collectionState === 'unavailable') return value.items === null && value.nextCursor === null;
  return value.collectionState === 'available' && Array.isArray(value.items) && value.items.every(item) &&
    (value.reportStatus === undefined || value.reportStatus === null || status(value.reportStatus));
}
function detail(value: unknown): boolean {
  if (!record(value) || !summary(value.summary) || !Array.isArray(value.diagnostics) || !Array.isArray(value.artifactReferences)) return false;
  if (!value.diagnostics.every((entry: unknown) => record(entry) && typeof entry.code === 'string' &&
    typeof entry.message === 'string' && ['repository', 'workspace'].includes(String(entry.scope)))) return false;
  if (!value.artifactReferences.every((entry: unknown) => record(entry) && typeof entry.digest === 'string' &&
    ['verified', 'missing', 'digest_mismatch', 'unavailable'].includes(String(entry.integrity)))) return false;
  const entry = value.summary as RecordValue;
  if (value.report === null) return entry.sourceScanID === null;
  if (entry.sourceScanID === null || !record(value.report) || typeof value.report.sourceSchemaVersion !== 'string' ||
    value.report.sourceSchemaVersion !== entry.sourceReportSchemaVersion ||
    typeof value.report.repositoryState !== 'string' || !record(value.report.provenance)) return false;
  const provenance = value.report.provenance;
  return ['deprailVersion', 'scannerName', 'scannerVersion', 'scannerDatabaseVersion'].every(key => nullableString(provenance[key]));
}
function validResponse(path: string, value: unknown): boolean {
  if (path === '/api/v1/scans') return record(value) && Array.isArray(value.items) && value.items.every(summary) && nullableString(value.nextCursor);
  if (path.endsWith('/workspaces')) return collection(value, workspace);
  if (path.endsWith('/findings')) return collection(value, finding);
  return detail(value);
}
