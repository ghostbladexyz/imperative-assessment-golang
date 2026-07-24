# Imperative Go Assessment

This context describes the language used to identify and organize assessment exercises.

## Language

**Exercise Key**:
The immutable, source-qualified identity of an assessment exercise. It is independent of the exercise's title, difficulty, Curriculum Position, and subordinate test identities.
_Avoid_: Level ID, test ID, title-derived identity

**Curriculum Position**:
The learner-facing numeric place of an exercise in the current teaching order. It may change when the curriculum is reordered and is never an exercise's identity.
_Avoid_: Exercise ID

**Exercise Progress**:
The learner's saved work and completion state for one exercise, associated with its Exercise Key.
_Avoid_: Level progress, position-keyed progress
