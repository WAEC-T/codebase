# frozen_string_literal: true

task :console do
  require "irb"
  require "./myapp"
  ARGV.clear
  IRB.start
end
