package plugin

import (
	"log"

	"github.com/ipfs/go-ipfs/core/coredag"
	"github.com/ipfs/go-ipfs/plugin"
	"github.com/likecoin/iscn-ipld/plugin/block"
	"github.com/likecoin/iscn-ipld/plugin/block/iscn"

	ipld "github.com/ipfs/go-ipld-format"
)

// Plugins is exported list of plugins that will be loaded
var Plugins = []plugin.Plugin{
	&Plugin{},
}

// Plugin is the main structure.
type Plugin struct{}

// Static (compile time) check that Plugin satisfies the plugin.PluginIPLD interface.
var _ plugin.PluginIPLD = (*Plugin)(nil)

// Name returns the name of Plugin
func (*Plugin) Name() string {
	log.Println("ipldiscn-Name")
	return "ipld-iscn"
}

// Version returns the version of Plugin
func (*Plugin) Version() string {
	log.Println("ipldiscn-Version")
	return "0.4.23.0.1"
}

// Init Plugin
func (*Plugin) Init() error {
	log.Println("ipldiscn-Init")
	return nil
}

// RegisterBlockDecoders registers the decoder for different types of block
func (*Plugin) RegisterBlockDecoders(decoder ipld.BlockDecoder) error {
	log.Println("ipldiscn-RegisterBlockDecoders")
	decoder.Register(block.CodecISCN, iscn.BlockDecoder)
	return nil
}

// RegisterInputEncParsers registers the encode parsers needed to put the blocks into the DAG
func (*Plugin) RegisterInputEncParsers(encodingParsers coredag.InputEncParsers) error {
	log.Println("ipldiscn-RegisterInputEncParsers")
	return nil
}
