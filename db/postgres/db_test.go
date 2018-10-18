package postgres

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldDumpDatabase(t *testing.T) {
	Convey("Should dump database with docker", t, func() {
		manager := new(DBPostgres)
		reader, errBkp := manager.Backup("ec498841-59e5-47fd-8075-136d79155705")
		So(errBkp, ShouldBeNil)

		sample := make([]byte, 256)
		reader.Read(sample)
		So(strings.Contains(string(sample), "PostgreSQL database dump"), ShouldBeTrue)
	})
}
