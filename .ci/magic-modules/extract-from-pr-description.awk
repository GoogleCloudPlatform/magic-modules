#! /usr/bin/gawk
BEGIN {
# This tells us that each "record" is separated by 5 or more dashes.
# That keeps the description of the magic-modules PR separate from
# the submodules.  It tells us that the "fields" are separated by at
# least one '#' symbol - this means that the PR request is two
# records, the second of which contains one field per title.
  RS = "-{5,}\n"; FS = "#+ ";
# This takes the command line argument 'fetch', a comma-separated
# list of which fields to fetch, and splits it into the array
# 'tofetch'.
  split(fetch, tofetch, ",")
}

{ 
  if (NR == 1) {
    # The first record is the magic-modules portion.  Skip it.
    next
  }
  # Here we are hunting for the item which has the earliest position
  # in 'tofetch'.  We iterate through the 'fields'...
  x = 2; best_num = 99999; best = ""
  while (x <= NF) {
    # Extract the tag (anything between []s)...
    match($x, /\[([-a-z]+)\]/, title)
    pref = 1
    # Iterate through the list of tags we're willing to fetch
    while (pref < length(tofetch) ) {
      if (title[1] == tofetch[pref] && pref < best_num) {
        # And, if it's higher-positioned than the best one
        # we have found so far, we store it.
        best = substr($x, index($x, "\n") + 1)
        best_num = pref
      }
      pref++
    }
    x++
  }
  # At the end, we print the string following the tag which appears
  # first in the 'fetch' comma-separated-list.  This may be an
  # empty string, and the caller will need to handle that.
  print best
}
