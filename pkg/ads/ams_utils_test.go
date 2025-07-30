package ads

import (
	"encoding/binary"
	"testing"

	"github.com/jarmocluyse/ads-go/pkg/ads/constants"
	"github.com/jarmocluyse/ads-go/pkg/ads/types"
)

func Test_createAmsTcpHeader(t *testing.T) {
	command := types.AMSHeaderFlag(1)
	dataLength := uint32(1234)
	buf := createAmsTcpHeader(command, dataLength)

	if len(buf) != constants.AMSTCPHeaderLength {
		t.Fatalf("expected buffer length %d, got %d", constants.AMSTCPHeaderLength, len(buf))
	}

	cmd := binary.LittleEndian.Uint16(buf[0:2])
	if cmd != uint16(command) {
		t.Errorf("expected command %d, got %d", command, cmd)
	}

	lenVal := binary.LittleEndian.Uint32(buf[2:6])
	if lenVal != dataLength {
		t.Errorf("expected dataLength %d, got %d", dataLength, lenVal)
	}
}

func Test_createAmsHeader_validInput(t *testing.T) {
	target := AmsAddress{NetID: "1.2.3.4.5.6", Port: 851}
	source := AmsAddress{NetID: "6.5.4.3.2.1", Port: 852}
	command := types.ADSCommand(2)
	dataLength := uint32(5678)
	invokeID := uint32(42)

	head, err := createAmsHeader(target, source, command, dataLength, invokeID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(head) != constants.AMSHeaderLength {
		t.Fatalf("expected AMS header length %d, got %d", constants.AMSHeaderLength, len(head))
	}
	// Spot-check a few fields
	targetPort := binary.LittleEndian.Uint16(head[6:8])
	if targetPort != target.Port {
		t.Errorf("expected target port %d, got %d", target.Port, targetPort)
	}
	commandField := binary.LittleEndian.Uint16(head[16:18])
	if commandField != uint16(command) {
		t.Errorf("expected command %d, got %d", command, commandField)
	}
	sentInvokeID := binary.LittleEndian.Uint32(head[28:32])
	if sentInvokeID != invokeID {
		t.Errorf("expected invokeID %d, got %d", invokeID, sentInvokeID)
	}
}

func Test_createAmsHeader_invalidNetID(t *testing.T) {
	invalidTarget := AmsAddress{NetID: "bad.netid", Port: 851}
	source := AmsAddress{NetID: "6.5.4.3.2.1", Port: 852}
	_, err := createAmsHeader(invalidTarget, source, 1, 0, 1)
	if err == nil {
		t.Error("expected error for invalid target NetID, got nil")
	}

	validTarget := AmsAddress{NetID: "1.2.3.4.5.6", Port: 851}
	invalidSource := AmsAddress{NetID: "not.a.netid", Port: 852}
	_, err = createAmsHeader(validTarget, invalidSource, 1, 0, 1)
	if err == nil {
		t.Error("expected error for invalid source NetID, got nil")
	}
}
