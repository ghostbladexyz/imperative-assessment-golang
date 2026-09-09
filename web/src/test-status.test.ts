import { describe, expect, it } from "vitest";
import { displayTestStatus, displayTestStatusLabel } from "./test-status";
import type { RunResult, TestResult } from "./types";

const result = {
  exerciseKey: "checkpoint/validate-stack",
  levelId: 1,
  passed: false,
  passedCount: 1,
  totalCount: 4,
  timedOut: false,
  stopped: false,
  stdout: "",
  stderr: "",
  durationMs: 1,
  results: [],
  sourceHash: "",
} satisfies RunResult;

function testResult(status: TestResult["status"], passed: boolean): TestResult {
  return {
    id: status,
    name: status,
    purpose: status,
    input: "",
    expected: "",
    actual: "",
    passed,
    status,
    durationMs: 0,
  };
}

describe("displayTestStatus", () => {
  it("keeps mixed executed, unexecuted, and error checks distinct", () => {
    const statuses = [
      displayTestStatus(result, testResult("pass", true)),
      displayTestStatus(result, testResult("assertion", false)),
      displayTestStatus(result, testResult("not_run", false)),
      displayTestStatus(result, testResult("runtime", false)),
    ];

    expect(statuses).toEqual(["pass", "fail", "not_run", "error"]);
    expect(statuses.map(displayTestStatusLabel)).toEqual([
      "PASS",
      "FAIL",
      "NOT RUN",
      "ERROR",
    ]);
  });

  it("treats a missing result as an error after a run exists", () => {
    expect(displayTestStatus(result, undefined)).toBe("error");
    expect(displayTestStatus(undefined, undefined)).toBe("pending");
  });
});
