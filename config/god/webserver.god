God.watch do |w|
  w.name  = "webserver"
  w.dir   = "/go/src/github.com/garethstokes/fourtyeight/"
  w.uid   = "ubuntu"
  w.start = "/go/src/github.com/garethstokes/fourtyeight/run_webserver.sh"
  w.log   = "/var/log/god/webserver.log"
  w.keepalive(:memory_max => 512.megabytes)
  w.env = {
    'GOPATH' => '/go/'
  }
end
