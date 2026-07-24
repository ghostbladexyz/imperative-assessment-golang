import type {
  Level,
  LevelProgress,
  SavedProgress,
  Settings,
  TimerState,
} from "./types";

export const STORAGE_KEY = "imperative-go-assessment:progress:v2";
export const ASSESSMENT_SECONDS = 6 * 60 * 60;

export const defaultSettings = (): Settings => ({
  theme:
    typeof window !== "undefined" &&
    window.matchMedia?.("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light",
  practiceMode: false,
  autoTest: false,
  fontSize: 15,
  reducedMotion:
    typeof window !== "undefined" &&
    window.matchMedia?.("(prefers-reduced-motion: reduce)").matches,
});

export const defaultTimer = (): TimerState => ({
  durationSeconds: ASSESSMENT_SECONDS,
  elapsedSeconds: 0,
  running: true,
  lastTickAt: Date.now(),
});

export function makeLevelProgress(level: Level): LevelProgress {
  return {
    code: level.starterCode,
    passed: false,
    attempts: 0,
    timeSpentSeconds: 0,
    hintsUsed: [],
    bestPassed: 0,
    totalTests: level.tests.length,
    history: [],
  };
}

export function createProgress(levels: Level[]): SavedProgress {
  return {
    schemaVersion: 2,
    updatedAt: Date.now(),
    currentLevelId: 1,
    levels: Object.fromEntries(
      levels.map((level) => [String(level.id), makeLevelProgress(level)]),
    ),
    timer: defaultTimer(),
    settings: defaultSettings(),
  };
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

export function validateImport(value: unknown, levels: Level[]): SavedProgress {
  if (!isRecord(value) || value.schemaVersion !== 2) {
    throw new Error("This backup does not use progress schema version 2.");
  }
  if (
    typeof value.currentLevelId !== "number" ||
    value.currentLevelId < 1 ||
    value.currentLevelId > levels.length ||
    !isRecord(value.levels) ||
    !isRecord(value.timer) ||
    !isRecord(value.settings)
  ) {
    throw new Error("The backup is missing required assessment fields.");
  }
  const clean = createProgress(levels);
  clean.currentLevelId = value.currentLevelId;
  for (const level of levels) {
    const candidate = value.levels[String(level.id)];
    if (!isRecord(candidate) || typeof candidate.code !== "string") {
      throw new Error(`The backup has invalid data for level ${level.id}.`);
    }
    const target = clean.levels[String(level.id)];
    target.code = candidate.code.slice(0, 192 * 1024);
    target.receipt =
      typeof candidate.receipt === "string" ? candidate.receipt : undefined;
    target.passed = candidate.passed === true && Boolean(target.receipt);
    target.attempts = finiteNonNegative(candidate.attempts);
    target.timeSpentSeconds = finiteNonNegative(candidate.timeSpentSeconds);
    target.bestPassed = Math.min(
      level.tests.length,
      finiteNonNegative(candidate.bestPassed),
    );
    target.totalTests = level.tests.length;
    target.hintsUsed = Array.isArray(candidate.hintsUsed)
      ? candidate.hintsUsed
          .filter(
            (item): item is number =>
              Number.isInteger(item) &&
              Number(item) >= 0 &&
              Number(item) < level.instructions.hints.length,
          )
          .slice(0, level.instructions.hints.length)
      : [];
    target.history = Array.isArray(candidate.history)
      ? candidate.history
          .filter(isRecord)
          .map((item) => ({
            id: typeof item.id === "string" ? item.id : crypto.randomUUID(),
            at: finiteNonNegative(item.at),
            passed: finiteNonNegative(item.passed),
            total: finiteNonNegative(item.total),
            durationMs: finiteNonNegative(item.durationMs),
            outcome:
              typeof item.outcome === "string" ? item.outcome.slice(0, 80) : "",
          }))
          .slice(0, 10)
      : [];
  }
  clean.timer = {
    durationSeconds: ASSESSMENT_SECONDS,
    elapsedSeconds: Math.min(
      ASSESSMENT_SECONDS,
      finiteNonNegative(value.timer.elapsedSeconds),
    ),
    running: value.timer.running === true,
    lastTickAt: Date.now(),
  };
  const defaults = defaultSettings();
  clean.settings = {
    theme: value.settings.theme === "dark" ? "dark" : "light",
    practiceMode: value.settings.practiceMode === true,
    autoTest: value.settings.autoTest === true,
    fontSize: Math.min(
      22,
      Math.max(12, finiteNonNegative(value.settings.fontSize) || defaults.fontSize),
    ),
    reducedMotion: value.settings.reducedMotion === true,
  };
  clean.updatedAt = Date.now();
  return clean;
}

function finiteNonNegative(value: unknown): number {
  return typeof value === "number" && Number.isFinite(value) && value >= 0
    ? Math.floor(value)
    : 0;
}

export function loadProgress(levels: Level[]): SavedProgress {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? validateImport(JSON.parse(raw), levels) : createProgress(levels);
  } catch {
    return createProgress(levels);
  }
}

export function saveProgress(progress: SavedProgress): void {
  localStorage.setItem(
    STORAGE_KEY,
    JSON.stringify({ ...progress, updatedAt: Date.now() }),
  );
}

export function settleTimer(timer: TimerState, now = Date.now()): TimerState {
  if (!timer.running) {
    return { ...timer, lastTickAt: now };
  }
  const delta = Math.max(0, Math.floor((now - timer.lastTickAt) / 1000));
  const elapsedSeconds = Math.min(
    timer.durationSeconds,
    timer.elapsedSeconds + delta,
  );
  return {
    ...timer,
    elapsedSeconds,
    running: elapsedSeconds < timer.durationSeconds,
    lastTickAt: now,
  };
}

export function sequentialCompleted(progress: SavedProgress): number {
  let completed = 0;
  for (let levelId = 1; levelId <= Object.keys(progress.levels).length; levelId += 1) {
    if (!progress.levels[String(levelId)]?.passed) break;
    completed = levelId;
  }
  return completed;
}

export function isLevelUnlocked(
  levelId: number,
  progress: SavedProgress,
): boolean {
  return (
    progress.settings.practiceMode ||
    levelId === 1 ||
    progress.levels[String(levelId - 1)]?.passed === true
  );
}
