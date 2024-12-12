# frozen_string_literal: true

desc 'Starts an IRB console with the myapp environment loaded'
task :console do
  require 'irb'
  require './myapp'
  ARGV.clear
  IRB.start
end
