#!/usr/bin/env ruby

require 'puccini'
require 'yaml'

if ARGV.length == 0
  puts 'no URL provided'
  exit 1
end

begin
  clout = Puccini::TOSCA.compile(ARGV[0])
  puts YAML.dump(clout)
rescue Puccini::TOSCA::Problems => e
  puts 'Problems:'
  for problem in e.problems
    puts YAML.dump(problem)
  end
  exit 1
end
