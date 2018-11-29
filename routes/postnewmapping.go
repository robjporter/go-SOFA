
package routes

import (
	"io"
	"os"
	"fmt"
  "bytes"
	"bufio"
	"strings"
	"strconv"
	"io/ioutil"
  "encoding/gob"
	"encoding/csv"
  "encoding/base64"
	"github.com/kataras/iris"
)

// GetIndexHandler handles the GET: /
func PostNewMappingHandler(ctx iris.Context) {
	technologies := []string{}
	regions := []string{}
	selected := []string{}

	securetechnologies := ctx.PostValue("securetechnologies")
	secureregions := ctx.PostValue("secureregions")

	step := ctx.PostValue("step")
	report := ctx.PostValue("report")
	regions = ctx.FormValues()["regions"]
	selected = ctx.FormValues()["selected"]
	technologies = ctx.FormValues()["technologies"]
	mapName := ctx.PostValue("mapname")

	if(securetechnologies == "" && technologies != nil) {
		securetechnologies = ToGOB64(technologies)
	}

	if(secureregions == "" && selected != nil) {
		secureregions = ToGOB64(selected)
	}

	if(step == "2") {
		processStep2(report)
		technologyFile := getSaveFilename("technology",report)
		data,err := loadDataFromFile(technologyFile)
		if(err == nil ) {
			technologies = data
		} else {
			fmt.Println(err)
		}
	}

	if(step == "3") {
		geographyFile := getSaveFilename("geography",report)
		data,err := loadDataFromFile(geographyFile)
		if(err == nil ) {
			regions = data
		} else {
			fmt.Println(err)
		}
	}

	if(step == "4") {
		geographyFile := getSaveFilename("geography",report)
		data,err := loadDataFromFile(geographyFile)

		if( err == nil ) {
			for i := 0; i < len(data); i++ {
				selected = append(selected, ctx.PostValue(data[i]))
			}
			secureregions = ToGOB64(selected)
		}

	}

	if(step == "5") {
		err := saveMapFile(report, mapName, secureregions, securetechnologies)
		if err != nil {
			fmt.Println(err)
		} else {
  		ctx.Redirect("/", iris.StatusSeeOther)
		}
	}

	ctx.ViewData("Title", "SOFA - Index Page")
	ctx.ViewData("Step", step)
	ctx.ViewData("Report", report)
	ctx.ViewData("Regions", regions)
	ctx.ViewData("Selected", selected)
	ctx.ViewData("SecureTechnologies", securetechnologies)
	ctx.ViewData("SecureRegions", secureregions)
	ctx.ViewData("Technologies", technologies)

	if err := ctx.View("newmapping.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

func ToGOB64(m []string) string {
    b := bytes.Buffer{}
    e := gob.NewEncoder(&b)
    err := e.Encode(m)
    if err != nil { fmt.Println(`failed gob Encode`, err) }
    return base64.StdEncoding.EncodeToString(b.Bytes())
}

func FromGOB64(str string) []string {
    m := []string{}
    by, err := base64.StdEncoding.DecodeString(str)
    if err != nil { fmt.Println(`failed base64 Decode`, err); }
    b := bytes.Buffer{}
    b.Write(by)
    d := gob.NewDecoder(&b)
    err = d.Decode(&m)
    if err != nil { fmt.Println(`failed gob Decode`, err); }
    return m
}

func stripArray(data string) []string {
	return FromGOB64(data)
}

func saveMapFile(reportName string, mapName string, selected string, technologies string) error {
	output := `{"INTERESTING": [`

	// CONVERT ENCODED TECHNOLOGIES TO STRING ARRAY
	arr := stripArray(technologies)

	// INTERESTING TECHNOLOGIES
	for i := 0; i < len(arr); i++ {
		output += `{"name":"`+arr[i]+`"},`
	}
	output = strings.TrimRight(output,",")
	output += `],`

	// CONVERT ENCODED SELECTED TO STRING ARRAY
	arr = stripArray(selected)

	// LOAD GEORGRAPHY DATA
	geographyFile := getSaveFilename("geography",reportName)
	data,err := loadDataFromFile(geographyFile)

	// REGION TO SALES SPECIALIST MAPPING
	output += `"ALIGNMENT": [`
	if( err == nil ) {
		for i := 0; i < len(data); i++ {
			output += `{"name":"`+data[i]+`","specialist":"`+arr[i]+`"},`
		}
	}
	output = strings.TrimRight(output,",")
	output += `]`

	// FINISH JSON STRING
	output += `}`

	return saveConfigFile(mapName,output)
}

func processStep2(reportName string) {
	count := loadCSVData("uploads/reports/"+reportName)
	fmt.Println("Completed successfully processing: " + strconv.Itoa(count) + " records.")
}

func loadDataFromFile(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

func loadCSVData(file string) int {
    f, _ := os.Open(file)
    r := csv.NewReader(bufio.NewReader(f))
		technologyValues := []string{}
		geographyValues := []string{}
    count := 0
    for {
        record, err := r.Read()
        if err == io.EOF {
            break
        }
        if count > 0 {
          technologyValues = processTechnologyRecord(record, technologyValues)
	        geographyValues = processGeographyRecord(record, geographyValues)
        }
        count += 1
    }
		saveValues(technologyValues, "technology", file)
		saveValues(geographyValues, "geography", file)
		return count
}

func getSaveFilename(typer string, fileName string) string {
	splits := strings.Split(fileName,".")
	splitss := strings.TrimPrefix(splits[0],"uploads/reports/")
	newFileName := "uploads/mappings/" + splitss + "-" + typer + "." + splits[1]
	return newFileName
}

func saveConfigFile(mapName string, data string) error {
  d1 := []byte(data)
  err := ioutil.WriteFile("uploads/maps/" + mapName + ".json", d1, 0644)
  return err
}

func saveValues(data []string, typer string, fileName string) {
	newFileName := getSaveFilename(typer,fileName)
	file, err := os.Create(newFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range data {
		fmt.Fprintln(w, line)
	}
	w.Flush()
}

func processTechnologyRecord(record []string, values []string) []string {
	splits := strings.Split(record[6],";")
	if(len(splits) > 0) {
		splitss := strings.Split(splits[0],"(")
		if(len(splitss) > 0) {
			values = AppendIfMissing(values,strings.TrimSpace(splitss[0]))
		}
	}
	return values
}

func processGeographyRecord(record []string, values []string) []string {
	values = AppendIfMissing(values,strings.TrimSpace(record[10]))
	return values
}

func AppendIfMissing(slice []string, i string) []string {
	if(i != "") {
    for _, ele := range slice {
        if ele == i {
            return slice
        }
    }
    return append(slice, i)
	}
	return slice
}
