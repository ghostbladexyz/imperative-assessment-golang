import type { VisibleTest } from "./types";

export function reorderTests(
  tests: VisibleTest[],
  movingID: string,
  targetID: string,
): VisibleTest[] {
  const movingIndex = tests.findIndex((test) => test.id === movingID);
  const targetIndex = tests.findIndex((test) => test.id === targetID);
  if (movingIndex < 0 || targetIndex < 0 || movingIndex === targetIndex) {
    return tests;
  }

  const reordered = [...tests];
  const [moving] = reordered.splice(movingIndex, 1);
  reordered.splice(targetIndex, 0, moving);
  return reordered;
}

export function constrainDragTop(
  desiredTop: number,
  viewportTop: number,
  viewportBottom: number,
  cardHeight: number,
  previousCardBottom?: number,
): number {
  const viewportMaximum = Math.max(viewportTop, viewportBottom - cardHeight);
  const constrainedMaximum =
    previousCardBottom === undefined
      ? viewportMaximum
      : Math.min(viewportMaximum, previousCardBottom);
  const maximum = Math.max(viewportTop, constrainedMaximum);
  return Math.round(Math.min(Math.max(desiredTop, viewportTop), maximum));
}
