import type {
  Catalogue,
  ExerciseKey,
  RunnerConfig,
  RunResult,
} from "./types";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
  });
  const data = (await response.json()) as T & { error?: string };
  if (!response.ok) {
    throw new Error(data.error ?? `Request failed with status ${response.status}`);
  }
  return data;
}

export async function fetchCatalogue(): Promise<Catalogue> {
  return request<Catalogue>("/api/levels");
}

export async function fetchRunnerConfig(): Promise<RunnerConfig> {
  return request<RunnerConfig>("/api/health");
}

export async function runTests(
  exerciseKey: ExerciseKey,
  code: string,
  testIds: string[],
  signal: AbortSignal,
): Promise<RunResult> {
  return request<RunResult>("/api/run", {
    method: "POST",
    body: JSON.stringify({ exerciseKey, code, testIds }),
    signal,
  });
}

export async function formatCode(code: string): Promise<string> {
  const response = await request<{ code: string }>("/api/format", {
    method: "POST",
    body: JSON.stringify({ code }),
  });
  return response.code;
}

export async function validateReceipts(
  receipts: Record<ExerciseKey, string>,
): Promise<ExerciseKey[]> {
  const response = await request<{ validExerciseKeys: ExerciseKey[] }>(
    "/api/receipts/validate",
    {
      method: "POST",
      body: JSON.stringify({ receipts }),
    },
  );
  return response.validExerciseKeys;
}
