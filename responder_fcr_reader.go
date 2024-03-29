package rosco

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type ScenarioFCRReader struct {
	filepath string
	file     *os.File
	info     ResponderFileInfo
}

func NewScenarioFCRReader(filepath string) *ScenarioFCRReader {
	r := &ScenarioFCRReader{}
	r.filepath = filepath

	return r
}

func (r *ScenarioFCRReader) Load() (ResponderFileInfo, error) {
	var err error
	var fcrData ScenarioFile

	if err = r.openFile(); err == nil {
		if data, err := ioutil.ReadAll(r.file); err == nil {
			if err = json.Unmarshal(data, &fcrData); err != nil {
				log.Errorf("error parsing fcr file %s (%s)", r.filepath, err)
			} else {
				if fcrData.Count > 0 {
					r.info = ResponderFileInfo{
						Data: fcrData.RawData,
						Description: ScenarioDescription{
							Name:     fcrData.Name,
							Count:    fcrData.Count,
							Position: 0,
							Date:     fcrData.Date,
							Details:  ScenarioDetails{},
							Summary:  fcrData.Summary,
							FileType: "FCR",
						},
					}

					if len(data) > 0 {
						r.info.Description.Duration, err = getScenarioDuration(fcrData.RawData[0].Time, fcrData.RawData[r.info.Description.Count-1].Time)
					}

					log.Infof("successfully parsed %s, %d records read", r.filepath, len(r.info.Data))
				} else {
					err = fmt.Errorf("file contains no data")
					log.Errorf("%s", err)
				}
			}
		} else {
			err = fmt.Errorf("error reading file (%s)", err)
			log.Errorf("%s", err)
		}
	}

	return r.info, err
}

func (r *ScenarioFCRReader) openFile() error {
	var err error

	if r.file, err = os.OpenFile(r.filepath, os.O_RDWR|os.O_CREATE, os.ModePerm); err != nil {
		log.Errorf("error opening fcr file %s (%s)", r.filepath, err)
	}

	return err
}
