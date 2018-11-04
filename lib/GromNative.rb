require 'ffi'
require "GromNative/version"

module GromNative
  extend FFI::Library
  ffi_lib File.expand_path("../ext/native.so", File.dirname(__FILE__))
  attach_function :get, [:string], :string
end
