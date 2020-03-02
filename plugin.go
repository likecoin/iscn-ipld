package main

import (
	goplugin "github.com/ipfs/go-ipfs/plugin"
	"github.com/likecoin/iscn-ipld/plugin"
)

// Plugins is an exported list of plugins that will be loaded by go-ipfs.
var Plugins = []goplugin.Plugin{
	&plugin.Plugin{},
}
