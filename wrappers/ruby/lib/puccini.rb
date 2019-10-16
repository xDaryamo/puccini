require 'fiddle'
require 'fiddle/import'
require 'yaml'

module Puccini
  extend Fiddle::Importer
  dlload 'libpuccini.so'
  extern 'char *Compile(char *)'

  module TOSCA
    extend self
    def compile(url)
      return YAML.load(Puccini::Compile(url).to_s)
    end
  end
end
