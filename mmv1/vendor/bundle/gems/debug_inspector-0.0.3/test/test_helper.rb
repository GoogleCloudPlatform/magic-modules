lib_dir = File.dirname(__FILE__) + '/../lib'

require 'minitest/autorun'
$:.unshift lib_dir unless $:.include?(lib_dir)
require 'debug_inspector'
