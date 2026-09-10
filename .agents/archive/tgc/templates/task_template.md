# Tasks - TGC Resource Workflow

## Phase 1: Session Setup
- [ ] [MANDATORY] Read AGENTS.md and TGC_WORKFLOWS.md <!-- id: 1 -->
- [ ] [MANDATORY] Read tgc_add.md or tgc_fix.md <!-- id: 2 -->
- [ ] [MANDATORY] Read skill: tgc-sync-provider <!-- id: 3 -->
- [ ] [MANDATORY] Read skill: tgc-build-skill <!-- id: 4 -->
- [ ] [MANDATORY] Read skill: tgc-run-unit-tests-skill <!-- id: 5 -->
- [ ] [MANDATORY] Read skill: tgc-run-integration-tests-skill <!-- id: 6 -->
- [ ] [MANDATORY] Read skill: tgc-add-new-generated-resource-skill or tgc-fix-handwritten-resources-tests-skill <!-- id: 7 -->
- [ ] Verify environment setup and TGC_DIR variable <!-- id: 8 -->
- [ ] Run the tgc-sync-provider skill to align repositories <!-- id: 9 -->
- [ ] Obtain User Approval for the Implementation Plan <!-- id: 10 -->

## Phase 2: Implementation
- [ ] Add or modify resource definition in Magic Modules <!-- id: 11 -->
- [ ] Ensure correct field ordering in YAML configuration <!-- id: 12 -->
- [ ] Build TGC (run build_tgc.sh) <!-- id: 13 -->

## Phase 3: Unit Testing
- [ ] Run selective unit tests (run_changed_folders_tests.sh) <!-- id: 14 -->

## Phase 4: Integration Testing
- [ ] [MANDATORY] Inspect generated downstream Go test files to identify exact test function names <!-- id: 15 -->
- [ ] Run integration tests with WRITE_FILES=true <!-- id: 16 -->
- [ ] Verify generated test files exist and not all are skipped <!-- id: 17 -->

## Phase 5: Fix Failures (If any)
- [ ] Apply necessary decoders/flatteners or skip rules in MMv1 <!-- id: 18 -->
- [ ] Verify test roundtrip files (Test_roundtrip.tf, Test_roundtrip.json, Test_export.tf) <!-- id: 19 -->

## Phase 6: Final Verification
- [ ] Commit green changes (mmv1/ folder only) <!-- id: 20 -->
