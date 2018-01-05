# This is a poor example that will not pass puppet-lint.
# The arrows are not properly aligned.
file { '/tmp/test':
  ensure => file,   # this should have => aligned with longest
  owner  => 'root', # ditto
  content => 'test content',
}
