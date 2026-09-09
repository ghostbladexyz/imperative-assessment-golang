import { describe, expect, it } from "vitest";
import {
  createProgress,
  isLevelUnlocked,
  loadProgress,
  saveProgress,
  sequentialCompleted,
  settleTimer,
  validateImport,
} from "./storage";
import type { Catalogue, Level } from "./types";

const slugs = ["validate-stack", "safe-sum", "tetris"];
const levels = slugs.map((slug, index) => ({
  key: `checkpoint/${slug}`,
  id: index + 1,
  track: "core",
  trackPosition: index + 1,
  title: slug,
  subject: `# ${slug}`,
  tests: [{ id: "a" }],
  starterCode: "package main\nfunc main() {}\n",
  instructions: { hints: [] },
})) as unknown as Level[];

const catalogue: Catalogue = {
  levels,
  progressSchemaVersion: 6,
  legacyProgress: { schemaVersion: 4, exerciseKeys: [] },
};

function installLocalStorage(): () => void {
  const previous = Object.getOwnPropertyDescriptor(globalThis, "localStorage");
  const values = new Map<string, string>();
  const storage = {
    clear: () => values.clear(),
    getItem: (key: string) => values.get(key) ?? null,
    key: (index: number) => [...values.keys()][index] ?? null,
    removeItem: (key: string) => values.delete(key),
    setItem: (key: string, value: string) => values.set(key, value),
    get length() {
      return values.size;
    },
  } as Storage;
  Object.defineProperty(globalThis, "localStorage", {
    configurable: true,
    value: storage,
  });
  return () => {
    if (previous) {
      Object.defineProperty(globalThis, "localStorage", previous);
    } else {
      Reflect.deleteProperty(globalThis, "localStorage");
    }
  };
}

describe("progress state", () => {
  it("unlocks checkpoint exercises sequentially", () => {
    const progress = createProgress(catalogue);
    expect(isLevelUnlocked(levels[0].key, levels, progress)).toBe(true);
    expect(isLevelUnlocked(levels[1].key, levels, progress)).toBe(false);
    progress.exercises[levels[0].key].passed = true;
    expect(isLevelUnlocked(levels[1].key, levels, progress)).toBe(true);
  });

  it("practice mode unlocks every exercise", () => {
    const progress = createProgress(catalogue);
    progress.settings.practiceMode = true;
    expect(isLevelUnlocked(levels[2].key, levels, progress)).toBe(true);
  });

  it("counts only the highest sequential completion", () => {
    const progress = createProgress(catalogue);
    progress.exercises[levels[0].key].passed = true;
    progress.exercises[levels[2].key].passed = true;
    expect(sequentialCompleted(levels, progress)).toBe(1);
  });

  it("settles a persisted running timer", () => {
    const timer = { durationSeconds: 100, elapsedSeconds: 10, running: true, lastTickAt: 1000 };
    expect(settleTimer(timer, 6000).elapsedSeconds).toBe(15);
  });

  it("reconciles removed and newly added exercise keys", () => {
    const saved = createProgress(catalogue);
    saved.exercises[levels[0].key].code = "preserved";
    saved.exercises["removed/old"] = saved.exercises[levels[0].key];
    const added = { ...levels[0], key: "checkpoint/new", id: 4 };
    const changed = { ...catalogue, levels: [...levels, added] };
    const reconciled = validateImport(saved, changed);
    expect(reconciled.exercises[levels[0].key].code).toBe("preserved");
    expect(reconciled.exercises["removed/old"]).toBeUndefined();
    expect(reconciled.exercises[added.key].code).toBe(added.starterCode);
  });

  it("refreshes untouched code when the starter changes", () => {
    const saved = createProgress(catalogue);
    const changedLevel = { ...levels[0], starterCode: "package main\nfunc main() { /* revised */ }\n" };
    const reconciled = validateImport(saved, { ...catalogue, levels: [changedLevel, ...levels.slice(1)] });
    expect(reconciled.exercises[changedLevel.key].code).toBe(changedLevel.starterCode);
  });

  it("preserves learner edits when the starter changes", () => {
    const saved = createProgress(catalogue);
    saved.exercises[levels[0].key].code = "package main\nfunc main() { /* learner */ }\n";
    const changedLevel = { ...levels[0], starterCode: "package main\nfunc main() { /* revised */ }\n" };
    const reconciled = validateImport(saved, { ...catalogue, levels: [changedLevel, ...levels.slice(1)] });
    expect(reconciled.exercises[changedLevel.key].code).toContain("learner");
  });

  it("preserves saved edits when checkpoint hints are omitted", () => {
    const checkpointCatalogue = {
      ...catalogue,
      levels: levels.map((level) => ({
        ...level,
        instructions: { ...level.instructions, hints: null },
      })),
    } as unknown as Catalogue;
    const saved = createProgress(checkpointCatalogue);
    saved.exercises[levels[0].key].code += "\n// learner";
    saved.exercises[levels[0].key].receipt = "receipt";
    saved.exercises[levels[0].key].passed = true;

    const reconciled = validateImport(saved, checkpointCatalogue);

    expect(reconciled.exercises[levels[0].key].code).toContain("learner");
    expect(reconciled.exercises[levels[0].key].passed).toBe(true);
  });

  it("restores edited code and a later exercise from saved local progress", () => {
    const restoreLocalStorage = installLocalStorage();
    try {
      const saved = createProgress(catalogue);
      saved.currentExerciseKey = levels[1].key;
      saved.trackExerciseKeys.core = levels[1].key;
      saved.exercises[levels[0].key].code = "package main\n// learner edit\n";

      saveProgress(saved);

      const reloaded = loadProgress(catalogue);
      expect(reloaded.currentExerciseKey).toBe(levels[1].key);
      expect(reloaded.trackExerciseKeys.core).toBe(levels[1].key);
      expect(reloaded.exercises[levels[0].key].code).toContain("learner edit");
    } finally {
      restoreLocalStorage();
    }
  });

  it("rejects a foreign schema", () => {
    expect(() => validateImport({ schemaVersion: 99 }, catalogue)).toThrow(/schema version 6 or 4/);
  });
});
