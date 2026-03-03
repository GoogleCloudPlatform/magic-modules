import subprocess
import os

def run_command(command):
    print(f"Running: {command}")
    result = subprocess.run(command, shell=True, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    return result.stdout.strip()

def clean_file(filepath):
    with open(filepath, 'r') as f:
        lines = f.readlines()
    
    new_lines = []
    last_line_empty = False
    for line in lines:
        is_empty = not line.strip()
        if is_empty:
            if not last_line_empty:
                new_lines.append(line)
            last_line_empty = True
        else:
            new_lines.append(line)
            last_line_empty = False
            
    # Remove trailing empty lines
    while new_lines and not new_lines[-1].strip():
        new_lines.pop()
    
    # Ensure one newline at end if file is not empty
    if new_lines and not new_lines[-1].endswith('\n'):
        new_lines[-1] += '\n'
        
    with open(filepath, 'w') as f:
        f.writelines(new_lines)

def main():
    branches = ["tgc-ai-compute", "refactor-compute-part-2", "refactor-compute-part-3", "refactor-compute-part-4"]
    
    for branch in branches:
        print(f"Processing branch: {branch}")
        run_command(f"git checkout {branch}")
        
        # Get modified files in the last commit (or HEAD if we just committed)
        # We want to check files that were touched in the 'Refactor...' commit.
        # Since we might have multiple commits now (e.g. the merge/split), let's look at the files modified in the last commit.
        files = run_command("git show --name-only --format='' HEAD").splitlines()
        files = [f for f in files if f.strip() and f.endswith('.yaml')]
        
        if not files:
            print(f"No YAML files found in last commit of {branch}")
            continue
            
        print(f"Checking {len(files)} files...")
        for f in files:
            if os.path.exists(f):
                clean_file(f)
        
        # Check if any changes
        status = run_command("git status --porcelain")
        if status:
            print(f"Changes found in {branch}, committing...")
            run_command("git add .")
            run_command("git commit -m 'Remove empty lines'")
            run_command(f"git push origin {branch}")
        else:
            print(f"No changes in {branch}")

if __name__ == "__main__":
    main()
