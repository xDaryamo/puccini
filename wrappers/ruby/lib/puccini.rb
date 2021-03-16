require 'fiddle'
require 'fiddle/import'
require 'yaml'

module Puccini
  extend Fiddle::Importer
  dlload File.join(__dir__, 'libpuccini.so')
  extern 'char *Compile(char *, char *)'

  module TOSCA
    extend self

    class Problems < StandardError
      def initialize(problems)
        @problems = problems
      end
      attr_reader :problems
    end

    def compile(url, inputs: {})
      inputs = YAML.dump inputs
      result = YAML.load Puccini::Compile(url, inputs).to_s
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
