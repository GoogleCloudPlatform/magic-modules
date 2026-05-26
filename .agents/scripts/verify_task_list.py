#!/usr/bin/env python3
import sys
import os
import re

def main():
    task_path = sys.argv[1] if len(sys.argv) > 1 else os.getenv("TASK_PATH")
    if not task_path:
        print("Error: TASK_PATH not specified.", file=sys.stderr)
        sys.exit(1)

    if not os.path.exists(task_path):
        print(f"Error: task.md not found at {task_path}", file=sys.stderr)
        sys.exit(1)

    with open(task_path, "r") as f:
        content = f.read()

    # Required checklist patterns (must be marked with [x] or [X])
    checks = {
        "GEMINI.md": r"-\s*\[([xX])\]\s*\[MANDATORY\]\s*Read\s+GEMINI\.md",
        "GEMINI_ADD.md/GEMINI_FIX.md": r"-\s*\[([xX])\]\s*\[MANDATORY\]\s*Read\s+GEMINI_(ADD|FIX)\.md",
        "skill: sync-provider": r"-\s*\[([xX])\]\s*\[MANDATORY\]\s*Read\s+skill:\s*sync-provider",
        "skill: tgc-build-skill": r"-\s*\[([xX])\]\s*\[MANDATORY\]\s*Read\s+skill:\s*tgc-build-skill",
        "skill: tgc-run-unit-tests-skill": r"-\s*\[([xX])\]\s*\[MANDATORY\]\s*Read\s+skill:\s*tgc-run-unit-tests-skill",
        "skill: tgc-run-integration-tests-skill": r"-\s*\[([xX])\]\s*\[MANDATORY\]\s*Read\s+skill:\s*tgc-run-integration-tests-skill",
        "skill: resource specialized skill": r"-\s*\[([xX])\]\s*\[MANDATORY\]\s*Read\s+skill:\s*(tgc-add-new-generated-resource-skill|tgc-fix-handwritten-resources-tests-skill)",
        "Go unit tests executed": r"-\s*\[([xX])\]\s*Run\s+Go\s+unit\s+tests"
    }

    errors = []
    for label, pattern in checks.items():
        if not re.search(pattern, content):
            errors.append(f"Mandatory task item '{label}' has not been completed [x] in task.md.")

    if errors:
        print("❌ TASK LIST VERIFICATION FAILED:", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        sys.exit(1)

    print("✅ Task list verification passed: All mandatory entrypoints and skill checks are successfully completed.")
    sys.exit(0)

if __name__ == "__main__":
    main()
