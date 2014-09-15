gosudar
=======

A very early stage, not at all ready for primetime yet, and at this moment
Mac-only go program that works like ohai. (Currently the name "gohai" is kind of
taken, so right now this is using "gosudar".)

Installation
============

`go get github.com/ctdk/gosudar`

Usage
=====

*Remember:* Currently this only works on Macs.

First, just run `gosudar`. It should print out a bunch of JSON describing your
Mac's processor and kernel. It's a small sample of what ohai would print out.

Next, make sure to create a directory named `/tmp/plugins`. There will be a 
better directory for that later, but this is only far enough along to showcase
its future functionality. Next, as an illustration, download the gist at 
https://gist.github.com/ctdk/5df969333bc4ed811708. However you download it,
you'll end up with a file named `velikij.go`. Put it in a directory named 
`velikij`, and in that directory run `go build`. Take the binary and put it in 
the `/tmp/plugins` directory. Then run `gosudar` again. If all has gone well,
the JSON will include `"zzz": "zuz"` in its output. This is a gosudar plugin.

A More Detailed Explanation
===========================

Go doesn't allow dynamically linking shared libraries into a binary, so to
support plugins gosudar uses sockets and JSON-RPC to communicate with its 
plugins. The `velikij` plugin in the above gist is a template for an extremely
simple gosudar plugin. All the plugin needs to do is connect to the socket at
call the SendInfo RPC call to send the hash of information the plugin built back
to gosudar.

Gosudar uses build tags to build the appropriate system specific functions for
extracting the system's information. Only the Mac one exists right now, but if
there were more each platform would have a file or set of files that would only
be built on that platform. In `gosudar_mac.go`, there are a lot of sysctl calls,
while on Linux those functions would have the same name, but instead read the
information from `/proc`.

To Do
=====

Because there's quite a bit.

* The rest of the information ohai provides.
* Strip out debugging log statements.
* Some configuration options.
* Better plugin directory default.
* Mutex for the info hash.
* Possibly use a unique socket for each run. This socket name would need to be
  passed to the plugins, though.
* Support other platforms, especially Linux.
* Plugins should have a smooth installation process.


This is only a start. It's a long ways from being ready for any real work, I'm
afraid.
