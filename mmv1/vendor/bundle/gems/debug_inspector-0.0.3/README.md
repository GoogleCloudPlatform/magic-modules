debug_inspector [![Build Status](https://travis-ci.org/banister/debug_inspector.svg?branch=master)](https://travis-ci.org/banister/debug_inspector)
===============

_A Ruby wrapper for the new MRI 2.0 debug\_inspector API_

The `debug_inspector` C extension and API were designed and built by [Koichi Sasada](https://github.com/ko1), this project
is just a gemification of his work.

**NOTES:**

* This library makes use of the new debug inspector API in MRI 2.0.0, **do not use this library outside of debugging situations**.
* Only works on MRI 2.0. Requiring it on unsupported Rubies will result in a no-op

Usage
-----

```ruby
require 'debug_inspector'

# Open debug context
# Passed `dc' is only active in a block
RubyVM::DebugInspector.open { |dc|
  # backtrace locations (returns an array of Thread::Backtrace::Location objects)
  locs = dc.backtrace_locations

  # you can get depth of stack frame with `locs.size'
  locs.size.times do |i|
    # binding of i-th caller frame (returns a Binding object or nil)
    p dc.frame_binding(i)

    # iseq of i-th caller frame (returns a RubyVM::InstructionSequence object or nil)
    p dc.frame_iseq(i)

    # class of i-th caller frame
    p dc.frame_class(i)
  end
}
```

Contact
-------

Problems or questions contact me at [github](http://github.com/banister)

License
-------

(The MIT License)

Copyright (c) 2012-2013 (John Mair)

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
