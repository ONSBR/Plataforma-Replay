package recorder

import (
	"os"
	"testing"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldRec(t *testing.T) {
	Convey("should rec an event for a system", t, func() {
		rec := newDefaultRecorder("ec498841-59e5-47fd-8075-136d79155705")
		evt := domain.Event{Name: "test", SystemID: "ec498841-59e5-47fd-8075-136d79155705"}
		err := rec.Rec(&evt)
		So(err, ShouldBeNil)
		So(rec.IsRecording(), ShouldBeTrue)
		os.RemoveAll("./ec498841-59e5-47fd-8075-136d79155705")
	})
}
