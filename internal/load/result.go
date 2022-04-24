package load

import (
	"benchutil/pkg/httploader"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"math"
)

const (
	outputYaml  = "yaml"
	outputJson  = "json"
	outputHuman = "human"
)

var outputFormats = map[string]struct{}{
	outputJson:  {},
	outputYaml:  {},
	outputHuman: {},
}

type report struct {
	Success     int `json:"success" yaml:"success"`
	Canceled    int `json:"canceled" yaml:"canceled"`
	Errors      int `json:"errors" yaml:"errors"`
	All         int `json:"all" yaml:"all"`
	AvgRespTime int `json:"avgRespTime" yaml:"avgRespTime"`
}

func (rep report) toBytes(format string) ([]byte, error) {
	if format == outputJson {
		return rep.toJson()
	}

	if format == outputYaml {
		return rep.toYaml()
	}

	if format == outputHuman {
		return rep.toHuman()
	}

	return nil, fmt.Errorf("unknown format: %s", format)
}

func (rep report) toHuman() ([]byte, error) {
	messageFormat := "Всего запросов: %d \nИз них \nУспешно: %d \nС ошибкой: %d \nОтменённых: %d \nСреднее время запроса(сек): %d"

	return []byte(fmt.Sprintf(messageFormat, rep.All, rep.Success, rep.Errors, rep.Canceled, rep.AvgRespTime)), nil
}

func (rep report) toJson() ([]byte, error) {
	return json.MarshalIndent(&rep, "", " ")
}

func (rep report) toYaml() ([]byte, error) {
	return yaml.Marshal(&rep)
}

func readResult(report httploader.Report, format string) ([]byte, error) {
	rep := loaderReportToInternal(report)
	return rep.toBytes(format)
}

func loaderReportToInternal(loaderRep httploader.Report) report {
	rep := report{
		Success:  loaderRep.Success,
		Errors:   loaderRep.Errors,
		Canceled: loaderRep.Cancelled,
		All:      loaderRep.All,
	}

	avgT := int(math.Round(loaderRep.AvgResponseTime.Seconds()))

	rep.AvgRespTime = avgT

	return rep
}
