package tape

import (
	"bufio"
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
		So(tape.Exist(), ShouldBeFalse)
		os.Mkdir("./1223", os.ModePerm)
		So(tape.Exist(), ShouldBeTrue)
		os.RemoveAll("./1223")
	})
}

func TestShouldDeleteTape(t *testing.T) {
	Convey("should delete tape", t, func() {
		tape := new(Tape)
		tape.Path = "./tapes/"
		tape.SystemID = "1223"
		os.MkdirAll("./tapes/1223", os.ModePerm)
		So(Delete("1223"), ShouldBeNil)
		_, err := os.Open("./tapes/1223")
		So(err, ShouldNotBeNil)
		os.RemoveAll("./tapes")

	})
}

func TestShouldGetOrCreateTape(t *testing.T) {
	Convey("should get or create a new tape", t, func() {
		os.Mkdir("./tapes", os.ModePerm)
		defer os.RemoveAll("./tapes")
		tape, err := GetOrCreateTape("123", GetTapesPath())
		So(err, ShouldBeNil)
		So(tape.Exist(), ShouldBeTrue)
		fd, err := os.OpenFile("./tapes/123/tape.json", os.O_RDONLY, os.ModePerm)
		So(err, ShouldBeNil)
		if err == nil {
			fd.Close()
		}

	})
}

func TestShouldSaveSegments(t *testing.T) {
	os.Mkdir("./tapes", os.ModePerm)
	defer os.RemoveAll("./tapes")
	Convey("should save segments on tape", t, func() {

		tape, err := GetOrCreateTape("1123", "./tapes")
		So(err, ShouldBeNil)
		evt := domain.Event{
			Name: "event",
		}
		err = tape.RecordEvent(&evt)
		So(err, ShouldBeNil)
		So(len(tape.Segments), ShouldEqual, 1)
		fd, _ := os.Open("./tapes/1123")
		infos, _ := fd.Readdir(10)
		So(len(infos), ShouldEqual, 2)

	})

	Convey("should save reader on tape", t, func() {
		tape, err := GetOrCreateTape("1234", "./tapes")
		So(err, ShouldBeNil)
		evt := domain.Event{
			Name: "event",
		}

		d, _ := json.Marshal(evt)
		r := bufio.NewReader(bytes.NewReader(d))
		err = tape.RecordReader("dump.txt", "dump", r)
		So(err, ShouldBeNil)
		So(len(tape.Segments), ShouldEqual, 1)
		fd, _ := os.Open("./tapes/1234")
		infos, _ := fd.Readdir(10)
		So(len(infos), ShouldEqual, 2)
		os.RemoveAll("./tapes/1234")
	})

	Convey("should get existing tape", t, func() {
		GetOrCreateTape("b1234", "./tapes")
		_, err := GetOrCreateTape("b1234", "./tapes")
		So(err, ShouldBeNil)
	})

	Convey("should close tape", t, func() {
		tape, err := GetOrCreateTape("a1234", "./tapes")
		So(err, ShouldBeNil)
		evt := domain.Event{
			Name: "event",
		}
		d, _ := json.Marshal(evt)
		r := bufio.NewReader(bytes.NewReader(d))
		err = tape.RecordReader("dump.txt", "dump", r)
		So(err, ShouldBeNil)
		So(tape.Close(), ShouldBeNil)
		fd1, _ := os.Open("./tapes")
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
		So(tape.Exist(), ShouldBeFalse)
	})
}

func TestSegments(t *testing.T) {
	SystemID := "ec498841-59e5-47fd-8075-136d79155705"
	TapeID := "ec498841-59e5-47fd-8075-136d79155705_1544212971.rec"
	Convey("should check if segment is an event", t, func() {
		seg := Segment{
			SegmentType: "event",
		}
		So(seg.IsEvent(), ShouldBeTrue)
		seg.SegmentType = "other"
		So(seg.IsEvent(), ShouldBeFalse)
	})

	Convey("should get tape", t, func() {
		os.Setenv("TAPES_PATH", "../test_files")
		So(Restore(SystemID, TapeID), ShouldBeNil)
		_, err := GetTape(SystemID, "../test_files")
		So(err, ShouldBeNil)
		os.RemoveAll("../test_files/" + SystemID)
	})
}
