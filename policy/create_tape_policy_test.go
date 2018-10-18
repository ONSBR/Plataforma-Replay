package policy

import (
	"testing"

	"github.com/ONSBR/Plataforma-Replay/tape"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldVerifyCreateTapePolicy(t *testing.T) {
	Convey("should return true if a tape has a dump segment", t, func() {
		policy := CreateTapePolicy{}
		tp := tape.Tape{
			Segments: []*tape.Segment{
				&tape.Segment{
					SegmentType: "dump",
				},
			},
		}
		So(policy.DatabaseBackupAlreadyCreated(&tp), ShouldBeTrue)

		tpEvt := tape.Tape{
			Segments: []*tape.Segment{
				&tape.Segment{
					SegmentType: "event",
				},
			},
		}

		So(policy.DatabaseBackupAlreadyCreated(&tpEvt), ShouldBeFalse)
	})
}
