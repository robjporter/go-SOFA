
package routes

import (
	"io"
	"os"
	"fmt"
  //"bytes"
	"time"
	"bufio"
	"strings"
	"strconv"
	"io/ioutil"
  //"encoding/gob"
	"encoding/csv"
  //"encoding/base64"
	"github.com/kataras/iris"
	"github.com/tidwall/gjson"
)

var input = []string{}
var interesting = []string{}
var config = ""
var output = ""
var region = []string{}
var specialist = []string{}
var titlebar = []string{}

const (
	technologyColumn = 6
)

func removeExtension(name string) string {
	splits := strings.Split(name,".")
	return splits[0]
}

// GetIndexHandler handles the GET: /
func PostProcessReportHandler(ctx iris.Context) {
	var maps = []string{}
	step := ctx.PostValue("step")
	report := ctx.PostValue("report")
	mapp := ctx.PostValue("map")
	outputFileName := ctx.PostValue("output")

	if(step == "2") {
		files, err := ioutil.ReadDir("./uploads/maps")
		if err != nil {
				fmt.Println(err)
		}

		for _, f := range files {
			if(!f.IsDir()) {
					maps = append(maps,removeExtension(f.Name()))
			}
		}
	}

	if(step == "4") {
		mapFileName := "uploads/maps/" + mapp + ".json"
		reportFileName := "uploads/reports/" + report + ".csv"
		outputFileName = "uploads/processed/" + processFilename() + ".csv"


		loadMapFile(mapFileName)
		processMapFile()
	  loadDataFile(reportFileName)
  	saveCSVData(outputFileName)
	}

	if(step == "5") {
		fmt.Println(outputFileName)
		ctx.SendFile(outputFileName, "download.csv")
	}

	ctx.ViewData("Step", step)
	ctx.ViewData("Maps",maps)
	ctx.ViewData("Mapped",mapp)
	ctx.ViewData("Output",outputFileName)
	ctx.ViewData("Report", report)
  if err := ctx.View("processreport.html"); err != nil {
    ctx.Application().Logger().Infof(err.Error())
  }
}

func processMapFile() {
	processConfigInteresting()
	processConfigAlignment()
}

func processConfigInteresting() {
		for i := 0; i < int(gjson.Get(config,"INTERESTING.#").Int()); i++ {
			interesting = append(interesting,gjson.Get(config,"INTERESTING."+ strconv.Itoa(i) + ".name").String())
		}
}

func processConfigAlignment() {
		for i := 0; i < int(gjson.Get(config,"ALIGNMENT.#").Int()); i++ {
			region = append(region,gjson.Get(config,"ALIGNMENT."+ strconv.Itoa(i) + ".name").String())
			specialist = append(specialist,gjson.Get(config,"ALIGNMENT."+ strconv.Itoa(i) + ".specialist").String())
		}
}

func processFilename() string {
	t := time.Now()
	return "output-" + t.Format("20060102150405")
}

func loadMapFile(filename string) error {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	config = string(byteValue)

	return nil
}

func loadDataFile(filename string) int {
	f, _ := os.Open(filename)
	r := csv.NewReader(bufio.NewReader(f))
	count := 0
	for {
			record, err := r.Read()
			if err == io.EOF {
					break
			}
			if count == 0 {
				titlebar = record
			} else {
				value := processRecord(record)
				input = append(input,value)
			}
			count += 1
	}
	return count
}

func processRecord(record []string) string {
	for i := 0; i < len(interesting); i++ {
    if strings.Contains(strings.ToLower(record[technologyColumn]),strings.ToLower(interesting[i])) {
      proccessRecordDetail(record)
    }
  }
  return ""
}

func getJsonSize(json string, section string) int {
  return int(gjson.Get(json,section+".#").Int())
}

func splitJson(json string, section string) string {
  value := gjson.Get(json, section)
  return value.String()
}

func quoteCommas(text string) string {
  if strings.Contains(text,",") {
    return `"` + text + `"`
  }
  return text
}

func proccessRecordDetail(record []string) {
  splits := strings.Split(strings.ToLower(record[technologyColumn]),";")
  if len(splits) > 0 {
    for i := 0; i < len(splits); i++ {
      for j := 0; j < len(interesting); j++ {
        if strings.Contains(strings.ToLower(splits[i]),strings.ToLower(interesting[j])) {
          processInterestingRecord(record, strings.ToLower(splits[i]))
        }
      }
    }
  }
}

func processInterestingRecord(record []string, element string) {
  res := ""
  splits := strings.Split(element,"(")
  percent := ""
  if len(splits) > 0 {
    percent = strings.TrimRight(splits[1],")")
  }
  res = strings.TrimSpace(record[0]) + ","
  res += strings.TrimSpace(quoteCommas(record[1])) + ","
  res += strings.TrimSpace(quoteCommas(record[2])) + ","
  res += strings.TrimSpace(record[3]) + ","
  res += strings.TrimSpace(record[4]) + ","
  res += strings.TrimSpace(record[5]) + ","
  res += strings.TrimSpace(splits[0]) + ","
  res += strings.TrimSpace(percent) + ","
  res += strings.TrimSpace(calculateRevenue(record[5],percent)) + ","
  res += strings.TrimSpace(record[7]) + ","
  res += strings.TrimSpace(record[8]) + ","
  res += strings.TrimSpace(record[9]) + ","
  res += strings.TrimSpace(record[10]) + ","
  res += strings.TrimSpace(alignment(record[10])) + ","
  res += strings.TrimSpace(record[11]) + ","
  res += strings.TrimSpace(record[12]) + "\n"

  output += res
}

func calculateRevenue(total string, percent string) string {
  if total != "" && percent != "" {
    e,_ := strconv.ParseFloat(percent,64)
    s,_ := strconv.Atoi(total)

    if e > 0 && s > 0 {
      a := float64(s) / 100
      b := float64(a) * e
      c := round(b)
      return strconv.Itoa(c)
    }
  }
  return ""
}

func alignment(terr string) string {
  ret := ""
  if terr != "" {
    terr = strings.ToLower(strings.TrimSpace(terr))
    for i := 0; i < len(region); i++ {
      if strings.ToLower(region[i]) == terr {
        return specialist[i]
      }
    }
  }
  return ret
}

func round(val float64) int {
    if val < 0 { return int(val-0.5) }
    return int(val+0.5)
}

func saveCSVData(filename string) {
  addTitleBar()
  d1 := []byte(output)
  err := ioutil.WriteFile(filename, d1, 0644)
  if err != nil {
    fmt.Println("SAVING Failed")
  }
}

func addTitleBar() {
  var result = ""

  for i := 0; i < 7; i++ {
    result += titlebar[i] + ","
  }

  result += "Percent,"
  result += "DC Expected Product ($000s),"

  for i := 7; i < 11; i++{
    result += titlebar[i] + ","
  }

  result += "DC Sales Specialist,"

  for i := 11; i < len(titlebar); i++{
    result += titlebar[i] + ","
  }

  result = strings.TrimRight(result,",")
  output = result + "\n" + output
}
