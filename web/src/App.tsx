import {
  type ChangeEvent,
  type ReactNode,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import CodeMirror from "@uiw/react-codemirror";
import { go } from "@codemirror/lang-go";
import {
  AlertTriangle,
  ArrowLeft,
  ArrowRight,
  Check,
  CheckCircle2,
  ChevronDown,
  Circle,
  Clock3,
  Copy,
  Download,
  FileJson,
  FlaskConical,
  History,
  Lightbulb,
  ListChecks,
  LockKeyhole,
  Maximize2,
  Menu,
  Minimize2,
  Minus,
  Moon,
  Pause,
  Play,
  Plus,
  RotateCcw,
  Save,
  ShieldAlert,
  ShieldCheck,
  Sparkles,
  Square,
  Sun,
  Trophy,
  Terminal,
  Upload,
  X,
} from "lucide-react";
import {
  fetchLevels,
  fetchRunnerConfig,
  formatCode,
  runTests,
  validateReceipts,
} from "./api";
import {
  ASSESSMENT_SECONDS,
  STORAGE_KEY,
  createProgress,
  isLevelUnlocked,
  loadProgress,
  saveProgress,
  sequentialCompleted,
  settleTimer,
  validateImport,
} from "./storage";
import type {
  Level,
  RunnerConfig,
  RunResult,
  SavedProgress,
  TestResult,
  VisibleTest,
} from "./types";

type Toast = { id: number; message: string; tone: "info" | "good" | "bad" };

function formatDuration(seconds: number): string {
  const safe = Math.max(0, Math.floor(seconds));
  const hours = Math.floor(safe / 3600);
  const minutes = Math.floor((safe % 3600) / 60);
  const remainder = safe % 60;
  return [hours, minutes, remainder]
    .map((part) => String(part).padStart(2, "0"))
    .join(":");
}

function formatSpent(seconds: number): string {
  if (seconds < 60) return `${seconds}s`;
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  return hours ? `${hours}h ${minutes}m` : `${minutes}m`;
}

function outcomeFor(result: RunResult): string {
  if (result.failureKind === "output") return "Output limit";
  if (result.failureKind === "capacity") return "Runner at capacity";
  if (result.failureKind === "cleanup") return "Cleanup error";
  if (result.timedOut) return "Timeout";
  if (result.stopped) return "Stopped";
  if (result.compileError) return "Compile error";
  if (result.runtimeError) return "Runtime error";
  if (result.passed) return "Level passed";
  return `${result.passedCount}/${result.totalCount} tests passed`;
}

function download(name: string, contents: string, type: string): void {
  const url = URL.createObjectURL(new Blob([contents], { type }));
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = name;
  anchor.click();
  URL.revokeObjectURL(url);
}

function App() {
  const [levels, setLevels] = useState<Level[]>([]);
  const [runnerConfig, setRunnerConfig] = useState<RunnerConfig | null>(null);
  const [progress, setProgress] = useState<SavedProgress | null>(null);
  const [runs, setRuns] = useState<Record<string, RunResult>>({});
  const [loadingError, setLoadingError] = useState("");
  const [running, setRunning] = useState(false);
  const [savedState, setSavedState] = useState<"saved" | "modified">("saved");
  const [editRevision, setEditRevision] = useState(0);
  const [expandedTests, setExpandedTests] = useState<Record<string, boolean>>({});
  const [showAllDetails, setShowAllDetails] = useState(false);
  const [showHistory, setShowHistory] = useState(false);
  const [showSummary, setShowSummary] = useState(false);
  const [mobileRail, setMobileRail] = useState(false);
  const [editorFullscreen, setEditorFullscreen] = useState(false);
  const [toasts, setToasts] = useState<Toast[]>([]);
  const controllerRef = useRef<AbortController | null>(null);
  const runVersionRef = useRef(0);
  const lastAutoRevisionRef = useRef(0);
  const importRef = useRef<HTMLInputElement>(null);
  const editorPanelRef = useRef<HTMLDivElement>(null);
  const toastIDRef = useRef(0);

  const toast = useCallback(
    (message: string, tone: Toast["tone"] = "info") => {
      const id = ++toastIDRef.current;
      setToasts((current) => [...current.slice(-2), { id, message, tone }]);
      window.setTimeout(
        () => setToasts((current) => current.filter((item) => item.id !== id)),
        3600,
      );
    },
    [],
  );

  useEffect(() => {
    let cancelled = false;
    void Promise.all([fetchLevels(), fetchRunnerConfig()])
      .then(async ([loadedLevels, loadedRunnerConfig]) => {
        if (cancelled) return;
        const loaded = loadProgress(loadedLevels);
        const receiptMap = Object.fromEntries(
          Object.entries(loaded.levels)
            .filter(([, item]) => Boolean(item.receipt))
            .map(([id, item]) => [id, item.receipt!]),
        );
        try {
          const valid = new Set(await validateReceipts(receiptMap));
          for (const [id, item] of Object.entries(loaded.levels)) {
            item.passed = Boolean(item.receipt) && valid.has(Number(id));
          }
        } catch {
          for (const item of Object.values(loaded.levels)) item.passed = false;
        }
        setLevels(loadedLevels);
        setRunnerConfig(loadedRunnerConfig);
        setProgress({ ...loaded, timer: settleTimer(loaded.timer) });
      })
      .catch((error: unknown) => {
        if (!cancelled) {
          setLoadingError(
            error instanceof Error
              ? error.message
              : "Could not load the assessment.",
          );
        }
      });
    return () => {
      cancelled = true;
    };
  }, []);

  useEffect(() => {
    if (!progress) return;
    document.documentElement.dataset.theme = progress.settings.theme;
    document.documentElement.dataset.motion = progress.settings.reducedMotion
      ? "reduced"
      : "full";
    const handle = window.setTimeout(() => {
      saveProgress(progress);
      setSavedState("saved");
    }, 180);
    return () => window.clearTimeout(handle);
  }, [progress]);

  useEffect(() => {
    const syncFullscreenState = () => {
      setEditorFullscreen(document.fullscreenElement === editorPanelRef.current);
    };
    document.addEventListener("fullscreenchange", syncFullscreenState);
    return () =>
      document.removeEventListener("fullscreenchange", syncFullscreenState);
  }, []);

  useEffect(() => {
    if (!progress?.timer.running) return;
    const handle = window.setInterval(() => {
      setProgress((current) => {
        if (!current) return current;
        const timer = settleTimer(current.timer);
        const levelKey = String(current.currentLevelId);
        const levelProgress = current.levels[levelKey];
        return {
          ...current,
          timer,
          levels: {
            ...current.levels,
            [levelKey]: {
              ...levelProgress,
              timeSpentSeconds: levelProgress.timeSpentSeconds + 1,
            },
          },
        };
      });
    }, 1000);
    return () => window.clearInterval(handle);
  }, [progress?.timer.running]);

  const currentLevel = useMemo(
    () =>
      levels.find((level) => level.id === progress?.currentLevelId) ?? levels[0],
    [levels, progress?.currentLevelId],
  );
  const currentLevelProgress =
    progress && currentLevel
      ? progress.levels[String(currentLevel.id)]
      : undefined;
  const currentRun =
    currentLevel && runs[String(currentLevel.id)]
      ? runs[String(currentLevel.id)]
      : undefined;
  const passedCount = progress
    ? Object.values(progress.levels).filter((item) => item.passed).length
    : 0;

  const updateCode = useCallback(
    (code: string) => {
      if (!currentLevel) return;
      if (controllerRef.current) {
        controllerRef.current.abort();
        controllerRef.current = null;
        runVersionRef.current += 1;
        setRunning(false);
      }
      setSavedState("modified");
      setEditRevision((revision) => revision + 1);
      setRuns((current) => {
        const key = String(currentLevel.id);
        if (!current[key]) return current;
        const next = { ...current };
        delete next[key];
        return next;
      });
      setProgress((current) => {
        if (!current) return current;
        const key = String(currentLevel.id);
        return {
          ...current,
          levels: {
            ...current.levels,
            [key]: { ...current.levels[key], code },
          },
        };
      });
    },
    [currentLevel],
  );

  const runSelection = useCallback(
    async (testIds: string[] = []) => {
      if (!currentLevel || !currentLevelProgress) return;
      controllerRef.current?.abort();
      const controller = new AbortController();
      controllerRef.current = controller;
      const version = ++runVersionRef.current;
      lastAutoRevisionRef.current = editRevision;
      setRunning(true);
      setProgress((current) => {
        if (!current) return current;
        const key = String(currentLevel.id);
        return {
          ...current,
          levels: {
            ...current.levels,
            [key]: {
              ...current.levels[key],
              attempts: current.levels[key].attempts + 1,
            },
          },
        };
      });
      try {
        const result = await runTests(
          currentLevel.id,
          currentLevelProgress.code,
          testIds,
          controller.signal,
        );
        if (version !== runVersionRef.current) return;
        setRuns((current) => {
          if (testIds.length === 0 || !current[String(currentLevel.id)]) {
            return { ...current, [String(currentLevel.id)]: result };
          }
          const prior = current[String(currentLevel.id)];
          const replacements = new Map(
            result.results.map((item) => [item.id, item]),
          );
          const merged = prior.results.map(
            (item) => replacements.get(item.id) ?? item,
          );
          return {
            ...current,
            [String(currentLevel.id)]: {
              ...result,
              results: merged,
              passedCount: merged.filter((item) => item.passed).length,
              totalCount: merged.length,
            },
          };
        });
        setProgress((current) => {
          if (!current) return current;
          const key = String(currentLevel.id);
          const item = current.levels[key];
          const fullRun = testIds.length === 0;
          const passed = item.passed || (fullRun && result.passed);
          const history = [
            {
              id: crypto.randomUUID(),
              at: Date.now(),
              passed: result.passedCount,
              total: result.totalCount,
              durationMs: result.durationMs,
              outcome: outcomeFor(result),
            },
            ...item.history,
          ].slice(0, 10);
          return {
            ...current,
            levels: {
              ...current.levels,
              [key]: {
                ...item,
                code: result.formattedCode ?? item.code,
                passed,
                receipt: result.receipt ?? item.receipt,
                bestPassed: Math.max(item.bestPassed, result.passedCount),
                totalTests: currentLevel.tests.length,
                history,
              },
            },
          };
        });
        if (result.passed) {
          toast(
            currentLevel.id === 7
              ? "Assessment target achieved: 7/9."
              : `Level ${currentLevel.id} passed. The next level is unlocked.`,
            "good",
          );
        } else if (result.compileError) {
          toast("Compilation needs attention.", "bad");
        } else if (result.timedOut) {
          toast("The run timed out and was terminated.", "bad");
        } else if (result.failureKind === "capacity") {
          toast("The runner is busy. Wait for the current run, then retry.", "bad");
        } else if (result.failureKind === "cleanup") {
          toast("The sandbox container could not be cleaned up safely.", "bad");
        } else if (result.failureKind === "output") {
          toast("Program output exceeded the 256 KiB limit.", "bad");
        } else if (result.failureKind === "startup") {
          toast("The sandbox container could not start. Check Docker Desktop.", "bad");
        } else if (result.runtimeError) {
          toast(result.runtimeError, "bad");
        } else {
          toast(`${result.passedCount}/${result.totalCount} tests passed.`);
        }
      } catch (error) {
        if (error instanceof DOMException && error.name === "AbortError") {
          if (version === runVersionRef.current) toast("Run stopped.");
          return;
        }
        toast(
          error instanceof Error ? error.message : "The test run failed.",
          "bad",
        );
      } finally {
        if (version === runVersionRef.current) {
          setRunning(false);
          controllerRef.current = null;
        }
      }
    },
    [currentLevel, currentLevelProgress?.code, editRevision, toast],
  );

  useEffect(() => {
    if (
      !progress?.settings.autoTest ||
      !currentLevelProgress ||
      running ||
      editRevision === 0 ||
      lastAutoRevisionRef.current === editRevision
    )
      return;
    const scheduledRevision = editRevision;
    const handle = window.setTimeout(() => {
      lastAutoRevisionRef.current = scheduledRevision;
      void runSelection();
    }, 1100);
    return () => window.clearTimeout(handle);
  }, [
    currentLevelProgress?.code,
    currentLevel?.id,
    editRevision,
    progress?.settings.autoTest,
    runSelection,
    running,
  ]);

  useEffect(() => {
    const keydown = (event: KeyboardEvent) => {
      if ((event.ctrlKey || event.metaKey) && event.key === "Enter") {
        event.preventDefault();
        void runSelection();
      }
    };
    window.addEventListener("keydown", keydown);
    return () => window.removeEventListener("keydown", keydown);
  }, [runSelection]);

  const patchSettings = useCallback(
    (patch: Partial<SavedProgress["settings"]>) => {
      setProgress((current) =>
        current
          ? { ...current, settings: { ...current.settings, ...patch } }
          : current,
      );
    },
    [],
  );

  const selectLevel = (level: Level) => {
    if (!progress || !isLevelUnlocked(level.id, progress)) {
      toast("Pass the previous level or enable practice mode to open this one.");
      return;
    }
    setProgress({ ...progress, currentLevelId: level.id });
    setEditRevision(0);
    lastAutoRevisionRef.current = 0;
    setMobileRail(false);
    setShowSummary(false);
  };

  const resetLevel = () => {
    if (!currentLevel || !progress) return;
    if (
      !window.confirm(
        `Reset Level ${currentLevel.id} to its starter code? Attempts and time are preserved.`,
      )
    )
      return;
    const key = String(currentLevel.id);
    setProgress({
      ...progress,
      levels: {
        ...progress.levels,
        [key]: {
          ...progress.levels[key],
          code: currentLevel.starterCode,
          passed: false,
          receipt: undefined,
          bestPassed: 0,
        },
      },
    });
    setRuns((current) => {
      const next = { ...current };
      delete next[key];
      return next;
    });
    toast("Starter code restored.");
  };

  const resetAssessment = () => {
    if (!progress) return;
    const answer = window.prompt(
      "This clears all code, receipts, history, and timer data. Type RESET to continue.",
    );
    if (answer !== "RESET") {
      if (answer !== null) toast("Reset cancelled: confirmation did not match.");
      return;
    }
    localStorage.removeItem(STORAGE_KEY);
    setProgress(createProgress(levels));
    setRuns({});
    setShowSummary(false);
    toast("Assessment reset.");
  };

  const exportProgress = () => {
    if (!progress) return;
    download(
      `imperative-go-progress-${new Date().toISOString().slice(0, 10)}.json`,
      JSON.stringify({ ...progress, updatedAt: Date.now() }, null, 2),
      "application/json",
    );
    toast("Progress backup downloaded.", "good");
  };

  const importProgress = async (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    event.target.value = "";
    if (!file) return;
    try {
      const imported = validateImport(JSON.parse(await file.text()), levels);
      const receiptMap = Object.fromEntries(
        Object.entries(imported.levels)
          .filter(([, item]) => Boolean(item.receipt))
          .map(([id, item]) => [id, item.receipt!]),
      );
      const valid = new Set(await validateReceipts(receiptMap));
      for (const [id, item] of Object.entries(imported.levels)) {
        item.passed = Boolean(item.receipt) && valid.has(Number(id));
      }
      setProgress(imported);
      setRuns({});
      toast("Progress imported and pass receipts revalidated.", "good");
    } catch (error) {
      toast(
        error instanceof Error ? error.message : "Could not import this backup.",
        "bad",
      );
    }
  };

  const useHint = (index: number) => {
    if (!currentLevel || !progress) return;
    const key = String(currentLevel.id);
    if (progress.levels[key].hintsUsed.includes(index)) return;
    setProgress({
      ...progress,
      levels: {
        ...progress.levels,
        [key]: {
          ...progress.levels[key],
          hintsUsed: [...progress.levels[key].hintsUsed, index].sort(),
        },
      },
    });
  };

  const gofmt = async () => {
    if (!currentLevelProgress) return;
    try {
      updateCode(await formatCode(currentLevelProgress.code));
      toast("Code formatted with gofmt.", "good");
    } catch (error) {
      toast(
        error instanceof Error ? error.message : "gofmt could not format code.",
        "bad",
      );
    }
  };

  const toggleEditorFullscreen = async () => {
    try {
      if (document.fullscreenElement === editorPanelRef.current) {
        await document.exitFullscreen();
      } else {
        await editorPanelRef.current?.requestFullscreen();
      }
    } catch {
      toast("This browser did not allow fullscreen mode.", "bad");
    }
  };

  if (loadingError) {
    return (
      <main className="fatal-state">
        <ShieldAlert size={36} />
        <h1>The local assessment server is not ready</h1>
        <p>{loadingError}</p>
        <button onClick={() => window.location.reload()}>Try again</button>
      </main>
    );
  }

  if (!progress || !currentLevel || !currentLevelProgress || !runnerConfig) {
    return (
      <main className="loading-state" aria-live="polite">
        <div className="loading-mark">
          <span />
          <span />
          <span />
        </div>
        <p>Preparing the assessment workspace…</p>
      </main>
    );
  }

  const remaining = progress.timer.durationSeconds - progress.timer.elapsedSeconds;
  const failedTestIDs =
    currentRun?.results
      .filter((result) => !result.passed)
      .map((result) => result.id) ?? [];

  return (
    <div className="app-shell">
      <header className="assessment-header">
        <div className="brand-block">
          <button
            className="icon-button mobile-menu"
            onClick={() => setMobileRail(true)}
            aria-label="Open level navigation"
          >
            <Menu />
          </button>
          <div className="brand-glyph" aria-hidden="true">
            <span>go</span>
          </div>
          <div>
            <p className="eyebrow">Local practice environment</p>
            <h1>Imperative Go Practice Assessment</h1>
            <div
              className={`mode-badge ${runnerConfig.runnerMode}`}
              title={`${runnerConfig.message} Execution toolchain: ${runnerConfig.goVersion}`}
            >
              {runnerConfig.runnerMode === "docker" ? (
                <ShieldCheck />
              ) : (
                <ShieldAlert />
              )}
              {runnerConfig.runnerMode === "docker"
                ? "Docker sandbox"
                : "Local runner — trusted code only"}
            </div>
          </div>
        </div>

        <div className="header-progress" aria-label="Assessment progress">
          <div className="progress-copy">
            <span>Level {currentLevel.id} of 9</span>
            <strong>{passedCount}/9 passed</strong>
          </div>
          <div
            className="progress-track"
            role="progressbar"
            aria-valuemin={0}
            aria-valuemax={9}
            aria-valuenow={passedCount}
          >
            <span style={{ width: `${(passedCount / 9) * 100}%` }} />
          </div>
          {passedCount >= 7 && (
            <button
              className="target-chip"
              onClick={() => setShowSummary(true)}
            >
              <Trophy size={15} /> Target achieved
            </button>
          )}
        </div>

        <div className="timer-block">
          <div className={`timer ${remaining < 1800 ? "timer-warning" : ""}`}>
            <Clock3 size={17} />
            <span>{formatDuration(remaining)}</span>
          </div>
          <button
            className="icon-button"
            onClick={() =>
              setProgress({
                ...progress,
                timer: {
                  ...settleTimer(progress.timer),
                  running: !progress.timer.running,
                  lastTickAt: Date.now(),
                },
              })
            }
            aria-label={progress.timer.running ? "Pause timer" : "Resume timer"}
            title={progress.timer.running ? "Pause timer" : "Resume timer"}
          >
            {progress.timer.running ? <Pause /> : <Play />}
          </button>
          <button
            className="icon-button"
            onClick={() => {
              if (window.confirm("Reset the six-hour timer?")) {
                setProgress({
                  ...progress,
                  timer: {
                    durationSeconds: ASSESSMENT_SECONDS,
                    elapsedSeconds: 0,
                    running: true,
                    lastTickAt: Date.now(),
                  },
                });
              }
            }}
            aria-label="Reset timer"
            title="Reset timer"
          >
            <RotateCcw />
          </button>
          <button
            className="icon-button"
            onClick={() =>
              patchSettings({
                theme: progress.settings.theme === "light" ? "dark" : "light",
              })
            }
            aria-label={`Use ${progress.settings.theme === "light" ? "dark" : "light"} theme`}
          >
            {progress.settings.theme === "light" ? <Moon /> : <Sun />}
          </button>
        </div>
      </header>

      <div className="assessment-layout">
        <LevelRail
          levels={levels}
          progress={progress}
          selected={currentLevel.id}
          open={mobileRail}
          onClose={() => setMobileRail(false)}
          onSelect={selectLevel}
          onPractice={(practiceMode) => patchSettings({ practiceMode })}
          onSummary={() => {
            setShowSummary(true);
            setMobileRail(false);
          }}
        />

        {showSummary ? (
          <Summary
            levels={levels}
            progress={progress}
            onClose={() => setShowSummary(false)}
            onSelect={selectLevel}
            onExport={exportProgress}
          />
        ) : (
          <main className="workspace">
            <section className="level-heading">
              <div>
                <div className="level-kicker">
                  <span>LEVEL {String(currentLevel.id).padStart(2, "0")}</span>
                  <span>{currentLevel.difficulty}</span>
                  {currentLevel.stretch && (
                    <span className="stretch-badge">
                      <Sparkles size={13} /> Stretch level
                    </span>
                  )}
                </div>
                <h2>{currentLevel.title}</h2>
                <p>{currentLevel.topic}</p>
              </div>
              <div className="level-metrics">
                <Metric
                  label="Attempts"
                  value={String(currentLevelProgress.attempts)}
                />
                <Metric
                  label="Time here"
                  value={formatSpent(currentLevelProgress.timeSpentSeconds)}
                />
                <Metric
                  label="Best"
                  value={`${currentLevelProgress.bestPassed}/${currentLevel.tests.length}`}
                />
              </div>
            </section>

            <InstructionsPanel
              level={currentLevel}
              hintsUsed={currentLevelProgress.hintsUsed}
              onUseHint={useHint}
            />

            <section className="workbench">
              <div className="editor-panel panel" ref={editorPanelRef}>
                <div className="panel-toolbar">
                  <div className="panel-title">
                    <span className="status-dot" />
                    <strong>solution.go</strong>
                    <span className={`save-state ${savedState}`}>
                      {savedState === "saved" ? <Save /> : <Circle />}
                      {savedState === "saved" ? "Saved locally" : "Modified"}
                    </span>
                  </div>
                  <div className="toolbar-actions">
                    <button
                      className="text-button"
                      onClick={() => {
                        void navigator.clipboard.writeText(
                          `package main\n\n${currentLevelProgress.code}`,
                        );
                        toast("Code copied.", "good");
                      }}
                    >
                      <Copy /> Copy
                    </button>
                    <button
                      className="text-button"
                      onClick={() =>
                        download(
                          `level-${currentLevel.id}-solution.go`,
                          `package main\n\n${currentLevelProgress.code}`,
                          "text/x-go",
                        )
                      }
                    >
                      <Download /> .go
                    </button>
                    <button className="text-button" onClick={() => void gofmt()}>
                      <Sparkles /> gofmt
                    </button>
                    <button
                      className="text-button"
                      onClick={() => void toggleEditorFullscreen()}
                      aria-label={
                        editorFullscreen
                          ? "Exit editor fullscreen"
                          : "Open editor fullscreen"
                      }
                    >
                      {editorFullscreen ? <Minimize2 /> : <Maximize2 />}
                      {editorFullscreen ? "Exit full screen" : "Full screen"}
                    </button>
                  </div>
                </div>
                <div className="package-line" aria-label="Protected package line">
                  <span>1</span>
                  <code>package main</code>
                  <em>protected</em>
                </div>
                <div className="editor-canvas">
                  <CodeMirror
                    value={currentLevelProgress.code}
                    onChange={updateCode}
                    extensions={[go()]}
                    theme={progress.settings.theme}
                    basicSetup={{
                      lineNumbers: true,
                      highlightActiveLineGutter: true,
                      foldGutter: true,
                      allowMultipleSelections: true,
                      indentOnInput: true,
                      bracketMatching: true,
                      closeBrackets: true,
                      autocompletion: true,
                      rectangularSelection: true,
                      crosshairCursor: false,
                      highlightActiveLine: true,
                      highlightSelectionMatches: true,
                      closeBracketsKeymap: true,
                      defaultKeymap: true,
                      searchKeymap: true,
                      historyKeymap: true,
                      foldKeymap: true,
                      completionKeymap: true,
                      lintKeymap: true,
                    }}
                    height="100%"
                    style={{ fontSize: progress.settings.fontSize }}
                    aria-label={`Go editor for level ${currentLevel.id}`}
                  />
                </div>
                <div className="editor-footer">
                  <span>package main is injected and repaired before every run.</span>
                  <div className="font-controls" aria-label="Editor font size">
                    <button
                      onClick={() =>
                        patchSettings({
                          fontSize: Math.max(
                            12,
                            progress.settings.fontSize - 1,
                          ),
                        })
                      }
                      aria-label="Decrease editor font size"
                    >
                      <Minus />
                    </button>
                    <span>{progress.settings.fontSize}px</span>
                    <button
                      onClick={() =>
                        patchSettings({
                          fontSize: Math.min(
                            22,
                            progress.settings.fontSize + 1,
                          ),
                        })
                      }
                      aria-label="Increase editor font size"
                    >
                      <Plus />
                    </button>
                  </div>
                </div>
              </div>

              <ResultsPanel
                level={currentLevel}
                result={currentRun}
                running={running}
                runnerMode={runnerConfig.runnerMode}
                expanded={expandedTests}
                showAll={showAllDetails}
                onToggle={(id) =>
                  setExpandedTests((current) => ({
                    ...current,
                    [id]: !current[id],
                  }))
                }
                onShowAll={() => setShowAllDetails((value) => !value)}
                onRunOne={(id) => void runSelection([id])}
              />
            </section>

            <section className="test-dock">
              <div className="dock-primary">
                {running ? (
                  <button
                    className="test-button stop-button"
                    onClick={() => controllerRef.current?.abort()}
                  >
                    <Square /> Stop run
                  </button>
                ) : (
                  <button
                    className="test-button"
                    onClick={() => void runSelection()}
                  >
                    <FlaskConical /> Test
                    <kbd>Ctrl ↵</kbd>
                  </button>
                )}
                <button
                  className="secondary-button"
                  disabled={!failedTestIDs.length || running}
                  onClick={() => void runSelection(failedTestIDs)}
                >
                  <RotateCcw /> Run failed tests
                </button>
                <button
                  className="toggle-label"
                  role="switch"
                  aria-checked={progress.settings.autoTest}
                  onClick={() =>
                    patchSettings({ autoTest: !progress.settings.autoTest })
                  }
                >
                  <span
                    className={`toggle ${progress.settings.autoTest ? "checked" : ""}`}
                  />
                  Auto-test after pause
                </button>
              </div>
              <div className="dock-secondary">
                <button
                  className="text-button"
                  onClick={() => setShowHistory((value) => !value)}
                >
                  <History /> History
                </button>
                <button className="text-button danger-text" onClick={resetLevel}>
                  <RotateCcw /> Reset level
                </button>
              </div>
            </section>

            {showHistory && (
              <HistoryPanel
                entries={currentLevelProgress.history}
                onClose={() => setShowHistory(false)}
              />
            )}

            <nav className="level-pagination" aria-label="Level pagination">
              <button
                disabled={currentLevel.id === 1}
                onClick={() => selectLevel(levels[currentLevel.id - 2])}
              >
                <ArrowLeft /> Previous level
              </button>
              <span>
                {currentLevelProgress.passed ? (
                  <>
                    <CheckCircle2 /> Level passed
                  </>
                ) : (
                  <>Pass every visible test to complete this level</>
                )}
              </span>
              <button
                disabled={
                  currentLevel.id === 9 ||
                  !isLevelUnlocked(currentLevel.id + 1, progress)
                }
                onClick={() => selectLevel(levels[currentLevel.id])}
              >
                Next level <ArrowRight />
              </button>
            </nav>
          </main>
        )}
      </div>

      {runnerConfig.runnerMode === "local" && (
        <aside className="local-warning">
          <ShieldAlert />
          Local runner — trusted code only. Submitted Go code runs with your
          current user permissions; time and output limits are not a security
          sandbox.
        </aside>
      )}

      <div className="utility-corner">
        <button onClick={exportProgress}>
          <FileJson /> Export
        </button>
        <button onClick={() => importRef.current?.click()}>
          <Upload /> Import
        </button>
        <button className="danger-text" onClick={resetAssessment}>
          Reset all
        </button>
        <label>
          <input
            type="checkbox"
            checked={progress.settings.reducedMotion}
            onChange={(event) =>
              patchSettings({ reducedMotion: event.target.checked })
            }
          />
          Reduced motion
        </label>
      </div>
      <input
        ref={importRef}
        className="sr-only"
        type="file"
        accept="application/json,.json"
        onChange={(event) => void importProgress(event)}
      />

      <div className="toast-stack" aria-live="polite">
        {toasts.map((item) => (
          <div key={item.id} className={`toast ${item.tone}`}>
            {item.tone === "good" ? (
              <CheckCircle2 />
            ) : item.tone === "bad" ? (
              <AlertTriangle />
            ) : (
              <Circle />
            )}
            {item.message}
          </div>
        ))}
      </div>
    </div>
  );
}

function LevelRail({
  levels,
  progress,
  selected,
  open,
  onClose,
  onSelect,
  onPractice,
  onSummary,
}: {
  levels: Level[];
  progress: SavedProgress;
  selected: number;
  open: boolean;
  onClose: () => void;
  onSelect: (level: Level) => void;
  onPractice: (enabled: boolean) => void;
  onSummary: () => void;
}) {
  return (
    <>
      {open && <button className="rail-scrim" onClick={onClose} aria-label="Close navigation" />}
      <aside className={`level-rail ${open ? "open" : ""}`}>
        <div className="rail-heading">
          <span>Assessment path</span>
          <button
            className="icon-button rail-close"
            onClick={onClose}
            aria-label="Close level navigation"
          >
            <X />
          </button>
        </div>
        <div className="rail-rule">
          <span />
          <em>7 required · 2 stretch</em>
        </div>
        <nav>
          {levels.map((level) => {
            const item = progress.levels[String(level.id)];
            const unlocked = isLevelUnlocked(level.id, progress);
            const state = item.passed
              ? "passed"
              : !unlocked
                ? "locked"
                : selected === level.id
                  ? "active"
                  : item.attempts > 0
                    ? "progress"
                    : "available";
            return (
              <button
                key={level.id}
                className={`level-step ${state} ${level.stretch ? "stretch" : ""}`}
                onClick={() => onSelect(level)}
                aria-current={selected === level.id ? "step" : undefined}
                aria-label={`Level ${level.id}, ${level.title}, ${state}`}
              >
                <span className="step-node">
                  {item.passed ? (
                    <Check />
                  ) : !unlocked ? (
                    <LockKeyhole />
                  ) : (
                    String(level.id).padStart(2, "0")
                  )}
                </span>
                <span className="step-copy">
                  <strong>{level.title}</strong>
                  <small>{level.topic}</small>
                </span>
                {level.stretch && <Sparkles className="step-stretch" />}
              </button>
            );
          })}
        </nav>
        <div className="rail-actions">
          <label className="practice-card">
            <span>
              <strong>Practice mode</strong>
              <small>Unlock all nine levels</small>
            </span>
            <input
              type="checkbox"
              checked={progress.settings.practiceMode}
              onChange={(event) => onPractice(event.target.checked)}
            />
            <span className="toggle" />
          </label>
          <button className="summary-button" onClick={onSummary}>
            <ListChecks /> Final summary <ArrowRight />
          </button>
        </div>
      </aside>
    </>
  );
}

function InstructionsPanel({
  level,
  hintsUsed,
  onUseHint,
}: {
  level: Level;
  hintsUsed: number[];
  onUseHint: (index: number) => void;
}) {
  return (
    <details className="instructions panel">
      <summary className="instruction-summary">
        <span>
          <em>Exercise brief</em>
          <strong>{level.instructions.objective}</strong>
        </span>
        <code>{level.signature}</code>
        <span className="brief-action">
          Read full specification <ChevronDown />
        </span>
      </summary>
      <div className="instruction-content">
        <p className="starter-note">{level.instructions.starterNote}</p>
        <div className="contract-grid">
          <InfoBlock title="Input">{level.instructions.input}</InfoBlock>
          <InfoBlock title="Expected output">
            {level.instructions.output}
          </InfoBlock>
        </div>
        <div className="details-body">
          <h4>Contract</h4>
          <p>{level.instructions.contract}</p>
          <h4>Constraints</h4>
          <ul>
            {level.instructions.constraints.map((item) => (
              <li key={item}>{item}</li>
            ))}
          </ul>
          <div className="examples-grid">
            {level.instructions.examples.map((example, index) => (
              <div className="example-card" key={`${example.input}-${index}`}>
                <span>Example {index + 1}</span>
                <code>IN&nbsp; {example.input}</code>
                <code>OUT {example.output}</code>
              </div>
            ))}
          </div>
          <p className="comparison-note">
            <strong>Comparison:</strong> {level.instructions.whitespaceRules}
          </p>
          <div className="rules-columns">
            <div>
              <h4>Allowed Go built-ins</h4>
              <code className="allowlist">
                {level.instructions.allowedBuiltins.join(", ")}
              </code>
              <h4>Allowed standard-library packages</h4>
              <code className="allowlist">
                {level.instructions.allowedPackages.join(", ")}
              </code>
            </div>
            <div>
              <h4>Also allowed</h4>
              <ul>
                {level.instructions.allowed.map((item) => (
                  <li key={item}>{item}</li>
                ))}
              </ul>
            </div>
          </div>
          <h4>Disallowed</h4>
          <ul>
            {level.instructions.disallowed.map((item) => (
              <li key={item}>{item}</li>
            ))}
          </ul>
          <p className="console-note">
            <Terminal />
            <span>
              <strong>Debug output:</strong> use <code>fmt.Println(...)</code>,{" "}
              <code>print(...)</code>, or <code>println(...)</code>.{" "}
              <code>console.log</code> belongs to JavaScript and is not valid Go.
              Output appears in the Console after the program compiles.
            </span>
          </p>
          <h4>Official Go references</h4>
          <div className="doc-links">
            {level.instructions.documentation.map((link) => (
              <a key={link.url} href={link.url} target="_blank" rel="noreferrer">
                {link.label} <ArrowRight />
              </a>
            ))}
          </div>
        </div>
        <div className="instruction-disclosures">
          <details>
            <summary>
              <AlertTriangle /> Common pitfalls <ChevronDown />
            </summary>
            <ul>
              {level.instructions.commonPitfalls.map((pitfall) => (
                <li key={pitfall}>{pitfall}</li>
              ))}
            </ul>
          </details>
          <details>
            <summary>
              <Lightbulb /> Hints ({hintsUsed.length}/
              {level.instructions.hints.length} opened) <ChevronDown />
            </summary>
            <ol className="hints-list">
              {level.instructions.hints.map((hint, index) => {
                const used = hintsUsed.includes(index);
                return (
                  <li key={hint}>
                    {used ? (
                      <p>{hint}</p>
                    ) : (
                      <button onClick={() => onUseHint(index)}>
                        Reveal hint {index + 1}
                      </button>
                    )}
                  </li>
                );
              })}
            </ol>
          </details>
        </div>
      </div>
    </details>
  );
}

function InfoBlock({ title, children }: { title: string; children: ReactNode }) {
  return (
    <div className="info-block">
      <span>{title}</span>
      <p>{children}</p>
    </div>
  );
}

function Metric({ label, value }: { label: string; value: string }) {
  return (
    <div className="metric">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function ResultsPanel({
  level,
  result,
  running,
  runnerMode,
  expanded,
  showAll,
  onToggle,
  onShowAll,
  onRunOne,
}: {
  level: Level;
  result?: RunResult;
  running: boolean;
  runnerMode: RunnerConfig["runnerMode"];
  expanded: Record<string, boolean>;
  showAll: boolean;
  onToggle: (id: string) => void;
  onShowAll: () => void;
  onRunOne: (id: string) => void;
}) {
  const byID = new Map(result?.results.map((item) => [item.id, item]));
  const counts = result?.results.reduce(
    (total, item) => {
      total[item.status] = (total[item.status] ?? 0) + 1;
      return total;
    },
    {} as Record<string, number>,
  );
  return (
    <div className="results-panel panel">
      <div className="panel-toolbar results-toolbar">
        <div className="panel-title">
          <FlaskConical />
          <strong>Visible tests</strong>
          <span className="test-total">{level.tests.length} required</span>
        </div>
        <button className="text-button" onClick={onShowAll}>
          {showAll ? "Collapse details" : "Show all details"}
        </button>
      </div>
      {running && (
        <div className="running-banner" aria-live="polite">
          <span className="runner-pulse" />
          {runnerMode === "docker"
            ? "Compiling and running in a fresh Docker sandbox…"
            : "Compiling and running with the trusted local runner…"}
        </div>
      )}
      {result && (
        <div className="run-summary">
          <strong>
            {result.passedCount}/{result.totalCount} passed
          </strong>
          <span>{result.durationMs.toFixed(0)} ms total</span>
          {counts?.compile ? <em>{counts.compile} compile</em> : null}
          {counts?.runtime ? <em>{counts.runtime} runtime</em> : null}
          {counts?.assertion ? <em>{counts.assertion} assertion</em> : null}
          {result.timedOut ? <em>timeout</em> : null}
        </div>
      )}
      <ConsolePanel result={result} running={running} />
      <div className="test-list">
        {level.tests.map((test) => (
          <TestCase
            key={test.id}
            test={test}
            result={byID.get(test.id)}
            expanded={showAll || expanded[test.id]}
            disabled={running}
            onToggle={() => onToggle(test.id)}
            onRun={() => onRunOne(test.id)}
          />
        ))}
      </div>
    </div>
  );
}

function ConsolePanel({
  result,
  running,
}: {
  result?: RunResult;
  running: boolean;
}) {
  const hasOutput = Boolean(
    result?.compileError ||
      result?.runtimeError ||
      result?.stdout ||
      result?.stderr,
  );
  return (
    <details className="console-panel" open={running || hasOutput}>
      <summary>
        <span>
          <Terminal /> Console
        </span>
        <small>
          {running
            ? "Running…"
            : result?.compileError
              ? "Compilation failed"
              : hasOutput
                ? "Output available"
                : "No output"}
        </small>
        <ChevronDown />
      </summary>
      <div className="console-body" aria-live="polite">
        {!result && !running && (
          <p>
            Run the tests to see compiler diagnostics and output from{" "}
            <code>fmt.Println</code>, <code>print</code>, or <code>println</code>.
          </p>
        )}
        {running && <p>Waiting for the Go runner…</p>}
        {result?.compileError && (
          <ConsoleStream label="compiler" tone="error">
            {result.compileError}
          </ConsoleStream>
        )}
        {result?.runtimeError && (
          <ConsoleStream label="runtime" tone="error">
            {result.runtimeError}
          </ConsoleStream>
        )}
        {result?.stdout && (
          <ConsoleStream label="stdout">{result.stdout}</ConsoleStream>
        )}
        {result?.stderr && !result.compileError && (
          <ConsoleStream label="stderr" tone="error">
            {result.stderr}
          </ConsoleStream>
        )}
        {result && !hasOutput && (
          <p>The program completed without writing to stdout or stderr.</p>
        )}
      </div>
    </details>
  );
}

function ConsoleStream({
  label,
  tone = "",
  children,
}: {
  label: string;
  tone?: "error" | "";
  children: string;
}) {
  return (
    <section className={`console-stream ${tone}`}>
      <span>{label}</span>
      <pre>{children}</pre>
    </section>
  );
}

function TestCase({
  test,
  result,
  expanded,
  disabled,
  onToggle,
  onRun,
}: {
  test: VisibleTest;
  result?: TestResult;
  expanded: boolean;
  disabled: boolean;
  onToggle: () => void;
  onRun: () => void;
}) {
  const status = result?.status ?? "pending";
  return (
    <article className={`test-case ${status}`}>
      <div className="test-case-head">
        <button className="test-expand" onClick={onToggle}>
          <span className="result-mark">
            {result?.passed ? (
              <Check />
            ) : status === "pending" ? (
              <Circle />
            ) : (
              <X />
            )}
          </span>
          <span>
            <strong>{test.name}</strong>
            <small>{test.purpose}</small>
          </span>
          {result && <em>{result.durationMs.toFixed(2)} ms</em>}
          <ChevronDown className={expanded ? "rotated" : ""} />
        </button>
        <button
          className="run-one"
          onClick={onRun}
          disabled={disabled}
          title={`Run ${test.name}`}
          aria-label={`Run only ${test.name}`}
        >
          <Play />
        </button>
      </div>
      {expanded && (
        <div className="test-detail">
          <ValueRow label="Input" value={test.input} />
          <ValueRow label="Expected" value={test.expected} tone="expected" />
          {result && (
            <ValueRow
              label="Actual"
              value={result.actual || result.failure || "No result"}
              tone={result.passed ? "actual-pass" : "actual-fail"}
            />
          )}
          {result && !result.passed && result.actual && (
            <div className="diff-view">
              <span>Readable diff</span>
              <pre>
                <del>- {test.expected}</del>
                {"\n"}
                <ins>+ {result.actual}</ins>
              </pre>
            </div>
          )}
          {result?.failure && (
            <p className="test-failure">{result.failure}</p>
          )}
        </div>
      )}
    </article>
  );
}

function ValueRow({
  label,
  value,
  tone = "",
}: {
  label: string;
  value: string;
  tone?: string;
}) {
  return (
    <div className={`value-row ${tone}`}>
      <span>{label}</span>
      <code>{value}</code>
    </div>
  );
}

function HistoryPanel({
  entries,
  onClose,
}: {
  entries: SavedProgress["levels"][string]["history"];
  onClose: () => void;
}) {
  return (
    <section className="history-panel panel">
      <div className="history-heading">
        <div>
          <History />
          <strong>Recent runs</strong>
        </div>
        <button
          className="icon-button"
          onClick={onClose}
          aria-label="Close test history"
        >
          <X />
        </button>
      </div>
      {entries.length ? (
        <div className="history-list">
          {entries.map((entry) => (
            <div key={entry.id}>
              <span>{new Date(entry.at).toLocaleString()}</span>
              <strong>{entry.outcome}</strong>
              <em>{entry.durationMs.toFixed(0)} ms</em>
            </div>
          ))}
        </div>
      ) : (
        <p className="empty-copy">No runs yet. Your last ten runs appear here.</p>
      )}
    </section>
  );
}

function Summary({
  levels,
  progress,
  onClose,
  onSelect,
  onExport,
}: {
  levels: Level[];
  progress: SavedProgress;
  onClose: () => void;
  onSelect: (level: Level) => void;
  onExport: () => void;
}) {
  const passed = Object.values(progress.levels).filter(
    (item) => item.passed,
  ).length;
  const sequential = sequentialCompleted(progress);
  const spent = Object.values(progress.levels).reduce(
    (sum, item) => sum + item.timeSpentSeconds,
    0,
  );
  const review = levels.filter(
    (level) => !progress.levels[String(level.id)].passed,
  );
  return (
    <main className="summary-page">
      <button className="summary-close" onClick={onClose}>
        <ArrowLeft /> Return to workspace
      </button>
      <section className="summary-hero">
        <div className={`summary-seal ${passed >= 7 ? "achieved" : ""}`}>
          {passed >= 7 ? <Trophy /> : <ListChecks />}
        </div>
        <p className="eyebrow">Practice assessment summary</p>
        <h2>
          {passed >= 7 ? "Assessment target achieved." : "Your path to 7/9."}
        </h2>
        <p>
          {passed >= 7
            ? "You passed at least seven server-verified levels. Stretch work is still available whenever you want it."
            : `${7 - passed} more passed level${7 - passed === 1 ? "" : "s"} will reach the assessment target.`}
        </p>
        <div className="summary-stats">
          <Metric label="Highest sequential" value={`${sequential}/9`} />
          <Metric label="Total passed" value={`${passed}/9`} />
          <Metric label="Time spent" value={formatSpent(spent)} />
          <Metric
            label="Target"
            value={passed >= 7 ? "Achieved" : "In progress"}
          />
        </div>
        <button className="test-button export-summary" onClick={onExport}>
          <FileJson /> Export progress
        </button>
      </section>

      <section className="summary-table panel">
        <div className="summary-table-head">
          <span>Level</span>
          <span>Status</span>
          <span>Best tests</span>
          <span>Attempts</span>
          <span>Time</span>
          <span>Hints</span>
        </div>
        {levels.map((level) => {
          const item = progress.levels[String(level.id)];
          return (
            <button key={level.id} onClick={() => onSelect(level)}>
              <span>
                <em>{String(level.id).padStart(2, "0")}</em>
                <strong>{level.title}</strong>
                <small>{level.topic}</small>
              </span>
              <span className={item.passed ? "summary-pass" : "summary-open"}>
                {item.passed ? <Check /> : <Circle />}
                {item.passed ? "Passed" : "Review"}
              </span>
              <span>
                {item.bestPassed}/{item.totalTests}
              </span>
              <span>{item.attempts}</span>
              <span>{formatSpent(item.timeSpentSeconds)}</span>
              <span>
                {item.hintsUsed.length}/{level.instructions.hints.length}
              </span>
            </button>
          );
        })}
      </section>

      {review.length > 0 && (
        <section className="review-card">
          <div>
            <p className="eyebrow">Suggested review</p>
            <h3>Curriculum areas to revisit</h3>
          </div>
          <div className="review-tags">
            {review.map((level) => (
              <button key={level.id} onClick={() => onSelect(level)}>
                Level {level.id} · {level.topic}
              </button>
            ))}
          </div>
        </section>
      )}
    </main>
  );
}

export default App;
