package tape

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldCheckIfTapeExist(t *testing.T) {
	Convey("should check if a tape exist for a system", t, func() {
		tape := new(Tape)
		tape.Path = "./"
		tape.SystemID = "1223"
		So(tape.exist(), ShouldBeFalse)
		os.Mkdir("./1223", os.ModePerm)
		So(tape.exist(), ShouldBeTrue)
		os.RemoveAll("./1223")
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
		tape, err := NewTape("1123", "./")
		So(err, ShouldBeNil)
		evt := domain.Event{
			Name: "event",
		}
		err = tape.RecordEvent(&evt)
		So(err, ShouldBeNil)
		So(len(tape.Segments), ShouldEqual, 1)
		fd, _ := os.Open("./1123")
		infos, _ := fd.Readdir(10)
		So(len(infos), ShouldEqual, 2)

		os.RemoveAll("./1123")

	})

	Convey("should save reader on tape", t, func() {
		tape, err := NewTape("1234", "./")
		So(err, ShouldBeNil)
		evt := domain.Event{
			Name: "event",
		}

		d, _ := json.Marshal(evt)
		r := bytes.NewReader(d)
		err = tape.RecordReader("dump.txt", "dump", r)
		So(err, ShouldBeNil)
		So(len(tape.Segments), ShouldEqual, 1)
		fd, _ := os.Open("./1234")
		infos, _ := fd.Readdir(10)
		So(len(infos), ShouldEqual, 2)
		os.RemoveAll("./1234")
	})

	Convey("should close tape", t, func() {
		tape, err := NewTape("a1234", "./")
		So(err, ShouldBeNil)
		evt := domain.Event{
			Name: "event",
		}
		d, _ := json.Marshal(evt)
		r := bytes.NewReader(d)
		err = tape.RecordReader("dump.txt", "dump", r)
		So(err, ShouldBeNil)
		So(tape.Close(), ShouldBeNil)
		fd1, _ := os.Open("./")
		names, _ := fd1.Readdirnames(-1)

		exist := false
		fName := ""
		for _, name := range names {
			if strings.HasSuffix(name, ".rec") && strings.HasPrefix(name, tape.SystemID) {
				exist = true
				fName = name
			}
		}
		So(exist, ShouldBeTrue)
		os.Remove(fName)
		So(tape.exist(), ShouldBeFalse)
	})
}
