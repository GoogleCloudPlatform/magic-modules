require File.dirname(__FILE__) + '/test_helper'

class BasicTest < MiniTest::Test
  def test_version
    assert(DebugInspector::VERSION)
  end
end
