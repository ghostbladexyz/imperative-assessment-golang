# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [4.0.0](https://github.com/ghostbladexyz/imperative-assessment-golang/compare/v3.0.0...v4.0.0) (2026-09-09)

### Other

- consolidate official assessment workflow ([#178](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/94a2f32b492b1bffa6ea6f7b100e11e7d3cf1930))
- *(review)* Honor grader ordering and explicit amd64 emulation ([#179](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/e21591fb015291e9256cb37b7a20386735162e10))
- *(review)* Remove test reordering and reject failed grader aggregates ([#180](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/fa80b22ff7c2b8b83ef5385d506aacee8fcf87e9))
- *(review)* Enforce canonical runs and show per-test status ([#181](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/ac342287f1515e02e28f90d07f1b0dbacd393ffa))
- *(review)* Mount exercise resources read-only in grader ([#182](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/e9d390682618a50e4c58eb8c7da47fb3e69c6b8e))
- *(review)* Fix resource paths, permissions, and security documentation ([#183](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/91ef0ef02663487cb382420d168effd70dd85beb))
- *(test)* Preserve progress across reloads with null-safe hint reconciliation ([#184](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/ee41e318b2be6970eab3e158f4df784028544acf))
- *(document)* Align domain glossary; lint checks pass ([#185](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/172f93558a28319d6df3432d2cd016112b15a9b9))
- *(review)* Fix test status mapping and regression coverage ([#186](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/4ac8533801ed3d35a0dee59792b870c2595f8075))
- *(review)* Prevent synthetic failures for missing aggregate checks ([#187](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/33912f014f0eec7bf27c2ac2d5600c02b84028e5))
- *(review)* Make grader cache writable and resources learner-readable ([#188](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/65ad0ae422e3c5f5794e800cbb26828b65fcb2bb))
- *(review)* Reject missing grader results and reset incompatible progress ([#189](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/fabe39b8013bda8855ea027cf6a791b59b992d4e))
- *(review)* Require failure evidence before marking checks not run ([#190](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/2d3a68f8ba30732c87d4680213103fb030a677d3))
- *(test)* Fix missing grader outcomes and official compile resource limits ([#191](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/bab015763352fc14ad1a3650aed62ec827f059f6))
- *(document)* Refresh stale checkpoint documentation naming ([#192](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/ecb3a227ffbb1ebacacb63c8c2f90fd4bd47c452))
- *(ci)* Regenerated web/package-lock.json to synchronize the @emnapi dependency graph. Verified npm ci, Go vet/tests, frontend lint/tests/build, gofmt, and git diff checks pass ([#193](https://github.com/ghostbladexyz/imperative-assessment-golang/commit/0e8e4143bb47ab43650a471d528f0a4256419dad))
