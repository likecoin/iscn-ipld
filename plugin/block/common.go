package block

import "fmt"

// IPLD Codecs for ISCN
// See the authoritative document:
// https://github.com/multiformats/multicodec/blob/master/table.csv
const (
	CodecISCN         = 0x0264
	CodecRights       = 0x0265
	CodecStakeholders = 0x0266
	CodecContent      = 0x0267
	CodecEntity       = 0x0268

	// Internal Codec
	CodecRight       = 0x02BD
	CodecStakeholder = 0x02D1
	CodecTimePeriod  = 0x033F
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

// ValidateParent between version and parent CID
func ValidateParent(version *Number, parent *Cid) error {
	ver, err := version.GetUint64()
	if err != nil {
		return err
	}

	if ver == 1 {
		if parent.IsDefined() {
			return fmt.Errorf("Parent should not be set as version <= 1")
		}
	} else if ver > 1 {
		if !parent.IsDefined() {
			return fmt.Errorf("Parent missed as version > 1")
		}
	}

	return nil
}
