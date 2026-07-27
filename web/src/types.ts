export type Theme = "light" | "dark";
export type RunnerMode = "docker" | "local";
export type ExerciseKey = string;

export interface RunnerConfig {
  ok: boolean;
  runnerMode: RunnerMode;
  sandboxReady: boolean;
  goVersion: string;
  dockerImage?: string;
  message: string;
}

export interface DocumentationLink {
  label: string;
  url: string;
}

export interface Example {
  input: string;
  output: string;
}

export interface Instructions {
  objective: string;
  contract: string;
  input: string;
  output: string;
  constraints: string[];
  examples: Example[];
  documentation: DocumentationLink[];
  allowedBuiltins: string[];
  allowedPackages: string[];
  allowed: string[];
  disallowed: string[];
  starterNote: string;
  whitespaceRules: string;
  commonPitfalls: string[];
  hints: string[];
}

export interface VisibleTest {
  id: string;
  name: string;
  purpose: string;
  input: string;
  expected: string;
}

export interface Level {
  key: ExerciseKey;
  id: number;
  title: string;
  topic: string;
  difficulty: string;
  stretch: boolean;
  signature: string;
  starterCode: string;
  instructions: Instructions;
  tests: VisibleTest[];
}

export interface Catalogue {
  levels: Level[];
  progressSchemaVersion: number;
  legacyProgress: {
    schemaVersion: 4;
    exerciseKeys: ExerciseKey[];
  };
}

export type ResultStatus =
  | "pending"
  | "pass"
  | "assertion"
  | "compile"
  | "runtime";

export interface TestResult extends VisibleTest {
  actual: string;
  passed: boolean;
  status: ResultStatus;
  failure?: string;
  durationMs: number;
}

export interface RunResult {
  exerciseKey: ExerciseKey;
  levelId: number;
  passed: boolean;
  passedCount: number;
  totalCount: number;
  compileError?: string;
  runtimeError?: string;
  failureKind?: "capacity" | "cleanup" | "internal" | "output" | "startup";
  timedOut: boolean;
  stopped: boolean;
  stdout: string;
  stderr: string;
  durationMs: number;
  results: TestResult[];
  formattedCode?: string;
  sourceHash: string;
  receipt?: string;
}

export interface HistoryEntry {
  id: string;
  at: number;
  passed: number;
  total: number;
  durationMs: number;
  outcome: string;
}

export interface ExerciseProgress {
  code: string;
  starterSnapshot: string;
  receipt?: string;
  passed: boolean;
  attempts: number;
  timeSpentSeconds: number;
  hintsUsed: number[];
  bestPassed: number;
  totalTests: number;
  history: HistoryEntry[];
}

export interface TimerState {
  durationSeconds: number;
  elapsedSeconds: number;
  running: boolean;
  lastTickAt: number;
}

export interface Settings {
  theme: Theme;
  practiceMode: boolean;
  autoTest: boolean;
  fontSize: number;
  reducedMotion: boolean;
}

export interface SavedProgress {
  schemaVersion: number;
  updatedAt: number;
  currentExerciseKey: ExerciseKey;
  exercises: Record<ExerciseKey, ExerciseProgress>;
  timer: TimerState;
  settings: Settings;
}
