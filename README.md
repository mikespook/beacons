Beacons
=======

[![Build Status][travis-img]][travis]

Beacons is a application that handling and passing data between 
systems, scripted by Lua and written in [Golang][golang].

Dependency
==========

 * [mikespook/golib][golib]
 * [aarzilli/golua][golua]
 * [stevedonovan/luar][luar]
 * [mgo/bson][bson]
 * liblua5.1-0-dev for Ubuntu

Installing
==========

All useful scripts were put at the directory [shell][shell].

Befor building, the proper lua librarie must be installed.
E.g. Ubuntu 14.04, it is `liblua5.1-0-dev`.

Then:

	go get github.com/mikespook/beacons/beacons

The beacons library can be embeded into other projects.

	go get github.com/mikespook/beacons

Have a fun!

Usage
=====

Service
-------

Executing following command:

	$ ${fullpath}/beacons -h

Some help information:

	Usage of ./beacons:
		-config="config.json": Configration file

Scripting
---------

Beacons uses Lua as the scripting language. Beacons will pass the log data into
a Lua script for handling and passing to next.

Authors
=======

 * Xing Xing <mikespook@gmail.com> [Blog][blog] [@Twitter][twitter]

Open Source
===========

See LICENSE for more information.

[golang]: http://golang.org
[golib]: https://github.com/mikespook/golib
[golua]: https://github.com/aarzilli/golua
[luar]: https://github.com/stevedonovan/luar
[blog]: http://mikespook.com
[twitter]: http://twitter.com/mikespook
[travis-img]: https://travis-ci.org/mikespook/beacons.png?branch=master
[travis]: https://travis-ci.org/mikespook/beacons
[shell]: https://github.com/mikespook/beacons/tree/master/shell 
[bson]: http://labix.org/v2/mgo/bson
