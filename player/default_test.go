package player

import (
	"os"
	"testing"

	"github.com/ONSBR/Plataforma-Replay/recorder"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPlayer(t *testing.T) {
	Convey("should not start playing while platform is recording", t, func() {
		So(os.Mkdir("./tapes", os.ModePerm), ShouldBeNil)
		systemID := "ec498841-59e5-47fd-8075-136d79155705"
		rec := recorder.GetRecorder(systemID)
		evt := domain.Event{Name: "test", SystemID: systemID}
		So(rec.Rec(&evt), ShouldBeNil)
		p := GetPlayer(systemID)
		So(p.Play("123").Error(), ShouldEqual, "Cannot start playing events while platform is recording")
		os.RemoveAll("./tapes")
	})
}
