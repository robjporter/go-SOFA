
package routes

import (
	"io"
	"os"
	"fmt"
	"bufio"
	"strings"
	"encoding/csv"
	"github.com/kataras/iris"
)

// GetIndexHandler handles the GET: /
func PostNewMappingHandler(ctx iris.Context) {
	step := ctx.PostValue("step")
	report := ctx.PostValue("report")
	mapping := ctx.PostValue("mapping")
	technologies := []string{}
	//technologies = ctx.PostValue("technologies")

	if(step == "") { fmt.Println("STEP: 1")}
	if(step == "2") {
		mapping = processStep2(report)
		technologyFile := getSaveFilename("technology",report)
		data,err := loadDataFromFile(technologyFile)
		if(err == nil ) {
			technologies = data
		} else {
			fmt.Println(err)
		}
	}
	if(step == "3") { fmt.Println("STEP: 3")}
	if(step == "4") { fmt.Println("STEP: 4")}

	ctx.ViewData("Title", "Index Page")
	ctx.ViewData("Step", step)
	ctx.ViewData("Report", report)
	ctx.ViewData("Mapping", mapping)
	ctx.ViewData("Technolgoies",technologies)
	if err := ctx.View("newmapping.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

func processStep2(reportName string) string {
	count := loadCSVData("uploads/reports/"+reportName)
	fmt.Println(count)
	return ""
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
