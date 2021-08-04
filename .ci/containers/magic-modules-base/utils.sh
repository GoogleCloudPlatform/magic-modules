# Usage: grep_files_modified grep_pattern
# Must be called while inside a downstream git repository.
# Requires 2 commits of depth.
# Returns 0 if there are matches and 1 otherwise.
function grep_files_modified {
	matching_files=$(git diff --name-only HEAD~1 | { grep "${1}" || test $? = 1; })
	if [[ -z $matching_files ]]; then
		return 1
	else
		return 0
	fi
}

# Usage: update_status context state target_url
# Expected env: GITHUB_TOKEN
function update_status {
	local post_body=$( jq -n \
		--arg context "${1}" \
		--arg state "${2}" \
		--arg target_url "${3}" \
		'{context: $context, target_url: $target_url, state: $state}')
	echo "Updating status ${1} to ${2} with target_url ${3} for sha ${mm_commit_sha}"
	curl \
	  -X POST \
	  -u "modular-magician:$GITHUB_TOKEN" \
	  -H "Accept: application/vnd.github.v3+json" \
	  "https://api.github.com/repos/GoogleCloudPlatform/magic-modules/statuses/${mm_commit_sha}" \
	  -d "$post_body"
}