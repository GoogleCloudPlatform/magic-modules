## Step 1: Delegate Workspace Syncing
Identify the absolute path of the downstream workspace we are working in. It should be the only other active workspace that is not the Magic Modules repo. If ambiguous, ask the user for the downstream path.
Once you have the path, you must delegate the synchronization task.

**Action:** Use the `invoke_subagent` tool to call the `session-setup` subagent.
- **Prompt to send:** "Please initialize and sync the workspace located at the following path: `<insert_absolute_path_here>`"
- Wait for the subagent to return its final completion message.