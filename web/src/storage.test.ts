import { describe, expect, it } from "vitest";
import {
  createProgress,
  isLevelUnlocked,
  sequentialCompleted,
  settleTimer,
  validateImport,
} from "./storage";
import type { Level } from "./types";

const levels = Array.from({ length: 171 }, (_, index) => ({
  id: index + 1,
  title: `Level ${index + 1}`,
  tests: [{ id: "a" }],
  starterCode: "func solve() {}",
  instructions: { hints: [] },
})) as unknown as Level[];

describe("progress state", () => {
  it("unlocks levels sequentially", () => {
    const progress = createProgress(levels);
    expect(isLevelUnlocked(1, progress)).toBe(true);
    expect(isLevelUnlocked(2, progress)).toBe(false);
    progress.levels["1"].passed = true;
    expect(isLevelUnlocked(2, progress)).toBe(true);
  });

  it("practice mode unlocks every level", () => {
    const progress = createProgress(levels);
    progress.settings.practiceMode = true;
    expect(isLevelUnlocked(171, progress)).toBe(true);
  });

  it("counts only the highest sequential completion", () => {
    const progress = createProgress(levels);
    progress.levels["1"].passed = true;
    progress.levels["3"].passed = true;
    expect(sequentialCompleted(progress)).toBe(1);
  });

  it("settles a persisted running timer", () => {
    const timer = {
      durationSeconds: 100,
      elapsedSeconds: 10,
      running: true,
      lastTickAt: 1000,
    };
    expect(settleTimer(timer, 6000).elapsedSeconds).toBe(15);
  });

  it("rejects a foreign schema", () => {
    expect(() => validateImport({ schemaVersion: 99 }, levels)).toThrow(
      /schema version 4/,
    );
  });
});
