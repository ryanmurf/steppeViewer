// steppeViewer project main.go
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	//prj := r.FormValue("project")
	var t *template.Template
	switch title {
	case "":
		t, _ = template.ParseFiles("html/index.html")
		t.Execute(w, "none")
	case "inputProd.html":
		t, _ = template.ParseFiles("html/inputProd.html")
		t.Execute(w, "")
	case "inputSoils.html":
		t, _ = template.ParseFiles("html/inputSoils.html")
		t.Execute(w, "")
	case "inputVars.html":
		t, _ = template.ParseFiles("html/inputVars.html")
		t.Execute(w, "")
	case "outputProd.html":
		t, _ = template.ParseFiles("html/outputProd.html")
		t.Execute(w, "")
	case "outputRGroup.html":
		t, _ = template.ParseFiles("html/outputRGroup.html")
		t.Execute(w, "")
	case "outputVars.html":
		t, _ = template.ParseFiles("html/outputVars.html")
		t.Execute(w, "")
	case "overallSummary.html":
		t, _ = template.ParseFiles("html/overallSummary.html")
		t.Execute(w, "")
	case "outputTranspiration.html":
		t, _ = template.ParseFiles("html/outputTranspiration.html")
		t.Execute(w, "")
	case "outputTransGroup.html":
		t, _ = template.ParseFiles("html/outputTransGroup.html")
		t.Execute(w, "")
	case "sxwXPhen.html":
		t, _ = template.ParseFiles("html/sxwXPhen.html")
		t.Execute(w, "")
	case "sxwRootsRel.html":
		t, _ = template.ParseFiles("html/sxwRootsRel.html")
		t.Execute(w, "")
	case "sxwRootsSum.html":
		t, _ = template.ParseFiles("html/sxwRootsSum.html")
		t.Execute(w, "")
	}
}

func scriptHandle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[len("/script/"):])
}

func setDB(w http.ResponseWriter, r *http.Request) {
	prj := projects[r.FormValue("project")]
	db := prj.sxwDB

	db.Disconnect()
	prj.getRunInfo()
	db.connect(prj.SXWDebugFile)
	fmt.Fprintf(w, "%s", prj.SXWDebugFile)
}

func getSXWInput(w http.ResponseWriter, r *http.Request) {
	table := r.URL.Path[len("/sxw/input/"):]
	db := projects[r.FormValue("project")].sxwDB

	switch table {
	case "prod":
		//column := r.FormValue("name")
		year, _ := strconv.Atoi(r.FormValue("Year"))
		iteration, _ := strconv.Atoi(r.FormValue("Iteration"))
		vegType, _ := strconv.Atoi(r.FormValue("VegType"))

		data := db.getInputProd(year, iteration, vegType)
		values, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(values)
	case "soils":
		year, _ := strconv.Atoi(r.FormValue("Year"))
		iteration, _ := strconv.Atoi(r.FormValue("Iteration"))

		data := db.getInputSoils(year, iteration)
		values, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(values)
	case "vars":
		//year, _ := strconv.Atoi(r.FormValue("Year"))
		iteration, _ := strconv.Atoi(r.FormValue("Iteration"))

		data := db.getInputFracs(iteration)
		values, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(values)
	}
}

type point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func getSXWOutput(w http.ResponseWriter, r *http.Request) {
	table := r.URL.Path[len("/sxw/output/"):]
	db := projects[r.FormValue("project")].sxwDB
	switch table {
	case "vars":
		column := r.FormValue("name")
		iteration, _ := strconv.Atoi(r.FormValue("Iteration"))

		switch column {
		case "MAP_mm":
			points := make([]*point, db.NYears, db.NYears)
			MAP := db.getMAP(iteration)
			for i := 0; i < db.NYears; i++ {
				points[i] = &point{X: db.Years[i], Y: MAP[i]}
			}
			values, err := json.Marshal(points)
			if err != nil {
				fmt.Println(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(values)
		case "AET_cm":
		case "AT_cm":
		case "table":
			iteration, _ := strconv.Atoi(r.FormValue("Iteration"))
			data := db.getOutputVars(iteration)
			values, err := json.Marshal(data)
			if err != nil {
				fmt.Println(err)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(values)
		}
	case "prod":
		year, _ := strconv.Atoi(r.FormValue("Year"))
		iteration, _ := strconv.Atoi(r.FormValue("Iteration"))

		data := db.getOutputProd(year, iteration)
		values, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(values)
	case "rgroup":
		rgroupid, _ := strconv.Atoi(r.FormValue("RGroupID"))
		iteration, _ := strconv.Atoi(r.FormValue("Iteration"))

		data := db.getOutputRGroup(iteration, rgroupid)
		values, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(values)
	case "transpiration":
		iteration, _ := strconv.Atoi(r.FormValue("Iteration"))

		data := db.getTranspSum(iteration)
		values, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(values)
	}
}

func getInfo(w http.ResponseWriter, r *http.Request) {
	value := r.URL.Path[len("/info/"):]
	db := projects[r.FormValue("project")].sxwDB
	switch value {
	case "":
		data := Info{Years: db.NYears, StartYear: db.StartYear, EndYear: db.StartYear + db.NYears, Iterations: db.Iterations, TranspirationLayers: db.TranspirationLayers, SoilLayers: db.SoilLayers, PlotSize: db.PlotSize, BVT: db.BVT}
		js, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	case "iterations":
		js, err := json.Marshal(db.Iterations)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	case "years":
		js, err := json.Marshal(db.Years)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	case "rgroups":
		VegType, err := strconv.Atoi(r.FormValue("VegType"))
		if err != nil {
			VegType = -1
		}

		js, err := json.Marshal(db.getRGroups(VegType))
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func getSXWSummary(w http.ResponseWriter, r *http.Request) {
	value := r.URL.Path[len("/sxw/summary/"):]
	db := projects[r.FormValue("project")].sxwDB
	switch value {
	case "rgroup":
		iteration, _ := strconv.Atoi(r.FormValue("Iteration"))
		js, err := json.Marshal(db.getRGroupSummary(iteration))
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	case "vars":
		iteration, _ := strconv.Atoi(r.FormValue("Iteration"))
		js, err := json.Marshal(db.getOutputVarsSummary(iteration))
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	case "fracs":
		iteration, _ := strconv.Atoi(r.FormValue("Iteration"))
		js, err := json.Marshal(db.getOutputFracSummary(iteration))
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func getProjectsHandle(w http.ResponseWriter, r *http.Request) {
	value := r.URL.Path[len("/projects/"):]
	switch value {
	case "":
		keys := make([]string, 0, len(projects))
		for k := range projects {
			keys = append(keys, k)
		}
		js, err := json.Marshal(keys)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	case "description":
		prj := projects[r.FormValue("project")]
		js, err := json.Marshal(prj.Description)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

}

func getSXW(w http.ResponseWriter, r *http.Request) {
	value := r.URL.Path[len("/sxw/"):]
	db := projects[r.FormValue("project")].sxwDB
	switch value {
	case "xphen":
		rgroupid, _ := strconv.Atoi(r.FormValue("RGroupID"))
		js, err := json.Marshal(db.getRootsXphen(rgroupid))
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	case "rootsrel":
		year, _ := strconv.Atoi(r.FormValue("Year"))
		rgroupid, _ := strconv.Atoi(r.FormValue("RGroupID"))
		iteration, _ := strconv.Atoi(r.FormValue("Iteration"))

		data := db.getRootsRelative(year, iteration, rgroupid)
		values, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(values)
	case "rootssum":
		year, _ := strconv.Atoi(r.FormValue("Year"))
		vegType, _ := strconv.Atoi(r.FormValue("VegType"))
		iteration, _ := strconv.Atoi(r.FormValue("Iteration"))

		data := db.getRootsSum(year, iteration, vegType)
		values, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(values)
	}
}

var projects = make(map[string]Project)
var defPrj string

func main() {
	prjs := getProjects()
	defPrj = prjs[0].Name
	for i, item := range prjs {
		projects[item.Name] = prjs[i]
	}

	http.HandleFunc("/info/", getInfo)
	http.HandleFunc("/sxw/summary/", getSXWSummary)
	http.HandleFunc("/sxw/input/", getSXWInput)
	http.HandleFunc("/sxw/output/", getSXWOutput)
	http.HandleFunc("/sxw/", getSXW)
	http.HandleFunc("/projects/", getProjectsHandle)
	http.HandleFunc("/connect/", setDB)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/script/", scriptHandle)
	http.ListenAndServe(":8090", nil)
}
