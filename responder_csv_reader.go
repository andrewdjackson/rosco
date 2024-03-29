package rosco

import (
	"fmt"
	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type ScenarioCSVReader struct {
	filepath string
	file     *os.File
	info     ResponderFileInfo
}

func NewScenarioCSVReader(filepath string) *ScenarioCSVReader {
	r := &ScenarioCSVReader{}
	r.filepath = filepath

	return r
}

func (r *ScenarioCSVReader) Load() (ResponderFileInfo, error) {
	var err error
	var data []*RawData
	var date time.Time

	if err = r.openFile(); err == nil {
		if err = gocsv.Unmarshal(r.file, &data); err != nil {
			err = fmt.Errorf("error parsing csv file %s (%s)", r.filepath, err)
			log.Errorf("%s", err)
		} else {
			log.Infof("successfully parsed %s, %d records read", r.filepath, len(data))

			if file, err := os.Stat(r.file.Name()); err == nil {
				date = file.ModTime()
			} else {
				date = time.Now()
				// don't report this as an error
				err = nil
			}

			r.info = ResponderFileInfo{
				Data: data,
				Description: ScenarioDescription{
					Name:     filepath.Base(r.file.Name()),
					Count:    len(data),
					Position: 0,
					Date:     date,
					Details:  ScenarioDetails{},
					Summary:  "MemsFCR Log File",
					FileType: "CSV",
				},
			}

			if len(data) > 0 {
				r.info.Description.Duration, err = getScenarioDuration(data[0].Time, data[r.info.Description.Count-1].Time)
			}
		}
	}

	return r.info, err
}

func (r *ScenarioCSVReader) openFile() error {
	var err error

	if r.file, err = os.OpenFile(r.filepath, os.O_RDWR|os.O_CREATE, os.ModePerm); err != nil {
		log.Errorf("error opening csv file %s (%s)", r.filepath, err)
	}

	return err
}
