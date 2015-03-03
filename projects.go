package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Project struct {
	Name           string
	Description    string
	StepWatVersion int
	Executed       bool
	ExecutedDate   time.Time
	SXWDebugFile   string
	sxwDB          SXWDebug
}

func getProjects() []Project {
	files, _ := ioutil.ReadDir("./Projects/")

	Projects := make([]Project, len(files))

	for i, f := range files {
		Projects[i].Name = f.Name()
		Projects[i].getDescription()
		Projects[i].getRunInfo()
		if Projects[i].Executed {
			Projects[i].sxwDB.connect(Projects[i].SXWDebugFile)
		}
	}

	return Projects
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (p *Project) getRunInfo() {
	file := "./Projects/" + p.Name + "/Output/sxwdebug.sqlite3"
	if info, err := os.Stat(file); err == nil {
		//Read file
		p.Executed = true
		p.ExecutedDate = info.ModTime()
		p.SXWDebugFile = file
	} else {
		//Create File
		p.Executed = false
		p.ExecutedDate = info.ModTime()
		p.SXWDebugFile = ""
	}
}

func (p *Project) getDescription() {
	file := "./Projects/" + p.Name + "/info.txt"
	if _, err := os.Stat(file); err == nil {
		//Read file
		dat, e := ioutil.ReadFile(file)
		check(e)
		p.Description = string(dat)
	} else {
		//Create File
		fmt.Println("creating file")
		os.Create(file)
		e := ioutil.WriteFile(file, []byte(""), 0644)
		check(e)
	}
}
