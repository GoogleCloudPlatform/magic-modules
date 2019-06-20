# Arguments to 'apply_patches' are:
# - name of patch directory
# - commit message
# - author
# - target branch
function apply_patches {
  # Apply necessary downstream patches.
  shopt -s nullglob
  for patch in "$1"/*; do
    # This is going to apply the patch as at least 1 commit, possibly more.
    git am --3way --signoff "$patch"
  done
  shopt -u nullglob
  # Now, collapse the patch commits into one.
  # This looks a little silly, but here's what we're doing.
  # We get rid of all the commits since we diverged from 'master',
  # We keep all the changes (--soft).
  git reset --soft "$(git merge-base HEAD "$4")"
  # Then we commit again.
  git commit -m "$2" --author="$3" --signoff || true  # don't crash if no changes
}
