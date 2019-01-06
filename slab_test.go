package gos

import (
	"reflect"
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSlab(t *testing.T) {
	newSlab(5, 10)
}

func TestSlabBitset(t *testing.T) {
	Convey("When creating a new slab", t, func() {
		objSize := uint8(5)
		objsPerSlab := uint(10000)
		slab := newSlab(objSize, objsPerSlab)
		So(slab.objSize, ShouldEqual, objSize)
		So(slab.objsPerSlab(), ShouldEqual, objsPerSlab)
		Convey("we should be able to obtain and use the bitset from it", func() {
			bitSet1 := slab.bitSet()
			bitSet1.Set(2)
			bitSet1.Set(4)
			bitSet1.Set(9)
			Convey("then we can obtain it a second time", func() {
				bitSet2 := slab.bitSet()
				So(bitSet2.Test(0), ShouldBeFalse)
				So(bitSet2.Test(1), ShouldBeFalse)
				So(bitSet2.Test(2), ShouldBeTrue)
				So(bitSet2.Test(3), ShouldBeFalse)
				So(bitSet2.Test(4), ShouldBeTrue)
				So(bitSet2.Test(5), ShouldBeFalse)
				So(bitSet2.Test(6), ShouldBeFalse)
				So(bitSet2.Test(7), ShouldBeFalse)
				So(bitSet2.Test(8), ShouldBeFalse)
				So(bitSet2.Test(9), ShouldBeTrue)
			})
		})
	})
}

func TestConversion(t *testing.T) {
	type myType struct {
		a     uint8
		value uint64
	}
	myVar1 := myType{a: 1, value: 12345}

	var copyFrom []byte
	copyFromHeader := (*reflect.SliceHeader)(unsafe.Pointer(&copyFrom))
	copyFromHeader.Data = uintptr(unsafe.Pointer(&myVar1))
	copyFromHeader.Cap = 16
	copyFromHeader.Len = 16

	copyTo := make([]byte, len(copyFrom))

	for i := range copyFrom {
		copyTo[i] = copyFrom[i]
	}

	myVar2 := (*myType)(unsafe.Pointer(&copyFrom[0]))
	myVar3 := (*myType)(unsafe.Pointer(&copyTo[0]))

	if myVar2.value != myVar3.value {
		t.Fatalf("Expected myVar3.value to be %d, but it is %d", myVar2.value, myVar3.value)
	}
}