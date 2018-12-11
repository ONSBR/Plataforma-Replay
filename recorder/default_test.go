package recorder

import (
	"os"
	"testing"

	"github.com/ONSBR/Plataforma-Replay/tape"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldRec(t *testing.T) {
	Convey("should rec an event for a system", t, func() {
		os.MkdirAll("./tapes", os.ModePerm)
		rec := newDefaultRecorder("ec498841-59e5-47fd-8075-136d79155705")
		evt := domain.Event{Name: "test", SystemID: "ec498841-59e5-47fd-8075-136d79155705"}
		err := rec.Rec(&evt)
		So(err, ShouldBeNil)
		So(rec.IsRecording(), ShouldBeTrue)
		os.RemoveAll("./tapes")
	})
}

func TestShouldEjectTape(t *testing.T) {
	systemID := "ec498841-59e5-47fd-8075-136d79155705"
	tapeID := "ec498841-59e5-47fd-8075-136d79155705_1544212971.rec"
	Convey("should rec an event for a system", t, func() {
		os.Setenv("TAPES_PATH", "../test_files")
		tape.Restore(systemID, tapeID)
		rec := newDefaultRecorder(systemID)
		_, err := rec.Eject()
		So(err, ShouldBeNil)
		tapes, _ := rec.AvailableTapesToDownload(systemID)
		So(len(tapes), ShouldBeGreaterThan, 1)
		for _, s := range tapes {
			if s != tapeID {
				So(tape.Delete(s), ShouldBeNil)
				break
			}
		}
	})
}

func TestShouldGetEventFromSegment(t *testing.T) {
	systemID := "ec498841-59e5-47fd-8075-136d79155705"
	tapeID := "ec498841-59e5-47fd-8075-136d79155705_1544212971.rec"
	Convey("should get an event from segment", t, func() {
		os.Setenv("TAPES_PATH", "../test_files")
		tape.Restore(systemID, tapeID)
		rec := newDefaultRecorder(systemID)
		currTape, _ := rec.GetOrCreateTape(systemID)
		for _, seg := range currTape.Segments {
			if seg.IsEvent() {
				evt, err := seg.Event(currTape)
				So(err, ShouldBeNil)
				So(evt, ShouldNotBeNil)
				break
			}

		}
		os.RemoveAll("../test_files/" + systemID)
	})
}
