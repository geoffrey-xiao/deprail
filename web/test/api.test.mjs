import assert from 'node:assert/strict';
import { afterEach, test } from 'node:test';
import { ApiFailure, get, setToken } from '../src/api.ts';

const originalFetch = globalThis.fetch;
const originalWindow = globalThis.window;
afterEach(() => {
  globalThis.fetch = originalFetch;
  globalThis.window = originalWindow;
  setToken(null);
});

const summary = {
  historyEntryID: '00000000-0000-4000-8000-000000000101',
  recordedAt: '2026-10-06T00:00:00Z',
  sourceScanID: null,
  operationOutcome: 'completed',
  reportStatus: 'complete',
  repositoryLabel: null,
  workspaceCount: null,
  findingCount: 0,
  sourceReportSchemaVersion: null,
};

async function readInjectedResponse(path, body) {
  globalThis.window = { location: { origin: 'http://127.0.0.1:8080' } };
  globalThis.fetch = async () => new Response(JSON.stringify(body), {
    status: 200, headers: { 'content-type': 'application/json' },
  });
  setToken('test-only-not-a-real-credential');
  return get(path);
}

test('contradictory source metadata cannot display a complete zero-finding history row', async () => {
  await assert.rejects(
    readInjectedResponse('/api/v1/scans', { items: [summary], nextCursor: null }),
    error => error instanceof ApiFailure && error.code === 'API_RESPONSE_INVALID',
  );
});

test('detail report presence must match the summary before rendering', async () => {
  await assert.rejects(
    readInjectedResponse(`/api/v1/scans/${summary.historyEntryID}`, {
      summary: { ...summary, sourceScanID: 'scan-1', sourceReportSchemaVersion: 'v1alpha' },
      report: null, diagnostics: [], artifactReferences: [],
    }),
    error => error instanceof ApiFailure && error.code === 'API_RESPONSE_INVALID',
  );
});
