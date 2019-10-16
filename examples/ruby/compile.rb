#!/usr/bin/env ruby

require 'puccini'
require 'yaml'

if ARGV.length == 0
  puts 'no URL provided'
  exit 1
end

clout = Puccini::TOSCA.compile(ARGV[0])

puts YAML.dump(clout)
