package tape

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/labstack/gommon/log"
)

//Tape is datalog struct to keep all events, dump file and metadata
type Tape struct {
	SystemID     string     `json:"system_id"`
	Path         string     `json:"path"`
	State        string     `json:"state"`
	CompressFile string     `json:"compress_file"`
	Segments     []*Segment `json:"segments,omitempty"`
}

type Segment struct {
	Timestamp   int64         `json:"timestamp"`
	FileName    string        `json:"fileName"`
	SegmentType string        `json:"segment_type"`
	Content     *bufio.Reader `json:"-"`
}

func (seg *Segment) IsEvent() bool {
	return seg.SegmentType == "event"
}

func (seg *Segment) Event(t *Tape) (*domain.Event, error) {
	if !seg.IsEvent() {
		return nil, fmt.Errorf("Segment is not an event")
	}
	if buf, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", t.Dest(), seg.FileName)); err != nil {
		return nil, err
	} else {
		evt := new(domain.Event)
		if err := json.Unmarshal(buf, evt); err != nil {
			return nil, err
		}
		return evt, nil
	}
}

//Record write segment into disk and update tape.json file
func (t *Tape) Record(seg *Segment) error {
	if seg.Content != nil {
		fd, err := os.Create(fmt.Sprintf("%s/%s", t.Dest(), seg.FileName))
		if err != nil {
			return err
		}
		if _, err := seg.Content.WriteTo(fd); err != nil {
			return err
		}
		fd.Close()
	}
	t.Segments = append(t.Segments, seg)
	return t.persist()
}

/*
Close O método close faz o fechamento da fita para gravação quando uma fita é fechada, não será mais possível
gravar mais nada nela, os arquivos existentes serão comprimidos no formato tar.gz
Esse arquivo .rec será disponibilizado para outras plataformas para reproduzir alguns cenários específicos
*/
func (t *Tape) Close() error {
	t.State = "closed"
	ts := time.Now().Unix()
	t.CompressFile = fmt.Sprintf("%s/%s_%d.rec", t.Path, t.SystemID, ts)
	err := t.persist()
	if err != nil {
		return err
	}
	origin := t.Dest()
	buf := bytes.NewBuffer(nil)
	if err := compress(origin, buf); err != nil {
		return err
	}
	if err := os.RemoveAll(origin); err != nil {
		return err
	}
	return ioutil.WriteFile(t.CompressFile, buf.Bytes(), 0600)
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
	return fmt.Sprintf("%s/%s", t.Path, t.SystemID)
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

func (t *Tape) Exist() bool {
	fd, err := os.Open(fmt.Sprintf("%s/%s", t.Path, t.SystemID))
	if err != nil {
		return false
	}
	fd.Close()
	return true
}

//GetOrCreateTape creates or load a tape for a systemId
func GetOrCreateTape(systemID, path string) (*Tape, error) {
	tape := new(Tape)
	tape.Path = path
	tape.SystemID = systemID
	tape.State = "recording"
	if !tape.Exist() {
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

func GetTapesPath() string {
	return env.Get("TAPES_PATH", "./tapes")
}

func GetTape(systemID, path string) (*Tape, error) {
	tape := new(Tape)
	tape.Path = path
	tape.SystemID = systemID
	if !tape.Exist() {
		return nil, fmt.Errorf("tape for system %s not exist", systemID)
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

//Delete tape from disk
func Delete(tapeID string) error {
	if tapeID == "" {
		return fmt.Errorf("empty tapeID")
	}
	path := GetTapesPath()
	str := fmt.Sprintf("%s/%s", path, tapeID)
	log.Info("removing file at:", str)
	return os.Remove(str)
}

//Restore tape file to a folder on  tape's path
func Restore(systemID string, tapeID string) error {
	return decompress(fmt.Sprintf("%s/%s", GetTapesPath(), tapeID), fmt.Sprintf("%s/%s", GetTapesPath(), systemID))
}
