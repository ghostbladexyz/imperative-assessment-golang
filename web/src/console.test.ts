import { describe, expect, it } from "vitest";
import { buildConsoleStreams } from "./console";
import type { Level, RunResult } from "./types";

const tests: Level["tests"] = [
  {
    id: "empty",
    name: "Empty input",
    purpose: "Handles empty input.",
    input: '""',
    expected: '["hi"]',
  },
];

function run(overrides: Partial<RunResult>): RunResult {
  return {
    levelId: 1,
    passed: false,
    passedCount: 0,
    totalCount: 1,
    timedOut: false,
    stopped: false,
    stdout: "",
    stderr: "",
    durationMs: 1,
    results: [],
    sourceHash: "",
    ...overrides,
  };
}

describe("buildConsoleStreams", () => {
  it("shows a compile error once without duplicating it as a test failure", () => {
    const streams = buildConsoleStreams(
      run({
        compileError: "solution.go:7:14: syntax error",
        results: [
          {
            ...tests[0],
            actual: "",
            passed: false,
            status: "compile",
            failure: "solution.go:7:14: syntax error",
            durationMs: 0,
          },
        ],
      }),
      tests,
    );

    expect(streams).toEqual([
      {
        label: "compiler",
        value: "solution.go:7:14: syntax error",
        error: true,
      },
    ]);
  });

  it("shows program output and only the first failed test", () => {
    const streams = buildConsoleStreams(
      run({
        stdout: "hello",
        results: [
          {
            ...tests[0],
            actual: "[]",
            passed: false,
            status: "assertion",
            durationMs: 0,
          },
        ],
      }),
      tests,
    );

    expect(streams).toEqual([
      { label: "output", value: "hello", error: false },
      {
        label: "test #1 failed: Empty input",
        value: 'Need: ["hi"], Actual: []',
        error: true,
      },
    ]);
  });
});
