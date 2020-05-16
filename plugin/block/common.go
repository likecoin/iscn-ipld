package block

// IPLD Codecs for ISCN
// See the authoritative document:
// https://github.com/multiformats/multicodec/blob/master/table.csv
const (
	CodecISCN         = 0x0264
	CodecRights       = 0x0265
	CodecStakeholders = 0x0266
	CodecEntity       = 0x0267
	CodecContent      = 0x0268
)

// IsIscnObject checks the codec whether belongs an ISCN object
func IsIscnObject(codec uint64) bool {
	switch codec {
	case CodecISCN, CodecRights, CodecStakeholders, CodecEntity, CodecContent:
		return true
	default:
		return false
	}
}
