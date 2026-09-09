import { describe, expect, it } from "vitest";
import { isMermaidSource } from "./MermaidDiagram";

describe("isMermaidSource", () => {
  it("recognizes the unlabeled flowcharts used by exercise briefs", () => {
    expect(
      isMermaidSource(`
        flowchart LR
        Args --> Valid
      `),
    ).toBe(true);
  });

  it("leaves ordinary fenced code blocks alone", () => {
    expect(isMermaidSource("package main\nfunc main() {}")).toBe(false);
  });
});
