import type { Level, RunResult } from "./types";

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

export async function fetchLevels(): Promise<Level[]> {
  const response = await request<{ levels: Level[] }>("/api/levels");
  return response.levels;
}

export async function runTests(
  levelId: number,
  code: string,
  testIds: string[],
  signal: AbortSignal,
): Promise<RunResult> {
  return request<RunResult>("/api/run", {
    method: "POST",
    body: JSON.stringify({ levelId, code, testIds }),
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
  receipts: Record<string, string>,
): Promise<number[]> {
  const response = await request<{ validLevelIds: number[] }>(
    "/api/receipts/validate",
    {
      method: "POST",
      body: JSON.stringify({ receipts }),
    },
  );
  return response.validLevelIds;
}

