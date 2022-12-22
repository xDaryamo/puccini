#!/usr/bin/env ruby

require 'optparse'
require 'puccini'
require 'yaml'

inputs = {}
quirks = []

option_parser = OptionParser.new do |opts|
  opts.on '-i', '--input=INPUT', 'specify input (format is name=value)' do |i|
    k, v = i.split('=')
    inputs[k] = YAML.load(v)
  end
  opts.on '-q', '--quirk=QUIRK', 'specify quirk' do |q|
    quirks.push(q)
  end
end

option_parser.parse!

if ARGV.length == 0
  puts 'no URL provided'
  exit 1
end

begin
  clout = Puccini::TOSCA.compile(ARGV[0], inputs, quirks)
  puts YAML.dump(clout)
rescue Puccini::TOSCA::Problems => e
  puts 'Problems:'
  for problem in e.problems
    puts YAML.dump(problem)
  end
  exit 1
end
