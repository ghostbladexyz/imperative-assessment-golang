import { describe, expect, it } from "vitest";
import { constrainDragTop, reorderTests } from "./test-order";
import type { VisibleTest } from "./types";

const tests = ["first", "second", "third"].map(
  (id): VisibleTest => ({
    id,
    name: id,
    purpose: "",
    input: "",
    expected: "",
  }),
);

describe("reorderTests", () => {
  it("moves a dragged test to the target position without mutating the catalogue", () => {
    const reordered = reorderTests(tests, "third", "first");

    expect(reordered.map((test) => test.id)).toEqual([
      "third",
      "first",
      "second",
    ]);
    expect(tests.map((test) => test.id)).toEqual([
      "first",
      "second",
      "third",
    ]);
  });

  it("keeps the order unchanged when either identifier is unknown", () => {
    expect(reorderTests(tests, "missing", "first")).toBe(tests);
    expect(reorderTests(tests, "first", "missing")).toBe(tests);
  });
});

describe("constrainDragTop", () => {
  it("does not let the final card separate from the preceding card", () => {
    expect(constrainDragTop(520, 100, 600, 120, 350)).toBe(350);
  });

  it("stays inside the viewport and returns whole pixels", () => {
    expect(constrainDragTop(90.4, 100, 600, 120)).toBe(100);
    expect(constrainDragTop(429.6, 100, 600, 120)).toBe(430);
    expect(constrainDragTop(700, 100, 600, 120)).toBe(480);
    expect(constrainDragTop(200, 100, 600, 120, 80)).toBe(100);
  });
});
