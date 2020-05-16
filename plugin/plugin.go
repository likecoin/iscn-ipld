package plugin

import (
	"fmt"

	"github.com/ipfs/go-ipfs/core/coredag"
	"github.com/ipfs/go-ipfs/plugin"
	"github.com/likecoin/iscn-ipld/plugin/block"
	"github.com/likecoin/iscn-ipld/plugin/block/content"
	"github.com/likecoin/iscn-ipld/plugin/block/kernel"

	ipld "github.com/ipfs/go-ipld-format"
)

// Plugins is exported list of plugins that will be loaded
var Plugins = []plugin.Plugin{
	&Plugin{},
}

// ==================================================
// Plugin
// ==================================================

// Plugin is the main structure.
type Plugin struct{}

// Static (compile time) check that Plugin satisfies the plugin.PluginIPLD interface.
var _ plugin.PluginIPLD = (*Plugin)(nil)

// Name returns the name of Plugin
func (*Plugin) Name() string {
	return "ipld-iscn"
}

// Version returns the version of Plugin
func (*Plugin) Version() string {
	return "0.5.0.0.0"
}

// Init Plugin
func (*Plugin) Init(*plugin.Environment) error {
	fmt.Println("ISCN IPLD plugin loaded")
	kernel.Register()
	content.Register()
	return nil
}

// RegisterBlockDecoders registers the decoder for different types of block
func (*Plugin) RegisterBlockDecoders(decoder ipld.BlockDecoder) error {
	decoder.Register(block.CodecISCN, kernel.BlockDecoder)
	decoder.Register(block.CodecContent, content.BlockDecoder)
	return nil
}

// RegisterInputEncParsers registers the encode parsers needed to put the blocks into the DAG
func (*Plugin) RegisterInputEncParsers(encodingParsers coredag.InputEncParsers) error {
	return nil
}
