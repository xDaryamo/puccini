require 'fiddle'
require 'fiddle/import'
require 'yaml'

module Puccini
  extend Fiddle::Importer
  dlload File.join(__dir__, 'libpuccini.so')
  extern 'char *Compile(char *, char *, char *, char, char)'

  module TOSCA
    extend self

    class Problems < StandardError
      def initialize(problems)
        @problems = problems
      end
      attr_reader :problems
    end

    def compile(url, inputs=nil, quirks=nil, resolve=true, coerce=true)
      inputs = YAML.dump (inputs || {})
      quirks = YAML.dump (quirks || [])
      result = YAML.unsafe_load Puccini::Compile(url, inputs, quirks, resolve ? 1 : 0, coerce ? 1 : 0).to_s
      if result.key? 'problems'
        raise Problems.new result['problems']
      elsif result.key? 'error'
        raise StandardError.new result['error']
      else
        return result['clout']
      end
    end
  end
end
