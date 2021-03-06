
lib = File.expand_path('../lib', __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'grom_native/version'

Gem::Specification.new do |spec|
  spec.name          = 'grom_native'
  spec.version       = GromNative::VERSION
  spec.authors       = ['Matt Rayner']
  spec.email         = ['m@rayner.io']

  spec.summary       = %q{A GROM implementation using a GO native extension for heavy lifting.}
  spec.homepage      = 'https://github.com/ukparliament/gromnative'
  spec.license       = 'MIT'

  # Specify which files should be added to the gem when it is released.
  # The `git ls-files -z` loads the files in the RubyGem that have been added into git.
  spec.files         = Dir.chdir(File.expand_path('..', __FILE__)) do
    `git ls-files -z`.split("\x0").reject { |f| f.match(%r{^(test|spec|features)/}) }
  end
  spec.bindir        = 'exe'
  spec.executables   = spec.files.grep(%r{^exe/}) { |f| File.basename(f) }
  spec.require_paths = ['lib']

  spec.add_dependency 'ffi', '~> 1.9'
  spec.add_dependency 'grom'
  spec.add_dependency 'rdf', '~> 3.0'

  spec.add_development_dependency 'bundler', '~> 1.16'
  spec.add_development_dependency 'parliament-grom-decorators'
  spec.add_development_dependency 'rack', '~> 2.0'
  spec.add_development_dependency 'rake', '~> 10.0'
  spec.add_development_dependency 'rspec', '~> 3.0'
  spec.add_development_dependency 'webmock', '~> 3.4'
end
