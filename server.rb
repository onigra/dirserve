#! /usr/bin/env ruby

require 'webrick'
require 'optparse'

options = {
  port: 8000,
  directory: '.dirserve'
}

OptionParser.new do |opts|
  opts.banner = "Usage: server.rb [options]"

  opts.on("-p", "--port PORT", Integer, "Specify port number") do |port|
    options[:port] = port
  end

  opts.on("-d", "--directory DIR", "Specify server root directory") do |dir|
    options[:directory] = dir
  end
end.parse!

server = WEBrick::HTTPServer.new(
  Port: options[:port],
  DocumentRoot: options[:directory]
)

server.config[:DirectoryIndex] = ['response.json', 'index.html']

trap('INT') { server.shutdown }
server.start
