package postgres

import (
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldDumpDatabase(t *testing.T) {
	Convey("Should dump database with docker", t, func() {
		manager := new(DBPostgres)
		errBkp := manager.Backup("ec498841-59e5-47fd-8075-136d79155705", "./")
		So(errBkp, ShouldBeNil)
		fd, err := os.OpenFile("./dump.sql", os.O_RDONLY, os.ModePerm)
		So(err, ShouldBeNil)
		sample := make([]byte, 256)
		_, err = fd.Read(sample)
		So(strings.Contains(string(sample), "PostgreSQL database dump"), ShouldBeTrue)
		fd.Close()
		os.Remove("./dump.sql")
	})
}
