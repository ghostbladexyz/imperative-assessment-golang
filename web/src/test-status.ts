import type { RunResult, TestResult } from "./types";

export type DisplayTestStatus = "pending" | "pass" | "fail" | "not_run" | "error";

export function displayTestStatus(
  result: RunResult | undefined,
  testResult: TestResult | undefined,
): DisplayTestStatus {
  if (!result) return "pending";
  if (!testResult) return "error";
  if (testResult.status === "not_run") return "not_run";
  if (testResult.status === "pass" && testResult.passed) return "pass";
  if (testResult.status === "assertion" && !testResult.passed) return "fail";
  return "error";
}

export function displayTestStatusLabel(status: DisplayTestStatus): string {
  return status === "pass"
    ? "PASS"
    : status === "fail"
      ? "FAIL"
      : status === "not_run"
        ? "NOT RUN"
        : status === "error"
          ? "ERROR"
          : "PENDING";
}
