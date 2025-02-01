package protocol

import "testing"

func Test_HeaderMashall(t *testing.T) {
	header := Header{
		CompressType: 1,
		Method:       "love_taffy",
		ID:           1883,
		Len:          1023,
		Checksum:     666,
	}
	res := header.Mashall()
	t.Errorf("Marshalled Header: %v", res)
}

func Test_Unmashall(t *testing.T) {
	data := []byte{8, 1, 0, 10, 108, 111, 118, 101, 95, 116, 97, 102, 102, 121, 219, 14, 255, 7, 154, 2, 0, 0}

	var header Header
	err := header.Unmashall(data)
	if err != nil {
		t.Fatalf("Unmashall failed: %v", err)
	}
	expectedHeader := Header{
		MagicNumber:  magicNumber,
		CompressType: 1,
		Method:       "love_taffy",
		ID:           1883,
		Len:          1023,
		Checksum:     666,
	}

	if header != expectedHeader {
		t.Errorf("Unmashall result mismatch\nExpected: %+v\nGot: %+v", expectedHeader, header)
	} else {
		t.Logf("Unmashall success: %+v", header)
	}
}
