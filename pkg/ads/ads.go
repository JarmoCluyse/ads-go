package ads

import (
	"encoding/binary"

	"github.com/jarmocluyse/ads-go/pkg/ads/constants"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
	"github.com/jarmocluyse/ads-go/pkg/ads/utils"
)

// Struct for Ams address and port
type AmsAddress struct {
	NetID string // ams net id
	Port  uint16 // port number
}

// createAmsTcpHeader creates the AMS/TCP header.
func createAmsTcpHeader(command types.AMSHeaderFlag, dataLength uint32) []byte {
	buf := make([]byte, constants.AMSTCPHeaderLength)
	binary.LittleEndian.PutUint16(buf[0:2], uint16(command)) // Ams command
	binary.LittleEndian.PutUint32(buf[2:6], dataLength)      // length of the data
	return buf
}

// createAmsHeader creates the AMS header.
func createAmsHeader(target AmsAddress, source AmsAddress, command types.ADSCommand, dataLength uint32, invokeID uint32) ([]byte, error) {
	buf := make([]byte, constants.AMSHeaderLength)
	targetNetID, err := utils.AmsNetIdStrToByteArray(target.NetID)
	if err != nil {
		return nil, err
	}
	copy(buf[0:6], targetNetID)                          // Add Target Address
	binary.LittleEndian.PutUint16(buf[6:8], target.Port) // Add target port

	sourceNetID, err := utils.AmsNetIdStrToByteArray(source.NetID)
	if err != nil {
		return nil, err
	}
	copy(buf[8:14], sourceNetID)                           // Add Source Address
	binary.LittleEndian.PutUint16(buf[14:16], source.Port) // Add Source port

	binary.LittleEndian.PutUint16(buf[16:18], uint16(command))                      // ADS Command to send
	binary.LittleEndian.PutUint16(buf[18:20], uint16(types.ADSStateFlagAdsCommand)) // indicate we are sending an ADS command see state flags
	binary.LittleEndian.PutUint32(buf[20:24], dataLength)                           // length of the data
	binary.LittleEndian.PutUint32(buf[24:28], 0)                                    // ErrorCode (normally 0 while sending)
	binary.LittleEndian.PutUint32(buf[28:32], invokeID)                             // invoke id to send
	return buf, nil
}
