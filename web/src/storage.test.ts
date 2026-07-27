import { describe, expect, it } from "vitest";
import {
  createProgress,
  isLevelUnlocked,
  sequentialCompleted,
  settleTimer,
  validateImport,
} from "./storage";
import type { Catalogue, Level } from "./types";

const levels = Array.from({ length: 171 }, (_, index) => ({
  key: index < 21 ? `foundation/${index + 1}` : `piscine/${index + 980}`,
  id: index + 1,
  title: `Exercise ${index + 1}`,
  tests: [{ id: "a" }],
  starterCode: "func solve() {}",
  instructions: { hints: [] },
})) as unknown as Level[];

const catalogue: Catalogue = {
  levels,
  progressSchemaVersion: 5,
  legacyProgress: {
    schemaVersion: 4,
    exerciseKeys: levels.map((level) => level.key),
  },
};

describe("progress state", () => {
  it("unlocks exercises sequentially by Exercise Key", () => {
    const progress = createProgress(catalogue);
    expect(isLevelUnlocked(levels[0].key, levels, progress)).toBe(true);
    expect(isLevelUnlocked(levels[1].key, levels, progress)).toBe(false);
    progress.exercises[levels[0].key].passed = true;
    expect(isLevelUnlocked(levels[1].key, levels, progress)).toBe(true);
  });

  it("practice mode unlocks every exercise", () => {
    const progress = createProgress(catalogue);
    progress.settings.practiceMode = true;
    expect(isLevelUnlocked(levels[170].key, levels, progress)).toBe(true);
  });

  it("counts only the highest sequential completion", () => {
    const progress = createProgress(catalogue);
    progress.exercises[levels[0].key].passed = true;
    progress.exercises[levels[2].key].passed = true;
    expect(sequentialCompleted(levels, progress)).toBe(1);
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

  it("migrates schema-v4 positions through the frozen key map", () => {
    const legacy = {
      schemaVersion: 4,
      currentLevelId: 2,
      levels: {
        "1": { ...createProgress(catalogue).exercises[levels[0].key], code: "first" },
        "2": { ...createProgress(catalogue).exercises[levels[1].key], code: "second" },
      },
      timer: { elapsedSeconds: 12, running: true },
      settings: { theme: "dark" },
    };
    const reordered = {
      ...catalogue,
      levels: [levels[1], levels[0], ...levels.slice(2)].map((level, index) => ({
        ...level,
        id: index + 1,
      })),
    };

    const migrated = validateImport(legacy, reordered);

    expect(migrated.currentExerciseKey).toBe(levels[1].key);
    expect(migrated.exercises[levels[0].key].code).toBe("first");
    expect(migrated.exercises[levels[1].key].code).toBe("second");
    expect(migrated.schemaVersion).toBe(5);
  });

  it("reconciles removed and newly added Exercise Keys", () => {
    const saved = createProgress(catalogue);
    saved.exercises[levels[0].key].code = "preserved";
    saved.exercises["removed/999"] = saved.exercises[levels[0].key];
    const added = {
      ...levels[0],
      key: "foundation/172",
      id: levels.length + 1,
      starterCode: "func Added() {}",
    };
    const changedCatalogue = {
      ...catalogue,
      levels: [...levels, added],
    };

    const reconciled = validateImport(saved, changedCatalogue);

    expect(reconciled.exercises[levels[0].key].code).toBe("preserved");
    expect(reconciled.exercises["removed/999"]).toBeUndefined();
    expect(reconciled.exercises[added.key].code).toBe(added.starterCode);
  });

  it("refreshes an untouched legacy return starter for a print-only exercise", () => {
    const printedLevel = {
      ...levels[21],
      key: "piscine/1016",
      signature: "DescendComb()",
      starterCode: `package main

import "github.com/01-edu/z01"

func DescendComb() {
	// TODO: print the required output with z01.PrintRune.
	_ = z01.PrintRune
}
`,
      instructions: {
        hints: [],
        allowedPackages: ["github.com/01-edu/z01"],
      },
    } as unknown as Level;
    const printedCatalogue = {
      ...catalogue,
      levels: [printedLevel],
      legacyProgress: {
        ...catalogue.legacyProgress,
        exerciseKeys: [printedLevel.key],
      },
    };
    const saved = createProgress(printedCatalogue);
    saved.exercises[printedLevel.key].code = `package main

func DescendComb() string {
	// TODO: implement the checkpoint behavior.
	return ""
}
`;
    delete (
      saved.exercises[printedLevel.key] as unknown as Record<string, unknown>
    ).starterSnapshot;

    const reconciled = validateImport(saved, printedCatalogue);

    expect(reconciled.exercises[printedLevel.key].code).toBe(
      printedLevel.starterCode,
    );
  });

  it("refreshes untouched code when a tracked starter changes", () => {
    const saved = createProgress(catalogue);
    const changedLevel = {
      ...levels[0],
      starterCode: "func solve() { /* revised */ }",
    };
    const changedCatalogue = {
      ...catalogue,
      levels: [changedLevel, ...levels.slice(1)],
    };

    const reconciled = validateImport(saved, changedCatalogue);

    expect(reconciled.exercises[changedLevel.key].code).toBe(
      changedLevel.starterCode,
    );
    expect(reconciled.exercises[changedLevel.key].starterSnapshot).toBe(
      changedLevel.starterCode,
    );
  });

  it("preserves learner edits when a tracked starter changes", () => {
    const saved = createProgress(catalogue);
    saved.exercises[levels[0].key].code = "func solve() { learnerEdit() }";
    const changedLevel = {
      ...levels[0],
      starterCode: "func solve() { /* revised */ }",
    };
    const changedCatalogue = {
      ...catalogue,
      levels: [changedLevel, ...levels.slice(1)],
    };

    const reconciled = validateImport(saved, changedCatalogue);

    expect(reconciled.exercises[changedLevel.key].code).toBe(
      "func solve() { learnerEdit() }",
    );
    expect(reconciled.exercises[changedLevel.key].starterSnapshot).toBe(
      changedLevel.starterCode,
    );
  });

  it("rejects a foreign schema", () => {
    expect(() => validateImport({ schemaVersion: 99 }, catalogue)).toThrow(
      /schema version 5 or 4/,
    );
  });
});
