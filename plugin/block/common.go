package block

// IPLD Codecs for ISCN
// See the authoritative document:
// https://github.com/multiformats/multicodec/blob/master/table.csv
const (
	CodecISCN    = 0x0264
	CodecContent = 0x0777
)

// IsIscnObject checks the codec whether belongs an ISCN object
func IsIscnObject(codec uint64) bool {
	switch codec {
	case CodecISCN, CodecContent:
		return true
	default:
		return false
	}
}
