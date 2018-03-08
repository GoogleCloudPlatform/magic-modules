BEGIN { RS = "-{5,}\n"; FS = "#+ "; split(fetch, tofetch, ",") }

{ 
  if (NR == 1) {
    next
  }
  x = 2
  best_num = 99999
  best = ""
  while (x <= NF) {
    match($x, /\[([-a-z]+)\]/, title)
    pref = 1
    while (pref < length(tofetch) ) {
      if (title[1] == tofetch[pref] && pref < best_num) {
        best = substr($x, index($x, "\n") + 1)
        best_num = pref
      }
      pref++
    }
    x++
  }
  print best
}
