require "bundler/setup"
require "GromNative"
require 'webmock/rspec'

$rack_pid = nil

RSpec.configure do |config|
  # Enable flags like --only-failures and --next-failure
  config.example_status_persistence_file_path = ".rspec_status"

  # Disable RSpec exposing methods globally on `Module` and `main`
  config.disable_monkey_patching!

  config.expect_with :rspec do |c|
    c.syntax = :expect
  end

  # Configure test api lifecycle
  config.before :suite do
    command = 'cd spec/support; rackup -p 3333 --pid rack_server.pid'
    sleep_seconds = 3

    puts '-- Start API'

    fork do
      puts "  -- Running: #{command}"

      exec command
    end

    puts "  -- Sleeping for #{sleep_seconds} seconds so that server can come up"
    sleep sleep_seconds

    puts '  -- DONE'
  end

  config.after :suite do
    command = 'cd spec/support;if [ -f rack_server.pid ]; then kill `cat rack_server.pid`; fi'
    sleep_seconds = 2

    puts ''
    puts '-- Stop API'
    puts "  -- Running: #{command}"

    system(command)

    puts "  -- Sleeping for #{sleep_seconds} seconds so that server can shut down"
    sleep sleep_seconds

    puts '  -- DONE'
  end
end

Dir[File.expand_path(File.join(File.dirname(__FILE__),'support','**','*.rb'))].each {|f| require f}