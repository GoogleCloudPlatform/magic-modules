import subprocess
import os
import shutil

repo_root = "/Users/zhenhuali/Documents/workspace/container_node_pool"
keep_prefix = "mmv1/third_party/tgc_next/pkg/services/container"
target_dir = "mmv1/third_party/tgc_next"

def run_command(command):
    subprocess.run(command, shell=True, check=True, cwd=repo_root)

def cleanup():
    # Get status
    result = subprocess.run(["git", "status", "--porcelain"], cwd=repo_root, capture_output=True, text=True)
    lines = result.stdout.splitlines()

    files_to_restore = []
    files_to_clean = []

    for line in lines:
        if not line.strip():
            continue
        status = line[:2]
        path = line[3:]
        
        # Only care about files in mmv1/third_party/tgc_next
        if not path.startswith(target_dir):
            continue

        # Skip if in the keep directory
        if path.startswith(keep_prefix):
            continue

        if status.strip() == "??":
            files_to_clean.append(path)
        else:
            files_to_restore.append(path)

    # Restore modified/deleted/staged files
    if files_to_restore:
        print(f"Restoring {len(files_to_restore)} files...")
        # Restore in batches to avoid command line length limits
        batch_size = 50
        for i in range(0, len(files_to_restore), batch_size):
            batch = files_to_restore[i:i+batch_size]
            cmd = ["git", "restore", "--staged", "--worktree"] + batch
            subprocess.run(cmd, cwd=repo_root, check=True)

    # Remove untracked files
    if files_to_clean:
        print(f"Cleaning {len(files_to_clean)} files...")
        for path in files_to_clean:
            full_path = os.path.join(repo_root, path)
            if os.path.isdir(full_path):
                shutil.rmtree(full_path)
            elif os.path.isfile(full_path):
                os.remove(full_path)

if __name__ == "__main__":
    cleanup()
