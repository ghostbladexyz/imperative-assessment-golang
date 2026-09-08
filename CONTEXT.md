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

**Starter Snapshot**:
The exact starter code associated with saved Exercise Progress. It distinguishes untouched starter code from a learner's edits when an exercise contract changes.
_Avoid_: Starter version, default solution

**Answer Mode**:
The way an exercise exposes the value assessed by its tests: either a returned value or exact printed output.
_Avoid_: Test type, output type

**Assessment Track**:
A sequence of exercises that shares a teaching order. The current assessment exposes one Core track.
_Avoid_: Module, global level range

**Official Check**:
A pass/fail condition reported by the digest-pinned Zone01 grader. It may be a Published Fixture, a Seeded Check, a structural or race check, or an aggregate suite.
_Avoid_: Manual test, local harness test

**Check Need**:
The learner-facing condition an Official Check requires, such as exact output, status, state, structure, or race safety. It is distinct from the grader's machine result of pass or fail.
_Avoid_: Expected pass, test result

**Published Fixture**:
A stable input and expectation embedded in the official grader and visible in its diagnostics.
_Avoid_: Hidden test

**Seeded Check**:
An Official Check whose concrete data changes with the grader seed to discourage hard-coded answers. The emitted replay seed identifies a particular run.
_Avoid_: Random test, unknowable test

**Oracle-backed Expectation**:
An expected result produced inside the official grader by a compiled reference program and compared with learner output. The oracle implementation is not duplicated by this project.
_Avoid_: Manual expected value
