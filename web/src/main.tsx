import React, { useEffect, useRef, useState } from 'react';
import { createRoot } from 'react-dom/client';
import { ApiFailure, get, hasToken, setToken, type Collection, type Diagnostic, type Finding, type HistoryDetail, type Page, type Summary, type Workspace } from './api';
import './style.css';

// Scrub the bootstrap credential synchronously, before mounting or making any request.
const initial = new URL(window.location.href);
const bootstrap = new URLSearchParams(initial.hash.slice(1)).get('token');
if (initial.hash) {
  initial.hash = '';
  window.history.replaceState(window.history.state, '', initial.pathname + initial.search);
}
if (bootstrap) setToken(bootstrap);

function routeFromPath(path: string) {
  if (path === '/console/' || path === '/console/history') return { page: 'history' as const, id: '' };
  if (path === '/console/about') return { page: 'about' as const, id: '' };
  const match = path.match(/^\/console\/scans\/([0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12})$/i);
  if (match) return { page: 'detail' as const, id: match[1] };
  return { page: 'missing' as const, id: '' };
}
function App() {
  const [route, setRoute] = useState(() => routeFromPath(location.pathname));
  const currentPath = useRef(location.pathname);
  const historyRun = useRef(0);
  const detailRun = useRef(0);
  const childRun = useRef({ workspaces: 0, findings: 0 });
  const [authEpoch, setAuthEpoch] = useState(0);
  const [history, setHistory] = useState<Page<Summary> | null>(null);
  const [historyError, setHistoryError] = useState<ApiFailure | null>(null);
  const [historyBusy, setHistoryBusy] = useState(false);
  const [cursorStack, setCursorStack] = useState<(string | undefined)[]>([undefined]);
  const [cursorIndex, setCursorIndex] = useState(0);
  const [detail, setDetail] = useState<HistoryDetail | null>(null);
  const [detailError, setDetailError] = useState<ApiFailure | null>(null);
  const [workspaces, setWorkspaces] = useState<Collection<Workspace> | null>(null);
  const [findings, setFindings] = useState<Collection<Finding> | null>(null);
  const [childErrors, setChildErrors] = useState<Partial<Record<'workspaces' | 'findings', ApiFailure>>>({});
  const [childBusy, setChildBusy] = useState<Record<string, boolean>>({});
  const [workspaceCursor, setWorkspaceCursor] = useState<string | undefined>();
  const [findingCursor, setFindingCursor] = useState<string | undefined>();
  const [notice, setNotice] = useState('');

  function clearDetail() {
    setDetail(null); setDetailError(null); setWorkspaces(null); setFindings(null);
    setChildErrors({}); setChildBusy({});
    setWorkspaceCursor(undefined); setFindingCursor(undefined);
  }

  useEffect(() => {
    const update = () => {
      const params = new URLSearchParams(location.hash.slice(1));
      if (!params.has('token')) return;
      const fragment = params.get('token');
      const current = new URL(window.location.href);
      current.hash = '';
      window.history.replaceState(window.history.state, '', current.pathname + current.search);
      historyRun.current++;
      detailRun.current++;
      setToken(fragment);
      setAuthEpoch(value => value + 1);
      setNotice(fragment ? 'Local session credentials updated.' : 'Authentication was lost. Reopen the console from the active local DepRail session.');
    };
    window.addEventListener('hashchange', update);
    const pop = () => {
      if (location.pathname === currentPath.current) return;
      currentPath.current = location.pathname;
      historyRun.current++;
      detailRun.current++;
      clearDetail();
      setRoute(routeFromPath(location.pathname));
      setNotice('Page changed.');
    };
    window.addEventListener('popstate', pop);
    return () => { window.removeEventListener('hashchange', update); window.removeEventListener('popstate', pop); };
  }, []);

  useEffect(() => {
    const focus = window.setTimeout(() => document.querySelector<HTMLElement>('main .error-panel h2, main h1, main h2')?.focus(), 0);
    return () => window.clearTimeout(focus);
  }, [route.page, route.id, detail, detailError, historyError, authEpoch]);

  async function loadHistory(cursor = cursorStack[cursorIndex]) {
    const run = ++historyRun.current;
    setHistoryBusy(true); setHistoryError(null);
    try {
      const result = await get<Page<Summary>>('/api/v1/scans', cursor);
      if (run === historyRun.current) setHistory(result);
    } catch (error) {
      if (run === historyRun.current) setHistoryError(error instanceof ApiFailure ? error : new ApiFailure('API_INTERNAL_ERROR', 0, 'History could not be loaded.'));
    } finally {
      if (run === historyRun.current) setHistoryBusy(false);
    }
  }
  useEffect(() => {
    if (route.page === 'history' && hasToken()) void loadHistory();
    else if (route.page === 'history') { setHistory(null); setHistoryError(null); }
  }, [route.page, authEpoch, cursorIndex]);

  async function loadDetail() {
    const run = ++detailRun.current;
    if (!hasToken()) return;
    const id = route.id;
    clearDetail();
    try {
      const result = await get<HistoryDetail>(`/api/v1/scans/${encodeURIComponent(id)}`);
      if (run !== detailRun.current) return;
      setDetail(result);
      void loadChild('workspaces', undefined, run, id);
      void loadChild('findings', undefined, run, id);
    } catch (error) {
      if (run === detailRun.current) setDetailError(error instanceof ApiFailure ? error : new ApiFailure('API_INTERNAL_ERROR', 0, 'The history entry could not be loaded.'));
    }
  }
  async function loadChild(kind: 'workspaces' | 'findings', cursor?: string, run = detailRun.current, id = route.id) {
    if (!hasToken() || run !== detailRun.current) return;
    const child = ++childRun.current[kind];
    setChildBusy(previous => ({ ...previous, [kind]: true }));
    setChildErrors(previous => { const remaining = { ...previous }; delete remaining[kind]; return remaining; });
    try {
      if (kind === 'workspaces') {
        const result = await get<Collection<Workspace>>(`/api/v1/scans/${encodeURIComponent(id)}/workspaces`, cursor);
        if (run === detailRun.current && child === childRun.current[kind]) setWorkspaces(result);
      } else {
        const result = await get<Collection<Finding>>(`/api/v1/scans/${encodeURIComponent(id)}/findings`, cursor);
        if (run === detailRun.current && child === childRun.current[kind]) setFindings(result);
      }
    } catch (error) {
      if (run === detailRun.current && child === childRun.current[kind]) setChildErrors(previous => ({ ...previous, [kind]: error instanceof ApiFailure ? error : new ApiFailure('API_INTERNAL_ERROR', 0, 'This collection could not be loaded.') }));
    } finally {
      if (run === detailRun.current && child === childRun.current[kind]) setChildBusy(previous => ({ ...previous, [kind]: false }));
    }
  }
  useEffect(() => { if (route.page === 'detail') void loadDetail(); }, [route.page, route.id, authEpoch]);

  function navigate(event: React.MouseEvent<HTMLAnchorElement>, path: string) {
    if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey || event.button !== 0) return;
    event.preventDefault();
    if (location.pathname === path) return;
    currentPath.current = path;
    historyRun.current++;
    detailRun.current++;
    clearDetail();
    window.history.pushState(null, '', path);
    setRoute(routeFromPath(path));
    setNotice('Page changed.');
  }
  function nextHistory() {
    if (!history?.nextCursor) return;
    const next = [...cursorStack.slice(0, cursorIndex + 1), history.nextCursor];
    setCursorStack(next); setCursorIndex(next.length - 1);
  }
  function previousHistory() { if (cursorIndex > 0) setCursorIndex(cursorIndex - 1); }
  const recovery = !hasToken();

  return <>
    <a className="skip-link" href="#main">Skip to main content</a>
    <header className="topbar"><div className="topbar-inner">
      <a className="brand" href="/console/history" onClick={event => navigate(event, '/console/history')} aria-label="DepRail local history home"><span className="brand-mark" aria-hidden="true">D</span><span>DepRail <span className="brand-sub">Local history</span></span></a>
      <div className="header-actions"><span className="local-badge"><span aria-hidden="true">●</span> Local only</span><nav aria-label="Main navigation"><a href="/console/history" aria-current={route.page === 'history' ? 'page' : undefined} onClick={event => navigate(event, '/console/history')}>History</a><a href="/console/about" aria-current={route.page === 'about' ? 'page' : undefined} onClick={event => navigate(event, '/console/about')}>About</a></nav></div>
    </div></header>
    <div className="announcement" role="status" aria-live="polite" aria-atomic="true">{notice}</div>
    <main id="main" className="container">
      {recovery ? <section className="panel recovery" aria-labelledby="recovery-title"><p className="eyebrow">Local session required</p><h1 id="recovery-title" tabIndex={-1}>Reopen the local console</h1><p>This tab does not have an active in-memory session. For your safety, credentials are not saved across reloads or tabs.</p><p>Return to the active DepRail terminal and press Enter to reopen the console. For a non-interactive session, restart <code>deprail web --open</code>.</p></section> : route.page === 'history' ? <>
        <div className="page-title"><div><p className="eyebrow">Saved operations</p><h1 tabIndex={-1}>Scan history</h1><p className="lede">Read-only records from this local DepRail session.</p></div></div>
        <div className="local-notice"><strong>History stays local</strong><span>Saved records on this device. This console does not start scans or change results.</span></div>
        {historyBusy && !history ? <div className="panel loading" role="status">Loading saved scan history…</div> : historyError ? <ErrorPanel error={historyError} retry={() => void loadHistory()} /> : history?.items.length === 0 ? <div className="panel empty"><h2>No saved scans</h2><p>No scan history has been recorded for this local session.</p><p>Run a scan with history capture enabled, then refresh this page.</p><button type="button" onClick={() => void loadHistory()} disabled={historyBusy}>Refresh history</button></div> : history ? <>
          <section className="panel history-panel" aria-label="Saved scans">
            <div className="panel-heading"><div><h2>Recent scans</h2><p>Newest first · operation and report status are shown separately.</p></div><button type="button" className="quiet" onClick={() => void loadHistory()} disabled={historyBusy}>Refresh history</button></div>
            <div className="history-columns" aria-hidden="true"><span>Repository / history entry</span><span>Operation</span><span>Report</span><span>Findings</span><span>Recorded</span></div>
            <ul className="history-list">{history.items.map(item => <li key={item.historyEntryID}><a className="history-item" href={`/console/scans/${encodeURIComponent(item.historyEntryID)}`} onClick={event => navigate(event, `/console/scans/${encodeURIComponent(item.historyEntryID)}`)}>
              <div className="item-main"><div className="item-title"><strong>{item.repositoryLabel ?? 'Repository label unavailable'}</strong><span aria-hidden="true">→</span></div><p className="item-meta">Entry <code>{item.historyEntryID}</code></p><p className="item-source">Source scan: {item.sourceScanID ? <code>{item.sourceScanID}</code> : 'Unavailable'}{item.sourceReportSchemaVersion && <> · Schema {item.sourceReportSchemaVersion}</>}</p></div>
              <div className="history-cell"><span className="cell-label">Operation outcome</span><Status value={item.operationOutcome} /></div>
              <div className="history-cell"><span className="cell-label">Report completeness</span>{item.reportStatus ? <Status value={item.reportStatus} /> : <span className="status status-neutral">Unavailable</span>}</div>
              <div className="history-cell"><span className="cell-label">Findings</span><strong className="finding-count">{item.findingCount === null ? 'Unavailable' : item.findingCount}</strong><span className="cell-note">Workspaces: {item.workspaceCount === null ? 'Unavailable' : item.workspaceCount}</span></div>
              <div className="history-cell"><span className="cell-label">Recorded</span><time dateTime={item.recordedAt}>{formatDate(item.recordedAt)}</time></div>
            </a></li>)}</ul>
          </section>
          <nav className="pagination" aria-label="History pages"><div className="pagination-controls"><button type="button" onClick={previousHistory} disabled={cursorIndex === 0 || historyBusy}>Previous page</button><span className="page-number" aria-live="polite">Page {cursorIndex + 1}</span><button type="button" onClick={nextHistory} disabled={!history.nextCursor || historyBusy}>Next page</button></div><span className="page-note">{history.items.length} entries on this page</span></nav>
        </> : null}
      </> : route.page === 'detail' ? <>
        <p className="breadcrumb"><a href="/console/history" onClick={event => navigate(event, '/console/history')}>← Back to scan history</a></p>
        {detailError ? <ErrorPanel error={detailError} retry={() => void loadDetail()} /> : !detail ? <div className="panel loading" role="status">Loading scan details…</div> : <Detail detail={detail} workspaces={workspaces} findings={findings} workspaceError={childErrors.workspaces} findingError={childErrors.findings} workspaceBusy={childBusy.workspaces} findingBusy={childBusy.findings} retryDetail={() => void loadDetail()} retryWorkspace={(cursor) => void loadChild('workspaces', cursor)} retryFinding={(cursor) => void loadChild('findings', cursor)} workspaceCursor={workspaceCursor} findingCursor={findingCursor} setWorkspaceCursor={setWorkspaceCursor} setFindingCursor={setFindingCursor} />}
      </> : route.page === 'about' ? <section className="panel about"><p className="eyebrow">About this console</p><h1 tabIndex={-1}>Local history, read only</h1><p>DepRail local history lets you review saved scan-operation records from the active local process. It does not start scans or change saved results.</p><h2>Privacy and session</h2><p>History is provided by the DepRail process running on this device. This page does not save your session credential in browser storage. Reloading or opening a new tab requires reopening the console from the active terminal.</p><h2>Reading a record</h2><p>Operation outcome describes whether the scan operation completed, failed, or was cancelled. Report status describes the completeness of any captured report; they are shown independently.</p><p>Diagnostics, workspaces, findings and artifact integrity are displayed only when provided by the local API. Unavailable information is not treated as an empty or clean result.</p></section> : <section className="panel"><p className="eyebrow">Page not found</p><h1 tabIndex={-1}>This console page is unavailable</h1><p>Use the history page to choose a saved entry.</p><a href="/console/history" onClick={event => navigate(event, '/console/history')}>Go to scan history</a></section>}
    </main>
    <footer className="footer"><div className="container"><span>DepRail local console</span><a href="/console/about" onClick={event => navigate(event, '/console/about')}>About local history</a></div></footer>
  </>;
}
function label(value: string) { return value.charAt(0).toUpperCase() + value.slice(1); }
function formatDate(value: string) { const date = new Date(value); return Number.isNaN(date.getTime()) ? value : new Intl.DateTimeFormat(undefined, { dateStyle: 'medium', timeStyle: 'short' }).format(date); }
function Status({ value }: { value: string }) { return <span className={`status status-${value}`}><span className="status-dot" aria-hidden="true" />{label(value)}</span>; }
function ErrorPanel({ error, retry }: { error: ApiFailure; retry: () => void }) {
  const compatibility = ['HISTORY_SCHEMA_UNSUPPORTED', 'HISTORY_MIGRATION_FAILED', 'HISTORY_CORRUPT'].includes(error.code);
  const notFound = error.code === 'HISTORY_ENTRY_NOT_FOUND';
  return <section className="panel error-panel" role="alert"><p className="eyebrow">{compatibility ? 'History compatibility problem' : notFound ? 'Entry unavailable' : 'Unable to load data'}</p><h2 tabIndex={-1}>{compatibility ? 'History cannot be read safely' : notFound ? 'This scan entry was not found' : 'The request did not succeed'}</h2><p>{error.message}</p><p className="error-code">Error code: <code>{error.code}</code></p>{compatibility && <p>Existing history has not been changed. Restarting or deleting local data is not recommended.</p>}{error.status !== 401 && <button type="button" onClick={retry}>Retry this read</button>}</section>;
}
function CollectionState<T>({ title, data, error, busy, retry, children, unavailable, emptyContent }: { title: string; data: Collection<T> | null; error?: ApiFailure; busy?: boolean; retry: () => void; children: React.ReactNode; unavailable: string; emptyContent?: React.ReactNode }) {
  return <section className="panel section-panel"><h2>{title}</h2>{busy && !data ? <p role="status">Loading {title.toLowerCase()}…</p> : error ? <div className="inline-error" role="alert"><p>{error.message} <code>{error.code}</code></p><button type="button" onClick={retry}>Retry {title.toLowerCase()} read</button></div> : !data ? <p role="status">Not loaded.</p> : data.collectionState === 'unavailable' ? <p className="callout">{unavailable}</p> : data.items === null ? <p className="callout">This collection returned no usable data. It is not treated as empty.</p> : data.items.length === 0 ? emptyContent ?? <p className="muted">No {title.toLowerCase()} were recorded in this collection.</p> : children}</section>;
}
function Detail({ detail, workspaces, findings, workspaceError, findingError, workspaceBusy, findingBusy, retryDetail, retryWorkspace, retryFinding, workspaceCursor, findingCursor, setWorkspaceCursor, setFindingCursor }: { detail: HistoryDetail; workspaces: Collection<Workspace> | null; findings: Collection<Finding> | null; workspaceError?: ApiFailure; findingError?: ApiFailure; workspaceBusy?: boolean; findingBusy?: boolean; retryDetail: () => void; retryWorkspace: (cursor?: string) => void; retryFinding: (cursor?: string) => void; workspaceCursor?: string; findingCursor?: string; setWorkspaceCursor: (cursor?: string) => void; setFindingCursor: (cursor?: string) => void }) {
  const summary = detail.summary;
  const reportStatus = summary.reportStatus;

  const workspacePage = workspaces;
  const findingPage = findings;
  return <>
    <div className="page-title detail-title"><div><p className="eyebrow">Scan record</p><h1 tabIndex={-1}>{summary.repositoryLabel ?? 'Repository label unavailable'}</h1><p className="lede">Recorded <time dateTime={summary.recordedAt}>{formatDate(summary.recordedAt)}</time></p></div><button className="quiet" type="button" onClick={retryDetail}>Refresh detail</button></div>
    <section className="outcome-panel" aria-label="Operation and report status"><div className="outcome-card"><p className="eyebrow">Operation outcome</p><Status value={summary.operationOutcome} /><p>Execution outcome, independent of report completeness.</p></div><div className="outcome-card"><p className="eyebrow">Report completeness</p>{reportStatus ? <Status value={reportStatus} /> : <span className="status status-neutral">Report unavailable</span>}<p>Status from the captured report, not the operation outcome.</p></div></section>
    {reportStatus === 'partial' && <p className="callout warning" role="status">This report is partial. Some scope may be omitted or incomplete; review diagnostics below.</p>}
    {reportStatus === 'failed' && <p className="callout warning" role="status">The report failed. Findings are not presented as a clean result.</p>}
    {!detail.report && <p className="callout">No trustworthy report was captured. Findings are unavailable, not empty.</p>}
    <dl className="metric-panel" aria-label="Saved report summary"><div><dt>Findings</dt><dd>{summary.findingCount === null ? 'Unavailable' : summary.findingCount}</dd></div><div><dt>Workspaces</dt><dd>{summary.workspaceCount === null ? 'Unavailable' : summary.workspaceCount}</dd></div><div><dt>Recorded</dt><dd className="metric-date"><time dateTime={summary.recordedAt}>{formatDate(summary.recordedAt)}</time></dd></div></dl>
    <div className="detail-grid"><div className="detail-primary">
    <CollectionState title="Findings" data={findingPage} error={findingError} busy={findingBusy} unavailable="Findings are unavailable because no trustworthy source report exists." retry={() => retryFinding()} emptyContent={findingPage?.reportStatus === 'complete' && reportStatus === 'complete' ? <p className="muted">No findings are recorded for this complete report.</p> : <p className="callout warning">This findings page is empty, but report completeness is {findingPage?.reportStatus ?? 'unavailable'} and the entry report status is {reportStatus ?? 'unavailable'}; it is not a clean zero-finding result.</p>}>
      {findingPage?.reportStatus && findingPage.reportStatus !== reportStatus && <p className="callout warning">Findings collection report status: {label(findingPage.reportStatus)}.</p>}
      <ul className="record-list findings">{(findingPage?.items ?? []).map(item => <li key={item.stableFindingKey}>
        <div className="finding-heading"><strong>{item.componentName} <span className="muted">{item.version}</span></strong><span className="severity">{item.severity ? label(item.severity) : 'Severity unavailable'}</span></div>
        <p>Vulnerability: <code>{item.vulnerabilityID}</code></p>
        <p className="finding-context">Workspace: {item.relativePath} · {item.ecosystem}</p>
        {item.fixedVersion && <p className="finding-context">Fixed version: <strong>{item.fixedVersion}</strong></p>}
        <details className="finding-evidence"><summary>Identifiers and evidence</summary><dl>
          {item.componentPURL && <div><dt>Package URL</dt><dd><code className="wrap-text">{item.componentPURL}</code></dd></div>}
          {item.aliases.length > 0 && <div><dt>Aliases</dt><dd>{item.aliases.map((alias, index) => <React.Fragment key={`${alias}-${index}`}>{index > 0 && ', '}<code>{alias}</code></React.Fragment>)}</dd></div>}
          <div><dt>Stable finding key</dt><dd><code className="digest">{item.stableFindingKey}</code></dd></div>
        </dl></details>
      </li>)}</ul>
      {(findingPage?.nextCursor || findingCursor) && <div className="collection-actions">{findingPage?.nextCursor && <button type="button" onClick={() => { setFindingCursor(findingPage.nextCursor ?? undefined); retryFinding(findingPage.nextCursor ?? undefined); }}>Next finding page</button>}{findingCursor && <button type="button" className="quiet" onClick={() => { setFindingCursor(undefined); retryFinding(); }}>First finding page</button>}</div>}
    </CollectionState>
    <section className="panel section-panel"><h2>Diagnostics</h2>{detail.diagnostics.length === 0 ? <p className="muted">No diagnostics were recorded.</p> : <ul className="record-list diagnostics">{detail.diagnostics.map((diagnostic: Diagnostic, index) => <li key={`${diagnostic.code}-${index}`}><div className="finding-heading"><strong>{diagnostic.code}</strong><span>{label(diagnostic.scope)}{diagnostic.workspaceID ? ` · ${diagnostic.workspaceID}` : ''}</span></div><p>{diagnostic.message}</p></li>)}</ul>}</section>
    </div><aside className="detail-sidebar" aria-label="Supporting scan information">
    <CollectionState title="Workspaces" data={workspacePage} error={workspaceError} busy={workspaceBusy} unavailable="Workspace data is unavailable because validated workspace context was not captured with this operation." retry={() => retryWorkspace()}>
      <ul className="record-list">{(workspacePage?.items ?? []).map(item => <li key={item.workspaceID}><strong>{item.relativePath}</strong><span>{item.ecosystem ?? 'Ecosystem unavailable'} · {item.packageManager ?? 'Package manager unavailable'}</span><span>Discovery: {item.discoveryCompleteness ?? 'Unavailable'} (not scan completeness)</span></li>)}</ul>
      {(workspacePage?.nextCursor || workspaceCursor) && <div className="collection-actions">{workspacePage?.nextCursor && <button type="button" onClick={() => { setWorkspaceCursor(workspacePage.nextCursor ?? undefined); retryWorkspace(workspacePage.nextCursor ?? undefined); }}>Next workspace page</button>}{workspaceCursor && <button type="button" className="quiet" onClick={() => { setWorkspaceCursor(undefined); retryWorkspace(); }}>First workspace page</button>}</div>}
    </CollectionState>
    <section className="panel metadata"><h2>Record information</h2><dl><div><dt>History entry</dt><dd><code>{summary.historyEntryID}</code></dd></div><div><dt>Source scan</dt><dd>{summary.sourceScanID ? <code>{summary.sourceScanID}</code> : 'Unavailable'}</dd></div><div><dt>Recorded</dt><dd><time dateTime={summary.recordedAt}>{formatDate(summary.recordedAt)}</time></dd></div><div><dt>Report schema</dt><dd>{summary.sourceReportSchemaVersion ?? 'Unavailable'}</dd></div><div><dt>Workspace count</dt><dd>{summary.workspaceCount ?? 'Unavailable'}</dd></div></dl></section>
    {detail.report && <section className="panel metadata"><h2>Provenance</h2><dl><div><dt>Source schema version</dt><dd>{detail.report.sourceSchemaVersion}</dd></div><div><dt>Repository state digest</dt><dd><code className="digest">{detail.report.repositoryState}</code></dd></div>{Object.entries(detail.report.provenance).map(([key, value]) => <div key={key}><dt>{label(key.replace(/[A-Z]/g, letter => ` ${letter.toLowerCase()}`))}</dt><dd>{value ?? 'Not recorded'}</dd></div>)}</dl></section>}
    <section className="panel section-panel"><h2>Artifact integrity</h2>{detail.artifactReferences.length === 0 ? <p className="muted">No artifact references were recorded.</p> : <ul className="record-list">{detail.artifactReferences.map(item => <li key={item.digest}><span className="digest-label">SHA-256</span><code className="digest">{item.digest}</code><Status value={item.integrity} />{item.integrity !== 'verified' && <p className="callout warning">Artifact evidence is not verified and is not available as trusted raw evidence.</p>}</li>)}</ul>}</section>
    </aside></div>
  </>;
}

const root = document.getElementById('root');
if (!root) throw new Error('Missing application root');
createRoot(root).render(<App />);
