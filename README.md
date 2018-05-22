Pomfu is a Go library for uploading data to various file-servers
supporting the [Pomf standard][]. 

Examples on how to use Pomfu, can be found in the `cmd/pomfu` (a pomf
upload client) and `cmd/ppomf` (a pomf proxy server). How to use these
can be read about in their respective man pages.

All software that uses Pomfu, shares a common configuration file. This
lets a user add new servers. How to use this feature is documented in
the `pomfu.5` man page, which can be found in this directory.

pomfu was entirely written from scratch, and is in the public
domain. More detailed documentation and general information can be found
on [Pomfu homepage][].

[Pomf standard]: https://github.com/pomf/pomf-standard
[Pomfu homepage]: http://sub.god.jp/~xat/pomfu/
