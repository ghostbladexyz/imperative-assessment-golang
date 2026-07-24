import type { Level, RunResult } from "./types";

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

  const failed = result.results.find((testResult) => !testResult.passed);
  if (failed) {
    const index = tests.findIndex((test) => test.id === failed.id);
    streams.push({
      label: `test #${index + 1} failed: ${failed.name}`,
      value: failed.failure
        ? `Error: ${failed.failure}`
        : `Need: ${failed.expected}, Actual: ${failed.actual}`,
      error: true,
    });
  } else if (result.results.length > 0) {
    streams.push({
      label: "tests",
      value: `All ${result.results.length} tests passed.`,
      error: false,
    });
  }

  return streams;
}
