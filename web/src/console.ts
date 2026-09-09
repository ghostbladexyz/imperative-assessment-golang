import type { Level, RunResult } from "./types";
import { displayTestStatus } from "./test-status";

export type ConsoleStream = {
  label: string;
  value: string;
  error: boolean;
};

export function buildConsoleStreams(
  result: RunResult | undefined,
  tests: Level["tests"],
): ConsoleStream[] {
  if (!result) return [];

  const programOutput = [
    result.stdout,
    result.stderr && !result.compileError ? result.stderr : "",
  ]
    .filter(Boolean)
    .join("\n");
  const streams: ConsoleStream[] = [];

  if (programOutput) {
    streams.push({ label: "output", value: programOutput, error: false });
  }
  if (result.compileError) {
    streams.push({
      label: "compiler",
      value: result.compileError,
      error: true,
    });
    return streams;
  }
  if (result.runtimeError) {
    streams.push({
      label: "runtime",
      value: result.runtimeError,
      error: true,
    });
    return streams;
  }

  const states = tests.map((test, index) => {
    const testResult = result.results.find((item) => item.id === test.id);
    return {
      index,
      test,
      testResult,
      status: displayTestStatus(result, testResult),
    };
  });
  const failed = states.find((state) => state.status === "fail");
  if (failed) {
    streams.push({
      label: `test #${failed.index + 1} failed: ${failed.test.name}`,
      value: [
        `Input: ${failed.testResult?.input ?? failed.test.input}`,
        `Need: ${failed.test.expected}`,
        `Got: ${failed.testResult?.failure || failed.testResult?.actual || ""}`,
      ].join("\n"),
      error: true,
    });
  } else {
    const errored = states.find((state) => state.status === "error");
    if (errored) {
      streams.push({
        label: `test #${errored.index + 1} error: ${errored.test.name}`,
        value: `Error: ${errored.testResult?.failure || "The official grader did not return a result for this check."}`,
        error: true,
      });
    } else {
      const notRun = states.find((state) => state.status === "not_run");
      if (notRun) {
        streams.push({
          label: `test #${notRun.index + 1} not run: ${notRun.test.name}`,
          value:
            notRun.testResult?.failure ||
            "The official grader stopped before this check ran.",
          error: false,
        });
      }
    }
  }
  if (states.length > 0 && states.every((state) => state.status === "pass")) {
    streams.push({
      label: "tests",
      value: `All ${states.length} tests passed.`,
      error: false,
    });
  }

  return streams;
}
