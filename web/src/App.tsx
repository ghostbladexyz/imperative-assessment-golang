import {
  type PointerEvent as ReactPointerEvent,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import CodeMirror from "@uiw/react-codemirror";
import { go } from "@codemirror/lang-go";
import {
  ArrowLeft,
  ArrowRight,
  Check,
  Copy,
  Download,
  Eye,
  EyeOff,
  ListChecks,
  Maximize2,
  Minimize2,
  Play,
  RotateCcw,
  Sparkles,
  Square,
  Terminal,
} from "lucide-react";
import {
  fetchCatalogue,
  formatCode,
  runTests,
  validateReceipts,
} from "./api";
import { buildConsoleStreams } from "./console";
import { loadProgress, saveProgress } from "./storage";
import type { Level, RunResult, SavedProgress } from "./types";

type PaneSizes = {
  brief: number;
  output: number;
  tests: number;
};

type ResizeKind = keyof PaneSizes;

const PANE_SIZE_KEY = "imperative-go-assessment:pane-sizes:v1";
const TESTS_VISIBILITY_KEY = "imperative-go-assessment:tests-visible:v1";
const DEFAULT_PANE_SIZES: PaneSizes = {
  brief: 340,
  output: 250,
  tests: 480,
};

function App() {
  const [levels, setLevels] = useState<Level[]>([]);
  const [progress, setProgress] = useState<SavedProgress | null>(null);
  const [runs, setRuns] = useState<Record<string, RunResult>>({});
  const [loadingError, setLoadingError] = useState("");
  const [running, setRunning] = useState(false);
  const [saved, setSaved] = useState(true);
  const [editRevision, setEditRevision] = useState(0);
  const [fullscreen, setFullscreen] = useState(false);
  const [paneSizes, setPaneSizes] = useState(loadPaneSizes);
  const [testsVisible, setTestsVisible] = useState(loadTestsVisibility);
  const controllerRef = useRef<AbortController | null>(null);
  const runVersionRef = useRef(0);
  const lastAutoRevisionRef = useRef(0);
  const layoutRef = useRef<HTMLDivElement>(null);
  const workspaceRef = useRef<HTMLElement>(null);
  const outputGridRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    let cancelled = false;
    void fetchCatalogue()
      .then(async (catalogue) => {
        if (cancelled) return;
        const loadedLevels = catalogue.levels;
        const loaded = loadProgress(catalogue);
        const receipts = Object.fromEntries(
          Object.entries(loaded.exercises)
            .filter(([, item]) => Boolean(item.receipt))
            .map(([key, item]) => [key, item.receipt!]),
        );
        try {
          const valid = new Set(await validateReceipts(receipts));
          for (const [key, item] of Object.entries(loaded.exercises)) {
            item.passed = Boolean(item.receipt) && valid.has(key);
          }
        } catch {
          for (const item of Object.values(loaded.exercises)) item.passed = false;
        }
        for (const item of Object.values(loaded.exercises)) {
          item.code = ensurePackageDeclaration(item.code);
        }
        setLevels(loadedLevels);
        setProgress(loaded);
      })
      .catch((error: unknown) => {
        if (!cancelled) {
          setLoadingError(
            error instanceof Error
              ? error.message
              : "Could not load the exercises.",
          );
        }
      });
    return () => {
      cancelled = true;
    };
  }, []);

  useEffect(() => {
    if (!progress) return;
    const handle = window.setTimeout(() => {
      saveProgress(progress);
      setSaved(true);
    }, 180);
    return () => window.clearTimeout(handle);
  }, [progress]);

  useEffect(() => {
    const syncFullscreen = () => {
      setFullscreen(document.fullscreenElement === workspaceRef.current);
    };
    document.addEventListener("fullscreenchange", syncFullscreen);
    return () =>
      document.removeEventListener("fullscreenchange", syncFullscreen);
  }, []);

  useEffect(() => {
    localStorage.setItem(PANE_SIZE_KEY, JSON.stringify(paneSizes));
  }, [paneSizes]);

  useEffect(() => {
    localStorage.setItem(TESTS_VISIBILITY_KEY, String(testsVisible));
  }, [testsVisible]);

  const level = useMemo(
    () =>
      levels.find((item) => item.key === progress?.currentExerciseKey) ?? levels[0],
    [levels, progress?.currentExerciseKey],
  );
  const levelProgress =
    progress && level ? progress.exercises[level.key] : undefined;
  const result = level ? runs[level.key] : undefined;

  const updateCode = useCallback(
    (code: string) => {
      if (!level) return;
      controllerRef.current?.abort();
      controllerRef.current = null;
      runVersionRef.current += 1;
      setRunning(false);
      setSaved(false);
      setEditRevision((value) => value + 1);
      setRuns((current) => {
        const next = { ...current };
        delete next[level.key];
        return next;
      });
      setProgress((current) => {
        if (!current) return current;
        const key = level.key;
        return {
          ...current,
          exercises: {
            ...current.exercises,
            [key]: { ...current.exercises[key], code },
          },
        };
      });
    },
    [level],
  );

  const runAllTests = useCallback(async () => {
    if (!level || !levelProgress) return;
    controllerRef.current?.abort();
    const controller = new AbortController();
    controllerRef.current = controller;
    const version = ++runVersionRef.current;
    lastAutoRevisionRef.current = editRevision;
    setRunning(true);

    try {
      const nextResult = await runTests(
        level.key,
        levelProgress.code,
        [],
        controller.signal,
      );
      if (version !== runVersionRef.current) return;
      setRuns((current) => ({
        ...current,
        [level.key]: nextResult,
      }));
      setProgress((current) => {
        if (!current) return current;
        const key = level.key;
        const item = current.exercises[key];
        return {
          ...current,
          exercises: {
            ...current.exercises,
            [key]: {
              ...item,
              code: item.code,
              attempts: item.attempts + 1,
              passed: item.passed || nextResult.passed,
              receipt: nextResult.receipt ?? item.receipt,
              bestPassed: Math.max(item.bestPassed, nextResult.passedCount),
              totalTests: level.tests.length,
              history: [
                {
                  id: crypto.randomUUID(),
                  at: Date.now(),
                  passed: nextResult.passedCount,
                  total: nextResult.totalCount,
                  durationMs: nextResult.durationMs,
                  outcome: outcomeFor(nextResult),
                },
                ...item.history,
              ].slice(0, 10),
            },
          },
        };
      });
    } catch (error) {
      if (error instanceof DOMException && error.name === "AbortError") return;
      const message =
        error instanceof Error ? error.message : "The test run failed.";
      setRuns((current) => ({
        ...current,
        [level.key]: failedRun(level.key, level.id, level.tests.length, message),
      }));
    } finally {
      if (version === runVersionRef.current) {
        controllerRef.current = null;
        setRunning(false);
      }
    }
  }, [editRevision, level, levelProgress]);

  useEffect(() => {
    if (
      !progress?.settings.autoTest ||
      !levelProgress ||
      running ||
      editRevision === 0 ||
      lastAutoRevisionRef.current === editRevision
    ) {
      return;
    }
    const revision = editRevision;
    const handle = window.setTimeout(() => {
      lastAutoRevisionRef.current = revision;
      void runAllTests();
    }, 900);
    return () => window.clearTimeout(handle);
  }, [
    editRevision,
    levelProgress?.code,
    progress?.settings.autoTest,
    runAllTests,
    running,
  ]);

  useEffect(() => {
    const runShortcut = (event: KeyboardEvent) => {
      if ((event.ctrlKey || event.metaKey) && event.key === "Enter") {
        event.preventDefault();
        void runAllTests();
      }
    };
    window.addEventListener("keydown", runShortcut);
    return () => window.removeEventListener("keydown", runShortcut);
  }, [runAllTests]);

  const stopRun = () => {
    controllerRef.current?.abort();
    controllerRef.current = null;
    runVersionRef.current += 1;
    setRunning(false);
  };

  const setAutoTest = (autoTest: boolean) => {
    setProgress((current) =>
      current
        ? {
            ...current,
            settings: { ...current.settings, autoTest },
          }
        : current,
    );
  };

  const selectLevel = (id: number) => {
    if (!progress || id < 1 || id > levels.length) return;
    const selected = levels.find((item) => item.id === id);
    if (!selected) return;
    setProgress({ ...progress, currentExerciseKey: selected.key });
    setEditRevision(0);
    lastAutoRevisionRef.current = 0;
  };

  const resetCode = () => {
    if (!level || !levelProgress) return;
    if (!window.confirm("Restore the starter code for this exercise?")) return;
    updateCode(level.starterCode);
  };

  const applyGofmt = async () => {
    if (!levelProgress) return;
    try {
      updateCode(await formatCode(levelProgress.code));
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "gofmt could not format code.";
      setRuns((current) => ({
        ...current,
        [level.key]: failedRun(level.key, level.id, level.tests.length, message),
      }));
    }
  };

  const copyCode = async () => {
    if (!levelProgress) return;
    await navigator.clipboard.writeText(levelProgress.code);
  };

  const downloadCode = () => {
    if (!levelProgress) return;
    const url = URL.createObjectURL(
      new Blob([levelProgress.code], {
        type: "text/x-go",
      }),
    );
    const anchor = document.createElement("a");
    anchor.href = url;
    anchor.download = `level-${level.id}-solution.go`;
    anchor.click();
    URL.revokeObjectURL(url);
  };

  const toggleFullscreen = async () => {
    try {
      if (document.fullscreenElement === workspaceRef.current) {
        await document.exitFullscreen();
      } else {
        await workspaceRef.current?.requestFullscreen();
      }
    } catch {
      setRuns((current) => ({
        ...current,
        [level.key]: failedRun(
          level.key,
          level.id,
          level.tests.length,
          "This browser did not allow fullscreen mode.",
        ),
      }));
    }
  };

  const startResize = (
    kind: ResizeKind,
    event: ReactPointerEvent<HTMLDivElement>,
  ) => {
    if (window.matchMedia("(max-width: 900px)").matches) return;
    event.preventDefault();
    const layoutRect = layoutRef.current?.getBoundingClientRect();
    const workspaceRect = workspaceRef.current?.getBoundingClientRect();
    const outputRect = outputGridRef.current?.getBoundingClientRect();
    if (!layoutRect || !workspaceRect) return;

    document.body.classList.add("resizing", `resizing-${kind}`);
    const move = (moveEvent: PointerEvent) => {
      setPaneSizes((current) => {
        if (kind === "brief") {
          return {
            ...current,
            brief: clamp(
              moveEvent.clientX - layoutRect.left,
              240,
              Math.max(280, Math.min(560, layoutRect.width - 520)),
            ),
          };
        }
        if (kind === "output") {
          return {
            ...current,
            output: clamp(
              workspaceRect.bottom - moveEvent.clientY,
              170,
              Math.max(220, workspaceRect.height - 260),
            ),
          };
        }
        if (!outputRect) return current;
        return {
          ...current,
          tests: clamp(
            outputRect.right - moveEvent.clientX,
            280,
            Math.max(320, outputRect.width - 320),
          ),
        };
      });
    };
    const finish = () => {
      document.body.classList.remove("resizing", `resizing-${kind}`);
      window.removeEventListener("pointermove", move);
      window.removeEventListener("pointerup", finish);
      window.removeEventListener("pointercancel", finish);
    };
    window.addEventListener("pointermove", move);
    window.addEventListener("pointerup", finish);
    window.addEventListener("pointercancel", finish);
  };

  const nudgePane = (kind: ResizeKind, delta: number) => {
    setPaneSizes((current) => ({
      ...current,
      [kind]: Math.max(160, current[kind] + delta),
    }));
  };

  const resetPane = (kind: ResizeKind) => {
    setPaneSizes((current) => ({
      ...current,
      [kind]: DEFAULT_PANE_SIZES[kind],
    }));
  };

  if (loadingError) {
    return (
      <main className="fatal-screen">
        <span>workspace unavailable</span>
        <h1>Could not load the Go exercises.</h1>
        <pre>{loadingError}</pre>
        <button onClick={() => window.location.reload()}>Retry</button>
      </main>
    );
  }

  if (!progress || !level || !levelProgress) {
    return (
      <main className="loading-screen">
        <span />
        <p>Loading workspace…</p>
      </main>
    );
  }

  const canMoveBack = level.id > 1;
  const canMoveForward =
    level.id < levels.length && Boolean(levelProgress.passed);

  return (
    <main className="assessment">
      <header className="topbar">
        <div className="wordmark">
          <span>imperative</span>
          <strong>/ go</strong>
        </div>
        <div className="level-stepper">
          <button
            onClick={() => selectLevel(level.id - 1)}
            disabled={!canMoveBack}
            aria-label="Previous exercise"
          >
            <ArrowLeft />
          </button>
          <span>
            Exercise {String(level.id).padStart(2, "0")}{" "}
            <em>/ {String(levels.length).padStart(2, "0")}</em>
          </span>
          <button
            onClick={() => selectLevel(level.id + 1)}
            disabled={!canMoveForward}
            aria-label="Next exercise"
            title={
              canMoveForward
                ? "Next exercise"
                : "Pass this exercise to continue"
            }
          >
            <ArrowRight />
          </button>
        </div>
      </header>

      <div
        className="practice-layout"
        ref={layoutRef}
        style={
          {
            "--brief-width": `${paneSizes.brief}px`,
          } as React.CSSProperties
        }
      >
        <ExerciseBrief level={level} passed={levelProgress.passed} />
        <ResizeHandle
          orientation="vertical"
          label="Resize instructions"
          onPointerDown={(event) => startResize("brief", event)}
          onNudge={(delta) => nudgePane("brief", delta)}
          onReset={() => resetPane("brief")}
        />

        <section
          className="workspace"
          ref={workspaceRef}
          style={
            {
              "--output-height": `${paneSizes.output}px`,
            } as React.CSSProperties
          }
        >
          <div className="editor-toolbar">
            <div className="file-state">
              <span className="live-dot" />
              <strong>solution.go</strong>
              <small>
                {saved ? <Check /> : null}
                {saved ? "saved locally" : "saving…"}
              </small>
            </div>
            <div className="editor-actions">
              <ToolbarButton label="Copy" onClick={() => void copyCode()}>
                <Copy />
              </ToolbarButton>
              <ToolbarButton label=".go" onClick={downloadCode}>
                <Download />
              </ToolbarButton>
              <ToolbarButton label="gofmt" onClick={() => void applyGofmt()}>
                <Sparkles />
              </ToolbarButton>
              <ToolbarButton
                label={fullscreen ? "Exit full screen" : "Full screen"}
                onClick={() => void toggleFullscreen()}
              >
                {fullscreen ? <Minimize2 /> : <Maximize2 />}
              </ToolbarButton>
            </div>
          </div>

          <div className="editor-surface">
            <CodeMirror
              value={levelProgress.code}
              onChange={updateCode}
              extensions={[go()]}
              theme="dark"
              height="100%"
              basicSetup={{
                lineNumbers: true,
                highlightActiveLineGutter: true,
                foldGutter: true,
                bracketMatching: true,
                closeBrackets: true,
                autocompletion: true,
                highlightActiveLine: true,
                highlightSelectionMatches: true,
                history: true,
                tabSize: 4,
              }}
              style={{ fontSize: progress.settings.fontSize }}
              aria-label={`Go editor for exercise ${level.id}`}
            />
          </div>

          <ResizeHandle
            orientation="horizontal"
            label="Resize editor and output"
            onPointerDown={(event) => startResize("output", event)}
            onNudge={(delta) => nudgePane("output", -delta)}
            onReset={() => resetPane("output")}
          />

          <div className="run-area">
            <div className="run-controls">
              <div>
                {running ? (
                  <button className="run-button stop" onClick={stopRun}>
                    <Square /> Stop
                  </button>
                ) : (
                  <button
                    className="run-button"
                    onClick={() => void runAllTests()}
                  >
                    <Play /> Test
                    <kbd>Ctrl ↵</kbd>
                  </button>
                )}
                <label className="auto-test">
                  <input
                    type="checkbox"
                    checked={progress.settings.autoTest}
                    onChange={(event) => setAutoTest(event.target.checked)}
                  />
                  <span />
                  Auto test
                </label>
              </div>
              <div className="run-status">
                <RunStatus result={result} running={running} />
                <button onClick={resetCode} title="Restore starter code">
                  <RotateCcw />
                </button>
              </div>
            </div>

            <div
              className={`output-grid${testsVisible ? "" : " tests-hidden"}`}
              ref={outputGridRef}
              style={
                {
                  "--tests-width": `${paneSizes.tests}px`,
                } as React.CSSProperties
              }
            >
              <Console
                result={result}
                tests={level.tests}
                running={running}
                testsVisible={testsVisible}
                onShowTests={() => setTestsVisible(true)}
              />
              {testsVisible ? (
                <>
                  <ResizeHandle
                    orientation="vertical"
                    label="Resize console and tests"
                    onPointerDown={(event) => startResize("tests", event)}
                    onNudge={(delta) => nudgePane("tests", -delta)}
                    onReset={() => resetPane("tests")}
                  />
                  <ExerciseTests
                    tests={level.tests}
                    onHide={() => setTestsVisible(false)}
                  />
                </>
              ) : null}
            </div>
          </div>
        </section>
      </div>
    </main>
  );
}

function ResizeHandle({
  orientation,
  label,
  onPointerDown,
  onNudge,
  onReset,
}: {
  orientation: "horizontal" | "vertical";
  label: string;
  onPointerDown: (event: ReactPointerEvent<HTMLDivElement>) => void;
  onNudge: (delta: number) => void;
  onReset: () => void;
}) {
  return (
    <div
      className={`resize-handle ${orientation}`}
      role="separator"
      aria-label={label}
      aria-orientation={orientation}
      tabIndex={0}
      title="Drag to resize. Double-click to reset."
      onPointerDown={onPointerDown}
      onDoubleClick={onReset}
      onKeyDown={(event) => {
        const backward =
          orientation === "vertical"
            ? event.key === "ArrowLeft"
            : event.key === "ArrowUp";
        const forward =
          orientation === "vertical"
            ? event.key === "ArrowRight"
            : event.key === "ArrowDown";
        if (!backward && !forward) return;
        event.preventDefault();
        onNudge(backward ? -16 : 16);
      }}
    >
      <span />
    </div>
  );
}

function ExerciseBrief({ level, passed }: { level: Level; passed: boolean }) {
  const instructions = level.instructions;
  return (
    <aside className="brief">
      <div className="brief-heading">
        <span>Exercise {String(level.id).padStart(2, "0")}</span>
        {passed ? <em>passed</em> : <em>{level.difficulty}</em>}
        <h1>{level.title}</h1>
        <code>{level.signature}</code>
      </div>

      <BriefSection title="Instructions">
        <p className="objective">{instructions.objective}</p>
        <p className="task-lead">Your solution must:</p>
        <ul className="instruction-list">
          <li>
            <strong>Accept:</strong> {instructions.input}
          </li>
          <li>
            <strong>Return:</strong> {instructions.output}
          </li>
          {instructions.constraints.map((rule) => (
            <li key={rule}>{rule}</li>
          ))}
        </ul>
        <p className="starter-note">{instructions.starterNote}</p>
      </BriefSection>

      <BriefSection title="Examples">
        <div className="examples">
          {instructions.examples.map((example) => (
            <div key={`${example.input}-${example.output}`}>
              <code>{example.input}</code>
              <span>→</span>
              <code>{example.output}</code>
            </div>
          ))}
        </div>
      </BriefSection>

      <BriefSection title="Allowed">
        <p className="allowed-label">Built-ins</p>
        <code className="token-list">
          {instructions.allowedBuiltins.join(" · ")}
        </code>
        <p className="allowed-label">Packages</p>
        <code className="token-list">
          {instructions.allowedPackages.join(" · ")}
        </code>
        <p className="restriction">
          Everything not listed above is blocked for this exercise.
        </p>
        <p className="allowed-label">Avoid</p>
        <ul className="pitfall-list">
          {instructions.commonPitfalls.map((pitfall) => (
            <li key={pitfall}>{pitfall}</li>
          ))}
        </ul>
      </BriefSection>
    </aside>
  );
}

function BriefSection({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <section className="brief-section">
      <h2>{title}</h2>
      <div>{children}</div>
    </section>
  );
}

function ToolbarButton({
  label,
  onClick,
  children,
}: {
  label: string;
  onClick: () => void;
  children: React.ReactNode;
}) {
  return (
    <button onClick={onClick}>
      {children}
      <span>{label}</span>
    </button>
  );
}

function RunStatus({
  result,
  running,
}: {
  result?: RunResult;
  running: boolean;
}) {
  if (running) return <span className="status-running">running tests…</span>;
  if (!result) return <span>ready</span>;
  if (result.passed) {
    return (
      <span className="status-pass">
        <Check /> all tests passed
      </span>
    );
  }
  if (result.compileError) return <span className="status-fail">compile error</span>;
  if (result.runtimeError) return <span className="status-fail">runtime error</span>;
  return (
    <span className="status-fail">
      {result.passedCount}/{result.totalCount} tests passed
    </span>
  );
}

function Console({
  result,
  tests,
  running,
  testsVisible,
  onShowTests,
}: {
  result?: RunResult;
  tests: Level["tests"];
  running: boolean;
  testsVisible: boolean;
  onShowTests: () => void;
}) {
  const streams = buildConsoleStreams(result, tests);

  return (
    <section className="console">
      <h2>
        <span>
          <Terminal /> Console
        </span>
        {!testsVisible ? (
          <button onClick={onShowTests} title="Show exercise tests">
            <Eye /> Show tests
          </button>
        ) : null}
      </h2>
      <div className="console-content" aria-live="polite">
        {running ? <p className="console-muted">running assessment tests…</p> : null}
        {!running && streams.length === 0 ? (
          <p className="console-muted">
            Program output and errors appear here.
          </p>
        ) : null}
        {streams.map((stream) => (
          <div
            className={`console-stream${stream.error ? " error" : ""}`}
            key={stream.label}
          >
            <span>{stream.label}</span>
            <pre>{stream.value}</pre>
          </div>
        ))}
      </div>
    </section>
  );
}

function ExerciseTests({
  tests,
  onHide,
}: {
  tests: Level["tests"];
  onHide: () => void;
}) {
  return (
    <section className="exercise-tests">
      <h2>
        <span>
          <ListChecks /> Tests ({tests.length})
        </span>
        <button onClick={onHide} title="Hide exercise tests">
          <EyeOff /> Hide
        </button>
      </h2>
      <div className="exercise-test-list">
        {tests.map((test, index) => (
          <article className="exercise-test" key={test.id}>
            <div>
              <span>#{index + 1}</span>
              <strong>{test.name}</strong>
            </div>
            <p>{test.purpose}</p>
            <dl>
              <div>
                <dt>Input</dt>
                <dd>
                  <code>{test.input}</code>
                </dd>
              </div>
              <div>
                <dt>Need</dt>
                <dd>
                  <code>{test.expected}</code>
                </dd>
              </div>
            </dl>
          </article>
        ))}
      </div>
    </section>
  );
}

function outcomeFor(result: RunResult): string {
  if (result.passed) return "All tests passed";
  if (result.compileError) return "Compile error";
  if (result.runtimeError) return "Runtime error";
  if (result.timedOut) return "Timeout";
  return `${result.passedCount}/${result.totalCount} tests passed`;
}

function failedRun(
  exerciseKey: string,
  levelId: number,
  totalCount: number,
  message: string,
): RunResult {
  return {
    exerciseKey,
    levelId,
    passed: false,
    passedCount: 0,
    totalCount,
    runtimeError: message,
    timedOut: false,
    stopped: false,
    stdout: "",
    stderr: "",
    durationMs: 0,
    results: [],
    sourceHash: "",
  };
}

function loadPaneSizes(): PaneSizes {
  try {
    const saved = JSON.parse(localStorage.getItem(PANE_SIZE_KEY) ?? "");
    return {
      brief:
        typeof saved.brief === "number"
          ? saved.brief
          : DEFAULT_PANE_SIZES.brief,
      output:
        typeof saved.output === "number"
          ? saved.output
          : DEFAULT_PANE_SIZES.output,
      tests:
        typeof saved.tests === "number"
          ? saved.tests
          : DEFAULT_PANE_SIZES.tests,
    };
  } catch {
    return DEFAULT_PANE_SIZES;
  }
}

function loadTestsVisibility(): boolean {
  return localStorage.getItem(TESTS_VISIBILITY_KEY) !== "false";
}

function clamp(value: number, minimum: number, maximum: number): number {
  return Math.min(maximum, Math.max(minimum, Math.round(value)));
}

function ensurePackageDeclaration(source: string): string {
  return /^\s*package\s+\w+/m.test(source)
    ? source
    : `package main\n\n${source}`;
}

export default App;
