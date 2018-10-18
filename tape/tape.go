package tape

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ONSBR/Plataforma-EventManager/domain"
)

//Tape is datalog struct to keep all events, dump file and metadata
type Tape struct {
	SystemID string     `json:"system_id"`
	Path     string     `json:"path"`
	State    string     `json:"state"`
	Segments []*Segment `json:"segments,omitempty"`
}

type Segment struct {
	Timestamp   int64         `json:"timestamp"`
	FileName    string        `json:"fileName"`
	SegmentType string        `json:"segment_type"`
	Content     *bufio.Reader `json:"-"`
}

//Record write segment into disk and update tape.json file
func (t *Tape) Record(seg *Segment) error {
	fd, err := os.Create(fmt.Sprintf("%s/%s/%s", t.Path, t.SystemID, seg.FileName))
	if err != nil {
		return err
	}
	seg.Content.WriteTo(fd)
	fd.Close()

	t.Segments = append(t.Segments, seg)
	return t.persist()
}

/*
Close O método close faz o fechamento da fita para gravação quando uma fita é fechada, não será mais possível
gravar mais nada nela, os arquivos existentes serão comprimidos no formato tar.gz
Esse arquivo .rec será disponibilizado para outras plataformas para reproduzir alguns cenários específicos
*/
func (t *Tape) Close() error {
	origin := fmt.Sprintf("%s/%s", t.Path, t.SystemID)
	buf := bytes.NewBuffer(nil)
	if err := compress(origin, buf); err != nil {
		return err
	}
	if err := os.RemoveAll(origin); err != nil {
		return err
	}
	ts := time.Now().Unix()
	return ioutil.WriteFile(fmt.Sprintf("%s/%s_%d.rec", t.Path, t.SystemID, ts), buf.Bytes(), 0600)
}

//RecordEvent create a segment based on event and persist
func (t *Tape) RecordEvent(event *domain.Event) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	ts := time.Now().Unix()
	seg := Segment{
		Content:     bufio.NewReader(bytes.NewReader(eventJSON)),
		FileName:    fmt.Sprintf("event_%d", ts),
		SegmentType: "event",
		Timestamp:   ts,
	}
	return t.Record(&seg)
}

//RecordReader create a segment based on reader
func (t *Tape) RecordReader(fileName, segmentType string, reader *bufio.Reader) error {
	ts := time.Now().Unix()
	seg := Segment{
		Content:     reader,
		FileName:    fileName,
		SegmentType: segmentType,
		Timestamp:   ts,
	}
	return t.Record(&seg)
}

//Dest returns a destination path to tape files
func (t *Tape) Dest() string {
	return t.Path
}

func (t *Tape) persist() error {
	jsonData, err := json.Marshal(t)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s/tape.json", t.Path, t.SystemID), jsonData, 0600); err != nil {
		return err
	}
	return nil
}

//IsRecording returns if a tape is open to write segments
func (t *Tape) IsRecording() bool {
	return t.State == "recording"
}

func (t *Tape) exist() bool {
	fd, err := os.Open(fmt.Sprintf("%s/%s", t.Path, t.SystemID))
	if err != nil {
		return false
	}
	fd.Close()
	return true
}

func GetTape(systemID string) *Tape {
	return nil
}

//NewTape creates or load a tape for a systemId
func NewTape(systemID, path string) (*Tape, error) {
	tape := new(Tape)
	tape.Path = path
	tape.SystemID = systemID
	tape.State = "recording"
	if !tape.exist() {
		err := os.Mkdir(fmt.Sprintf("%s/%s", tape.Path, systemID), os.ModePerm)
		if err != nil {
			return nil, err
		}
		if err := tape.persist(); err != nil {
			return nil, err
		}
	} else {
		tapeJSON, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/tape.json", tape.Path, systemID))
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(tapeJSON, tape); err != nil {
			return nil, err
		}
	}
	return tape, nil
}
