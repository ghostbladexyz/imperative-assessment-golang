import type {
  Catalogue,
  ExerciseKey,
  ExerciseProgress,
  Level,
  SavedProgress,
  Settings,
  TimerState,
} from "./types";

export const STORAGE_KEY = "imperative-go-assessment:progress";
export const LEGACY_STORAGE_KEY = "imperative-go-assessment:progress:v4";
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

export function makeExerciseProgress(level: Level): ExerciseProgress {
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

export function createProgress(catalogue: Catalogue): SavedProgress {
  const first = catalogue.levels[0];
  if (!first) throw new Error("The catalogue contains no exercises.");
  return {
    schemaVersion: catalogue.progressSchemaVersion,
    updatedAt: Date.now(),
    currentExerciseKey: first.key,
    exercises: Object.fromEntries(
      catalogue.levels.map((level) => [level.key, makeExerciseProgress(level)]),
    ),
    timer: defaultTimer(),
    settings: defaultSettings(),
  };
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

export function validateImport(
  value: unknown,
  catalogue: Catalogue,
): SavedProgress {
  if (!isRecord(value)) {
    throw new Error("The backup is not a progress record.");
  }
  if (value.schemaVersion === catalogue.progressSchemaVersion) {
    return reconcileCurrent(value, catalogue);
  }
  if (value.schemaVersion === catalogue.legacyProgress.schemaVersion) {
    return migrateLegacyProgress(value, catalogue);
  }
  throw new Error(
    `This backup does not use progress schema version ${catalogue.progressSchemaVersion} or ${catalogue.legacyProgress.schemaVersion}.`,
  );
}

function reconcileCurrent(
  value: Record<string, unknown>,
  catalogue: Catalogue,
): SavedProgress {
  if (!isRecord(value.exercises) || !isRecord(value.timer) || !isRecord(value.settings)) {
    throw new Error("The backup is missing required assessment fields.");
  }
  const clean = createProgress(catalogue);
  const requestedKey =
    typeof value.currentExerciseKey === "string"
      ? value.currentExerciseKey
      : clean.currentExerciseKey;
  clean.currentExerciseKey = catalogue.levels.some(
    (level) => level.key === requestedKey,
  )
    ? requestedKey
    : clean.currentExerciseKey;
  reconcileExercises(clean, value.exercises, catalogue.levels, (level) => level.key);
  reconcileSharedState(clean, value);
  return clean;
}

function migrateLegacyProgress(
  value: Record<string, unknown>,
  catalogue: Catalogue,
): SavedProgress {
  if (!isRecord(value.levels) || !isRecord(value.timer) || !isRecord(value.settings)) {
    throw new Error("The schema-v4 backup is missing required assessment fields.");
  }
  const clean = createProgress(catalogue);
  const legacyPosition =
    typeof value.currentLevelId === "number" &&
    Number.isInteger(value.currentLevelId)
      ? value.currentLevelId
      : 1;
  const requestedKey = catalogue.legacyProgress.exerciseKeys[legacyPosition - 1];
  if (
    requestedKey &&
    catalogue.levels.some((level) => level.key === requestedKey)
  ) {
    clean.currentExerciseKey = requestedKey;
  }
  const legacyPositionByKey = new Map(
    catalogue.legacyProgress.exerciseKeys.map((key, index) => [key, index + 1]),
  );
  reconcileExercises(clean, value.levels, catalogue.levels, (level) => {
    const position = legacyPositionByKey.get(level.key);
    return position === undefined ? undefined : String(position);
  });
  reconcileSharedState(clean, value);
  return clean;
}

function reconcileExercises(
  clean: SavedProgress,
  candidates: Record<string, unknown>,
  levels: Level[],
  candidateKey: (level: Level) => string | undefined,
): void {
  for (const level of levels) {
    const key = candidateKey(level);
    const candidate = key === undefined ? undefined : candidates[key];
    if (candidate === undefined) continue;
    if (!isRecord(candidate) || typeof candidate.code !== "string") {
      throw new Error(`The backup has invalid data for exercise ${level.key}.`);
    }
    const target = clean.exercises[level.key];
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
}

function reconcileSharedState(
  clean: SavedProgress,
  value: Record<string, unknown>,
): void {
  const timer = value.timer as Record<string, unknown>;
  clean.timer = {
    durationSeconds: ASSESSMENT_SECONDS,
    elapsedSeconds: Math.min(
      ASSESSMENT_SECONDS,
      finiteNonNegative(timer.elapsedSeconds),
    ),
    running: timer.running === true,
    lastTickAt: Date.now(),
  };
  const settings = value.settings as Record<string, unknown>;
  const defaults = defaultSettings();
  clean.settings = {
    theme: settings.theme === "dark" ? "dark" : "light",
    practiceMode: settings.practiceMode === true,
    autoTest: settings.autoTest === true,
    fontSize: Math.min(
      22,
      Math.max(12, finiteNonNegative(settings.fontSize) || defaults.fontSize),
    ),
    reducedMotion: settings.reducedMotion === true,
  };
  clean.updatedAt = Date.now();
}

function finiteNonNegative(value: unknown): number {
  return typeof value === "number" && Number.isFinite(value) && value >= 0
    ? Math.floor(value)
    : 0;
}

export function loadProgress(catalogue: Catalogue): SavedProgress {
  for (const key of [STORAGE_KEY, LEGACY_STORAGE_KEY]) {
    try {
      const raw = localStorage.getItem(key);
      if (raw) return validateImport(JSON.parse(raw), catalogue);
    } catch {
      // Try the older store before falling back to clean progress.
    }
  }
  return createProgress(catalogue);
}

export function saveProgress(progress: SavedProgress): void {
  localStorage.setItem(
    STORAGE_KEY,
    JSON.stringify({ ...progress, updatedAt: Date.now() }),
  );
  localStorage.removeItem(LEGACY_STORAGE_KEY);
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

export function sequentialCompleted(
  levels: Level[],
  progress: SavedProgress,
): number {
  let completed = 0;
  for (const level of levels) {
    if (!progress.exercises[level.key]?.passed) break;
    completed = level.id;
  }
  return completed;
}

export function isLevelUnlocked(
  exerciseKey: ExerciseKey,
  levels: Level[],
  progress: SavedProgress,
): boolean {
  if (progress.settings.practiceMode) return true;
  const index = levels.findIndex((level) => level.key === exerciseKey);
  return index === 0 || progress.exercises[levels[index - 1]?.key]?.passed === true;
}
