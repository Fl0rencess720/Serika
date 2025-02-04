package protocol

import "testing"

func Test_HeaderMashall(t *testing.T) {
	header := Header{
		CompressType:  1,
		Status:        0,
		ServiceMethod: "love_taffy",
		ServicePath:   "meowars",
		ID:            1883,
		PayloadLen:    1023,
		Checksum:      666,
	}
	res := header.Mashall()
	t.Errorf("Marshalled Header: %v", res)
}

func Test_Unmashall(t *testing.T) {
	data := []byte{8, 0, 1, 0, 10, 108, 111, 118, 101, 95, 116, 97, 102, 102, 121, 219, 14, 255, 7, 154, 2, 0, 0}

	var header Header
	err := header.Unmashall(data)
	if err != nil {
		t.Fatalf("Unmashall failed: %v", err)
	}
	expectedHeader := Header{
		MagicNumber:   magicNumber,
		Status:        0,
		CompressType:  1,
		ServiceMethod: "love_taffy",
		ServicePath:   "meowars",
		ID:            1883,
		PayloadLen:    1023,
		Checksum:      666,
	}

	if header != expectedHeader {
		t.Errorf("Unmashall result mismatch\nExpected: %+v\nGot: %+v", expectedHeader, header)
	} else {
		t.Logf("Unmashall success: %+v", header)
	}
}
