package tape

import (
	"os"
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldCheckIfTapeExist(t *testing.T) {
	Convey("should check if a tape exist for a system", t, func() {
		tape := new(Tape)
		tape.Path = "./"
		tape.SystemID = "123"
		So(tape.exist(), ShouldBeFalse)
		os.Mkdir("./123", os.ModePerm)
		So(tape.exist(), ShouldBeTrue)
		os.RemoveAll("./123")
	})
}

func TestShouldGetOrCreateTape(t *testing.T) {
	Convey("should get or create a new tape", t, func() {
		tape, err := NewTape("123", "./")
		So(err, ShouldBeNil)
		So(tape.exist(), ShouldBeTrue)
		fd, err := os.OpenFile("./123/tape.json", os.O_RDONLY, os.ModePerm)
		So(err, ShouldBeNil)
		if err == nil {
			fd.Close()
		}
		os.RemoveAll("./123")
	})
}

func TestShouldSaveSegments(t *testing.T) {
	Convey("should save segments on tape", t, func() {
		tape, err := NewTape("123", "./")
		So(err, ShouldBeNil)
		evt := domain.Event{
			Name: "event",
		}
		err = tape.RecordEvent(&evt)
		So(err, ShouldBeNil)
	})
}
