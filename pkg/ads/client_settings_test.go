package ads

import (
	"testing"
	"time"
)

func TestClientSettings_LoadDefaults_AllZero(t *testing.T) {
	cs := &ClientSettings{}
	cs.LoadDefaults()

	if cs.TargetNetID != "127.0.0.1.1.1" {
		t.Errorf("expected TargetNetID default, got %q", cs.TargetNetID)
	}
	if cs.RouterAddr != "127.0.0.1" {
		t.Errorf("expected RouterAddr default, got %q", cs.RouterAddr)
	}
	if cs.RouterPort != 48898 {
		t.Errorf("expected RouterPort default, got %d", cs.RouterPort)
	}
	if cs.Timeout != 2*time.Second {
		t.Errorf("expected Timeout default, got %v", cs.Timeout)
	}
}

func TestClientSettings_LoadDefaults_PartialSet(t *testing.T) {
	cs := &ClientSettings{
		TargetNetID: "1.2.3.4.5.6",
		Timeout:     5 * time.Second,
	}
	cs.LoadDefaults()

	if cs.TargetNetID != "1.2.3.4.5.6" {
		t.Errorf("should not overwrite set TargetNetID, got %q", cs.TargetNetID)
	}
	if cs.RouterAddr != "127.0.0.1" {
		t.Errorf("expected RouterAddr default, got %q", cs.RouterAddr)
	}
	if cs.RouterPort != 48898 {
		t.Errorf("expected RouterPort default, got %d", cs.RouterPort)
	}
	if cs.Timeout != 5*time.Second {
		t.Errorf("should not overwrite set Timeout, got %v", cs.Timeout)
	}
}

func TestClientSettings_LoadDefaults_AllSet(t *testing.T) {
	cs := &ClientSettings{
		TargetNetID: "2.2.2.2.2.2",
		RouterAddr:  "192.168.99.1",
		RouterPort:  12345,
		Timeout:     10 * time.Second,
	}
	cs.LoadDefaults()

	if cs.TargetNetID != "2.2.2.2.2.2" {
		t.Errorf("should not overwrite set TargetNetID, got %q", cs.TargetNetID)
	}
	if cs.RouterAddr != "192.168.99.1" {
		t.Errorf("should not overwrite set RouterAddr, got %q", cs.RouterAddr)
	}
	if cs.RouterPort != 12345 {
		t.Errorf("should not overwrite set RouterPort, got %d", cs.RouterPort)
	}
	if cs.Timeout != 10*time.Second {
		t.Errorf("should not overwrite set Timeout, got %v", cs.Timeout)
	}
}
