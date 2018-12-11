package player

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/PMoneda/whaler"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Replay/recorder"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPlayer(t *testing.T) {
	testPath := "../test_files"

	systemID := "ec498841-59e5-47fd-8075-136d79155705"
	Convey("should not start playing while platform is recording", t, func() {
		So(os.Mkdir("./tapes", os.ModePerm), ShouldBeNil)
		rec := recorder.GetRecorder(systemID)
		evt := domain.Event{Name: "test", SystemID: systemID}
		So(rec.Rec(&evt), ShouldBeNil)
		p := GetPlayer(systemID)
		So(p.Play("123").Error(), ShouldEqual, "Cannot start playing events while platform is recording")
		os.RemoveAll("./tapes")
	})

	Convey("should stop container by name", t, func() {
		p := GetPlayer(systemID)
		So(p.StopDomainContaners("replay"), ShouldBeNil)
		So(p.StartDomainContaners("replay"), ShouldBeNil)
	})

	Convey("play a test tape", t, func() {
		os.Setenv("TAPES_PATH", testPath)
		p := GetPlayer(systemID)
		So(p.Play("ec498841-59e5-47fd-8075-136d79155705_1544212971.rec"), ShouldBeNil)
		os.RemoveAll(fmt.Sprintf("%s/%s", testPath, systemID))
	})

	Convey("should return error when event manager is down", t, func() {
		os.Setenv("TAPES_PATH", testPath)
		var timeout = 30 * time.Second
		whaler.StopContainer("event_manager", &timeout)
		p := GetPlayer(systemID)
		So(p.Play("ec498841-59e5-47fd-8075-136d79155705_1544212971.rec"), ShouldNotBeNil)
		os.RemoveAll(fmt.Sprintf("%s/%s", testPath, systemID))
		whaler.StartContainer("event_manager")
	})

}
