package main

import (
	"database/sql"
	"fmt"
	"github.com/grd/statistics"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type SXWDebug struct {
	Connected           bool
	db                  *sql.DB
	StartYear           int
	NYears              int
	Iterations          int
	RGroups             int
	TranspirationLayers int
	SoilLayers          int
	Years               []int
	PlotSize            float32
	BVT                 float32
}

type Info struct {
	Years               int
	StartYear           int
	EndYear             int
	Iterations          int
	TranspirationLayers int
	SoilLayers          int
	PlotSize            float32
	BVT                 float32
}

type RGroupSummaryRow struct {
	RGroup            string
	VegType           string
	BiomassMean       float64
	BiomassSD         float64
	RealSizeMean      float64
	RealSizeSD        float64
	PRMean            float64
	PRSD              float64
	YearsZeroPR       int
	TranspirationMean float64
	TranspirationSD   float64
}

func (sxw *SXWDebug) connect(file string) {
	if sxw.db != nil {
		sxw.Disconnect()
	}
	sxw.Connected = true
	var err error
	sxw.db, err = sql.Open("sqlite3", file) //"sxwdebug.sqlite3"
	if err != nil {
		log.Fatal(err)
	}

	rows, err := sxw.db.Query("SELECT StartYear,Years,Iterations,RGroups,TranspirationLayers,SoilLayers,PlotSize,BVT FROM info;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&sxw.StartYear, &sxw.NYears, &sxw.Iterations, &sxw.RGroups, &sxw.TranspirationLayers, &sxw.SoilLayers, &sxw.PlotSize, &sxw.BVT)
	}

	sxw.Years = make([]int, sxw.NYears, sxw.NYears)
	for i := 0; i < sxw.NYears; i++ {
		sxw.Years[i] = sxw.StartYear + i
	}
}

func (sxw *SXWDebug) Disconnect() {
	sxw.db.Close()
}

func (sxw *SXWDebug) getOutputFracSummary(iteration int) []float64 {
	var rows *sql.Rows
	var err error
	var sqlString string
	var size int

	summary := make([]float64, 15)

	sqlString = "SELECT FracGrass, FracShrub, FracTree, FracForb, FracBareGround FROM sxwinputvars"
	if iteration == 0 {
		size = sxw.NYears * sxw.Iterations
		sqlString += ";"
		rows, err = sxw.db.Query(sqlString)
	} else {
		size = sxw.NYears
		sqlString += " WHERE Iteration=%d;"
		rows, err = sxw.db.Query(fmt.Sprintf(sqlString, iteration))
	}

	Grass := make(statistics.Float64, size)
	Shrub := make(statistics.Float64, size)
	Tree := make(statistics.Float64, size)
	Forb := make(statistics.Float64, size)
	BareGround := make(statistics.Float64, size)

	if err != nil {
		log.Fatal(err)
	}

	i := 0
	for rows.Next() {
		rows.Scan(&Grass[i], &Shrub[i], &Tree[i], &Forb[i], &BareGround[i])
		i++
	}
	rows.Close()

	if iteration != 0 {
		rows, err = sxw.db.Query(fmt.Sprintf("SELECT COUNT(FracGrass) FROM sxwinputvars WHERE FracGrass=0 AND Iteration=%d;", iteration))
		if err != nil {
			log.Fatal(err)
		}
		rows.Next()
		rows.Scan(&summary[2])
		rows.Close()

		rows, err = sxw.db.Query(fmt.Sprintf("SELECT COUNT(FracShrub) FROM sxwinputvars WHERE FracShrub=0 AND Iteration=%d;", iteration))
		if err != nil {
			log.Fatal(err)
		}
		rows.Next()
		rows.Scan(&summary[5])
		rows.Close()

		rows, err = sxw.db.Query(fmt.Sprintf("SELECT COUNT(FracForb) FROM sxwinputvars WHERE FracForb=0 AND Iteration=%d;", iteration))
		if err != nil {
			log.Fatal(err)
		}
		rows.Next()
		rows.Scan(&summary[8])
		rows.Close()

		rows, err = sxw.db.Query(fmt.Sprintf("SELECT COUNT(FracTree) FROM sxwinputvars WHERE FracTree=0 AND Iteration=%d;", iteration))
		if err != nil {
			log.Fatal(err)
		}
		rows.Next()
		rows.Scan(&summary[11])
		rows.Close()

		rows, err = sxw.db.Query(fmt.Sprintf("SELECT COUNT(FracBareGround) FROM sxwinputvars WHERE FracBareGround=0 AND Iteration=%d;", iteration))
		if err != nil {
			log.Fatal(err)
		}
		rows.Next()
		rows.Scan(&summary[14])
		rows.Close()
	} else {
		rows, err = sxw.db.Query("SELECT COUNT(FracGrass) FROM sxwinputvars WHERE FracGrass=0;")
		if err != nil {
			log.Fatal(err)
		}
		rows.Next()
		rows.Scan(&summary[2])
		rows.Close()

		rows, err = sxw.db.Query("SELECT COUNT(FracShrub) FROM sxwinputvars WHERE FracShrub=0;")
		if err != nil {
			log.Fatal(err)
		}
		rows.Next()
		rows.Scan(&summary[5])
		rows.Close()

		rows, err = sxw.db.Query("SELECT COUNT(FracForb) FROM sxwinputvars WHERE FracForb=0;")
		if err != nil {
			log.Fatal(err)
		}
		rows.Next()
		rows.Scan(&summary[8])
		rows.Close()

		rows, err = sxw.db.Query("SELECT COUNT(FracTree) FROM sxwinputvars WHERE FracTree=0;")
		if err != nil {
			log.Fatal(err)
		}
		rows.Next()
		rows.Scan(&summary[11])
		rows.Close()

		rows, err = sxw.db.Query("SELECT COUNT(FracBareGround) FROM sxwinputvars WHERE FracBareGround=0;")
		if err != nil {
			log.Fatal(err)
		}
		rows.Next()
		rows.Scan(&summary[14])
		rows.Close()
	}

	summary[0] = statistics.Mean(&Grass)
	summary[1] = statistics.Sd(&Grass)

	summary[3] = statistics.Mean(&Shrub)
	summary[4] = statistics.Sd(&Shrub)

	summary[6] = statistics.Mean(&Forb)
	summary[7] = statistics.Sd(&Forb)

	summary[9] = statistics.Mean(&Tree)
	summary[10] = statistics.Sd(&Tree)

	summary[12] = statistics.Mean(&BareGround)
	summary[13] = statistics.Sd(&BareGround)

	return summary
}

func (sxw *SXWDebug) getOutputVarsSummary(iteration int) []float64 {
	var rows *sql.Rows
	var err error
	var sqlString string
	var size int

	summary := make([]float64, 14)

	sqlString = "SELECT MAP_mm, MAT_C, AET_cm, AT_cm, TotalRelsize, TotalPR, TotalTransp FROM sxwoutputvars"
	if iteration == 0 {
		size = sxw.NYears * sxw.Iterations
		sqlString += ";"
		rows, err = sxw.db.Query(sqlString)
	} else {
		size = sxw.NYears
		sqlString += " WHERE Iteration=%d;"
		rows, err = sxw.db.Query(fmt.Sprintf(sqlString, iteration))
	}

	MAP := make(statistics.Float64, size)
	MAT := make(statistics.Float64, size)
	AET := make(statistics.Float64, size)
	AT := make(statistics.Float64, size)
	TRS := make(statistics.Float64, size)
	TPR := make(statistics.Float64, size)
	TTR := make(statistics.Float64, size)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		rows.Scan(&MAP[i], &MAT[i], &AET[i], &AT[i], &TRS[i], &TPR[i], &TTR[i])
		MAP[i] /= 10
		i++
	}

	summary[0] = statistics.Mean(&MAP)
	summary[1] = statistics.Sd(&MAP)
	summary[2] = statistics.Mean(&MAT)
	summary[3] = statistics.Sd(&MAT)
	summary[4] = statistics.Mean(&AET)
	summary[5] = statistics.Sd(&AET)
	summary[6] = statistics.Mean(&AT)
	summary[7] = statistics.Sd(&AT)
	summary[8] = statistics.Mean(&TRS)
	summary[9] = statistics.Sd(&TRS)
	summary[10] = statistics.Mean(&TPR)
	summary[11] = statistics.Sd(&TPR)
	summary[12] = statistics.Mean(&TTR)
	summary[13] = statistics.Sd(&TTR)

	return summary
}

func (sxw *SXWDebug) getRGroupSummary(iteration int) []RGroupSummaryRow {
	var rows *sql.Rows
	var err error

	summaryRows := make([]RGroupSummaryRow, sxw.RGroups)

	rgroups := sxw.getRGroups(-1)
	for i := range rgroups {
		summaryRows[i].RGroup = rgroups[i]
	}

	vegTypes := sxw.getRGroupsVegProdType()
	for i := range vegTypes {
		switch vegTypes[i] {
		case 1:
			summaryRows[i].VegType = "Tree"
		case 2:
			summaryRows[i].VegType = "Shrub"
		case 3:
			summaryRows[i].VegType = "Grass"
		case 4:
			summaryRows[i].VegType = "Forb"
		}
	}

	for i := range rgroups {
		if iteration == 0 {
			rows, err = sxw.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM sxwoutputrgroup WHERE RGroupID=%d AND PR>0;", i+1))
			if err != nil {
				log.Fatal(err)
			}
			YearsZeroPR := 0
			for rows.Next() {
				rows.Scan(&YearsZeroPR)
			}
			rows.Close()
			summaryRows[i].YearsZeroPR = sxw.NYears*sxw.Iterations - YearsZeroPR

			rows, err = sxw.db.Query(fmt.Sprintf("SELECT Biomass FROM sxwoutputrgroup WHERE RGroupID=%d;", i+1))
			if err != nil {
				log.Fatal(err)
			}
			Biomass := make(statistics.Float64, sxw.NYears*sxw.Iterations)
			j := 0
			for rows.Next() {
				rows.Scan(&Biomass[j])
				j++
			}
			rows.Close()
			summaryRows[i].BiomassMean = statistics.Mean(&Biomass)
			summaryRows[i].BiomassSD = statistics.Sd(&Biomass)

			rows, err = sxw.db.Query(fmt.Sprintf("SELECT Realsize FROM sxwoutputrgroup WHERE RGroupID=%d;", i+1))
			if err != nil {
				log.Fatal(err)
			}
			RealSize := make(statistics.Float64, sxw.NYears*sxw.Iterations)
			j = 0
			for rows.Next() {
				rows.Scan(&RealSize[j])
				j++
			}
			rows.Close()
			summaryRows[i].RealSizeMean = statistics.Mean(&RealSize)
			summaryRows[i].RealSizeSD = statistics.Sd(&RealSize)

			//0 values for PR will influence the average the wrong way
			rows, err = sxw.db.Query(fmt.Sprintf("SELECT PR FROM sxwoutputrgroup WHERE RGroupID=%d AND PR > 0;", i+1))
			if err != nil {
				log.Fatal(err)
			}
			PR := make(statistics.Float64, sxw.NYears*sxw.Iterations)
			j = 0
			for rows.Next() {
				rows.Scan(&PR[j])
				j++
			}
			rows.Close()
			summaryRows[i].PRMean = statistics.Mean(&PR)
			summaryRows[i].PRSD = statistics.Sd(&PR)

			rows, err = sxw.db.Query(fmt.Sprintf("SELECT Transpiration FROM sxwoutputrgroup WHERE RGroupID=%d;", i+1))
			if err != nil {
				log.Fatal(err)
			}
			Transpiration := make(statistics.Float64, sxw.NYears*sxw.Iterations)
			j = 0
			for rows.Next() {
				rows.Scan(&Transpiration[j])
				j++
			}
			rows.Close()
			summaryRows[i].TranspirationMean = statistics.Mean(&Transpiration)
			summaryRows[i].TranspirationSD = statistics.Sd(&Transpiration)

		} else {
			rows, err = sxw.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM sxwoutputrgroup WHERE RGroupID=%d AND PR>0 AND Iteration=%d;", i+1, iteration))
			if err != nil {
				log.Fatal(err)
			}
			YearsZeroPR := 0
			for rows.Next() {
				rows.Scan(&YearsZeroPR)
			}
			rows.Close()
			summaryRows[i].YearsZeroPR = sxw.NYears - YearsZeroPR

			rows, err = sxw.db.Query(fmt.Sprintf("SELECT Biomass FROM sxwoutputrgroup WHERE RGroupID=%d AND Iteration=%d;", i+1, iteration))
			if err != nil {
				log.Fatal(err)
			}
			Biomass := make(statistics.Float64, sxw.NYears)
			j := 0
			for rows.Next() {
				rows.Scan(&Biomass[j])
				j++
			}
			rows.Close()
			summaryRows[i].BiomassMean = statistics.Mean(&Biomass)
			summaryRows[i].BiomassSD = statistics.Sd(&Biomass)

			rows, err = sxw.db.Query(fmt.Sprintf("SELECT Realsize FROM sxwoutputrgroup WHERE RGroupID=%d AND Iteration=%d;", i+1, iteration))
			if err != nil {
				log.Fatal(err)
			}
			RealSize := make(statistics.Float64, sxw.NYears)
			j = 0
			for rows.Next() {
				rows.Scan(&RealSize[j])
				j++
			}
			rows.Close()
			summaryRows[i].RealSizeMean = statistics.Mean(&RealSize)
			summaryRows[i].RealSizeSD = statistics.Sd(&RealSize)

			//0 values for PR will influence the average the wrong way
			rows, err = sxw.db.Query(fmt.Sprintf("SELECT PR FROM sxwoutputrgroup WHERE RGroupID=%d AND PR > 0 AND Iteration=%d;", i+1, iteration))
			if err != nil {
				log.Fatal(err)
			}
			PR := make(statistics.Float64, sxw.NYears)
			j = 0
			for rows.Next() {
				rows.Scan(&PR[j])
				j++
			}
			rows.Close()
			summaryRows[i].PRMean = statistics.Mean(&PR)
			summaryRows[i].PRSD = statistics.Sd(&PR)

			rows, err = sxw.db.Query(fmt.Sprintf("SELECT Transpiration FROM sxwoutputrgroup WHERE RGroupID=%d AND Iteration=%d;", i+1, iteration))
			if err != nil {
				log.Fatal(err)
			}
			Transpiration := make(statistics.Float64, sxw.NYears)
			j = 0
			for rows.Next() {
				rows.Scan(&Transpiration[j])
				j++
			}
			rows.Close()
			summaryRows[i].TranspirationMean = statistics.Mean(&Transpiration)
			summaryRows[i].TranspirationSD = statistics.Sd(&Transpiration)
		}
	}
	return summaryRows
}

func (sxw *SXWDebug) getMAPvalue(year int, iteration int) int {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT MAP_mm FROM sxwoutputvars WHERE Year=%d AND Iteration=%d;", year, iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var MAP_mm int

	for rows.Next() {
		rows.Scan(&MAP_mm)
	}
	return MAP_mm
}

func (sxw *SXWDebug) getMAP(iteration int) []int {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT MAP_mm FROM sxwoutputvars WHERE Iteration=%d;", iteration))
	if err != nil {
		fmt.Println("ERROR\n\n")
		log.Fatal(err)
	}
	defer rows.Close()

	MAP_mm := make([]int, sxw.NYears)

	i := 0
	for rows.Next() {
		rows.Scan(&(MAP_mm[i]))
		i++
	}
	return MAP_mm
}

func (sxw *SXWDebug) getMATvalue(year int, iteration int) int {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT MAT_C FROM sxwoutputvars WHERE Year=%d AND Iteration=%d;", year, iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var MAT_C int

	for rows.Next() {
		rows.Scan(&MAT_C)
	}
	return MAT_C
}

func (sxw *SXWDebug) getMAT(iteration int) []int {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT MAT_C FROM sxwoutputvars WHERE Iteration=%d;", iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	MAT_C := make([]int, sxw.NYears)

	for rows.Next() {
		rows.Scan(&MAT_C)
	}
	return MAT_C
}

func (sxw *SXWDebug) getAETvalue(year int, iteration int) int {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT AET_cm FROM sxwoutputvars WHERE Year=%d AND Iteration=%d;", year, iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var AET int

	for rows.Next() {
		rows.Scan(&AET)
	}
	return AET
}

func (sxw *SXWDebug) getAET(iteration int) []int {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT AET_cm FROM sxwoutputvars WHERE Iteration=%d;", iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	AET := make([]int, sxw.NYears)

	for rows.Next() {
		rows.Scan(&AET)
	}
	return AET
}

func (sxw *SXWDebug) getATvalue(year int, iteration int) int {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT AT_cm FROM sxwoutputvars WHERE Year=%d AND Iteration=%d;", year, iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var AT int

	for rows.Next() {
		rows.Scan(&AT)
	}
	return AT
}

func (sxw *SXWDebug) getAT(iteration int) []int {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT AT_cm FROM sxwoutputvars WHERE Iteration=%d;", iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	AT := make([]int, sxw.NYears)

	for rows.Next() {
		rows.Scan(&AT)
	}
	return AT
}

func (sxw *SXWDebug) getTotalRelsize(iteration int) []float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT TotalRelsize FROM sxwoutputvars WHERE Iteration=%d;", iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	TotalRelsize := make([]float32, sxw.NYears)

	for rows.Next() {
		rows.Scan(&TotalRelsize)
	}
	return TotalRelsize
}

func (sxw *SXWDebug) getTotalPR(iteration int) []float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT TotalPR FROM sxwoutputvars WHERE Iteration=%d;", iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	TotalPR := make([]float32, sxw.NYears)

	for rows.Next() {
		rows.Scan(&TotalPR)
	}
	return TotalPR
}

func (sxw *SXWDebug) getTotalTransp(iteration int) []float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT TotalTransp FROM sxwoutputvars WHERE Iteration=%d;", iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	TotalTransp := make([]float32, sxw.NYears)

	for rows.Next() {
		rows.Scan(&TotalTransp)
	}
	return TotalTransp
}

func (sxw *SXWDebug) getOutputVars(iteration int) [][]float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT Year, MAP_mm, MAT_C, AET_cm, AT_cm FROM sxwoutputvars WHERE Iteration=%d ORDER BY Year;", iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	vars := make([][]float32, sxw.NYears)

	i := 0
	for rows.Next() {
		vars[i] = make([]float32, 5)
		rows.Scan(&vars[i][0], &vars[i][1], &vars[i][2], &vars[i][3], &vars[i][4])
		vars[i][1] /= 10
		i++
	}
	return vars
}

func (sxw *SXWDebug) getTransp(year int, iteration int, vegType int) [][]float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT January, February, March, April, May, June, July, August, September, October, November, December FROM sxwoutputtranspiration WHERE YEAR=%d AND Iteration=%d AND VegProdType=%d;", year, iteration, vegType))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	Transp := make([][]float32, sxw.SoilLayers)

	i := 0
	for rows.Next() {
		Transp[i] = make([]float32, 12)
		rows.Scan(&Transp[i][0], &Transp[i][1], &Transp[i][2], &Transp[i][3], &Transp[i][4], &Transp[i][5], &Transp[i][6], &Transp[i][7], &Transp[i][8], &Transp[i][9], &Transp[i][10], &Transp[i][11])
		i++
	}
	return Transp
}

func (sxw *SXWDebug) getTranspSum(iteration int) [][]float32 {
	Transp := make([][]float32, sxw.NYears)
	for i := 0; i < sxw.NYears; i++ {
		Transp[i] = make([]float32, 6)
		Transp[i][0] = float32(sxw.StartYear + i)
	}

	for i := 0; i < 5; i++ {
		rows, err := sxw.db.Query(fmt.Sprintf("SELECT SUM(January+February+March+April+May+June+July+August+September+October+November+December) AS Total FROM sxwoutputtranspiration WHERE Iteration=%d AND VegProdType=%d GROUP BY YEAR;", iteration, i))
		if err != nil {
			log.Fatal(err)
		}
		j := 0
		for rows.Next() {
			rows.Scan(&Transp[j][i+1])
			j++
		}
	}
	return Transp
}

func (sxw *SXWDebug) getSWCBulk(year int, iteration int) [][]float32 {
	//SELECT SUM(January+February+March+April+May+June+July+August+September+October+November+December) AS Total FROM sxwoutputswcbulk WHERE Iteration=1 GROUP BY YEAR;
	//SELECT January,February,March,April,May,June,July,August,September,October,November,December FROM sxwoutputswcbulk WHERE Iteration=1 GROUP BY YEAR;
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT January, February, March, April, May, June, July, August, September, October, November, December FROM sxwoutputswcbulk WHERE YEAR=%d AND Iteration=%d;", year, iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	SWCBulk := make([][]float32, sxw.SoilLayers)

	i := 0
	for rows.Next() {
		SWCBulk[i] = make([]float32, 12)
		rows.Scan(&SWCBulk[i][0], &SWCBulk[i][1], &SWCBulk[i][2], &SWCBulk[i][3], &SWCBulk[i][4], &SWCBulk[i][5], &SWCBulk[i][6], &SWCBulk[i][7], &SWCBulk[i][8], &SWCBulk[i][9], &SWCBulk[i][10], &SWCBulk[i][11])
		i++
	}
	return SWCBulk
}

func (sxw *SXWDebug) getBiomass(iteration int, RGroupID int) []float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT Biomass FROM sxwoutputrgroup WHERE Iteration=%d AND RGroupID=%d;", iteration, RGroupID))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	Biomass := make([]float32, sxw.NYears)

	for rows.Next() {
		rows.Scan(&Biomass)
	}
	return Biomass
}

func (sxw *SXWDebug) getRealsize(iteration int, RGroupID int) []float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT Realsize FROM sxwoutputrgroup WHERE Iteration=%d AND RGroupID=%d;", iteration, RGroupID))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	Realsize := make([]float32, sxw.NYears)

	for rows.Next() {
		rows.Scan(&Realsize)
	}
	return Realsize
}

func (sxw *SXWDebug) getPR(iteration int, RGroupID int) []float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT PR FROM sxwoutputrgroup WHERE Iteration=%d AND RGroupID=%d;", iteration, RGroupID))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	PR := make([]float32, sxw.NYears)

	for rows.Next() {
		rows.Scan(&PR)
	}
	return PR
}

func (sxw *SXWDebug) getRGroupTransp(iteration int, RGroupID int) []float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT Transpiration FROM sxwoutputrgroup WHERE Iteration=%d AND RGroupID=%d;", iteration, RGroupID))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	Transpiration := make([]float32, sxw.NYears)

	for rows.Next() {
		rows.Scan(&Transpiration)
	}
	return Transpiration
}

func (sxw *SXWDebug) getOutputRGroup(iteration int, RGroupID int) [][]float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT Year,Biomass,Realsize,PR,Transpiration FROM sxwoutputrgroup WHERE Iteration=%d AND RGroupID=%d;", iteration, RGroupID))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	rgroup := make([][]float32, sxw.NYears)
	i := 0
	for rows.Next() {
		rgroup[i] = make([]float32, 5)
		rows.Scan(&rgroup[i][0], &rgroup[i][1], &rgroup[i][2], &rgroup[i][3], &rgroup[i][4])
		i++
	}
	return rgroup
}

func (sxw *SXWDebug) getOutputProd(year int, iteration int) [][]float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT Month,BMass,PctLive,LAIlive,VegCov,TotAGB FROM sxwoutputprod WHERE YEAR=%d AND Iteration=%d ORDER BY Month;", year, iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	vals := make([][]float32, 12)

	i := 0
	for rows.Next() {
		vals[i] = make([]float32, 6)
		rows.Scan(&vals[i][0], &vals[i][1], &vals[i][2], &vals[i][3], &vals[i][4], &vals[i][5])
		i++
	}
	return vals
}

func (sxw *SXWDebug) getInputFracGrass(iteration int) []float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT FracGrass FROM sxwinputvars WHERE Iteration=%d;", iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fracGrass := make([]float32, sxw.NYears)

	for rows.Next() {
		rows.Scan(&fracGrass)
	}
	return fracGrass
}

func (sxw *SXWDebug) getInputFracShrub(iteration int) []float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT FracShrub FROM sxwinputvars WHERE Iteration=%d;", iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fracShrub := make([]float32, sxw.NYears)

	for rows.Next() {
		rows.Scan(&fracShrub)
	}
	return fracShrub
}

func (sxw *SXWDebug) getInputFracTree(iteration int) []float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT FracTree FROM sxwinputvars WHERE Iteration=%d;", iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fracTree := make([]float32, sxw.NYears)

	for rows.Next() {
		rows.Scan(&fracTree)
	}
	return fracTree
}

func (sxw *SXWDebug) getInputFracForb(iteration int) []float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT FracForb FROM sxwinputvars WHERE Iteration=%d;", iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	FracForb := make([]float32, sxw.NYears)

	for rows.Next() {
		rows.Scan(&FracForb)
	}
	return FracForb
}

func (sxw *SXWDebug) getInputFracBareGround(iteration int) []float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT FracBareGround FROM sxwinputvars WHERE Iteration=%d;", iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	FracBareGround := make([]float32, sxw.NYears)

	for rows.Next() {
		rows.Scan(&FracBareGround)
	}
	return FracBareGround
}

func (sxw *SXWDebug) getInputFracsYear(year int, iteration int) []float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT FracTree,FracShrub,FracGrass,FracForb,FracBareGround FROM sxwinputvars WHERE Year=%d AND Iteration=%d;", year, iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	Fracs := make([]float32, 5)

	for rows.Next() {
		rows.Scan(&Fracs[0], &Fracs[1], &Fracs[2], &Fracs[3], &Fracs[4])
	}
	return Fracs
}

func (sxw *SXWDebug) getInputFracs(iteration int) [][]float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT Year,FracTree,FracShrub,FracGrass,FracForb,FracBareGround FROM sxwinputvars WHERE Iteration=%d ORDER BY Year;", iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	Fracs := make([][]float32, sxw.NYears)
	i := 0
	for rows.Next() {
		Fracs[i] = make([]float32, 6)
		rows.Scan(&Fracs[i][0], &Fracs[i][1], &Fracs[i][2], &Fracs[i][3], &Fracs[i][4], &Fracs[i][5])
		i++
	}
	return Fracs
}

func (sxw *SXWDebug) getInputSoils(year int, iteration int) [][]float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT Layer,Tree_trco,Shrub_trco,Grass_trco,Forb_trco FROM sxwinputsoils WHERE Year=%d AND Iteration=%d ORDER BY Layer;", year, iteration))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	trco := make([][]float32, sxw.SoilLayers)
	i := 0
	for rows.Next() {
		trco[i] = make([]float32, 5)
		rows.Scan(&trco[i][0], &trco[i][1], &trco[i][2], &trco[i][3], &trco[i][4])
		i++
	}
	return trco
}

func (sxw *SXWDebug) getInputProd(year int, iteration int, vegType int) [][]float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT Month,Litter,Biomass,PLive,LAI_conv FROM sxwinputprod WHERE Year=%d AND Iteration=%d AND VegProdType=%d ORDER BY Month;", year, iteration, vegType))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	prod := make([][]float32, 12, 12)
	i := 0
	for rows.Next() {
		prod[i] = make([]float32, 5, 5)
		rows.Scan(&prod[i][0], &prod[i][1], &prod[i][2], &prod[i][3], &prod[i][4])
		i++
	}
	return prod
}

func (sxw *SXWDebug) getRGroups(VegType int) []string {
	var rows *sql.Rows
	var err error

	if VegType > 0 {
		rows, err = sxw.db.Query(fmt.Sprintf("SELECT NAME FROM rgroups WHERE VegProdType=%d ORDER BY ID;", VegType))
	} else {
		rows, err = sxw.db.Query("SELECT NAME FROM rgroups ORDER BY ID;")
	}

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	names := make([]string, 0, sxw.RGroups)
	i := 0
	for rows.Next() {
		names = append(names, "")
		rows.Scan(&names[i])
		i++
	}
	return names
}

func (sxw *SXWDebug) getRGroupsVegProdType() []int {
	rows, err := sxw.db.Query("SELECT VegProdType FROM rgroups ORDER BY ID;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	VegProdType := make([]int, sxw.RGroups)
	i := 0
	for rows.Next() {
		rows.Scan(&VegProdType[i])
		i++
	}
	return VegProdType
}

func (sxw *SXWDebug) getRootsRelative(year int, iteration int, RGroupID int) [][]float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT Layer,January, February, March, April, May, June, July, August, September, October, November, December FROM sxwrootsrelative WHERE Year=%d AND Iteration=%d AND RGroupID=%d ORDER BY Layer;", year, iteration, RGroupID))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	rootsRel := make([][]float32, sxw.TranspirationLayers)
	i := 0
	for rows.Next() {
		rootsRel[i] = make([]float32, 13)
		rows.Scan(&rootsRel[i][0], &rootsRel[i][1], &rootsRel[i][2], &rootsRel[i][3], &rootsRel[i][4], &rootsRel[i][5], &rootsRel[i][6], &rootsRel[i][7], &rootsRel[i][8], &rootsRel[i][9], &rootsRel[i][10], &rootsRel[i][11], &rootsRel[i][12])
		i++
	}
	return rootsRel
}

func (sxw *SXWDebug) getRootsSum(year int, iteration int, vegType int) [][]float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT Layer,January, February, March, April, May, June, July, August, September, October, November, December FROM sxwrootssum WHERE Year=%d AND Iteration=%d AND VegProdType=%d ORDER BY Layer;", year, iteration, vegType))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	rootsSum := make([][]float32, sxw.TranspirationLayers)
	i := 0
	for rows.Next() {
		rootsSum[i] = make([]float32, 13)
		rows.Scan(&rootsSum[i][0], &rootsSum[i][1], &rootsSum[i][2], &rootsSum[i][3], &rootsSum[i][4], &rootsSum[i][5], &rootsSum[i][6], &rootsSum[i][7], &rootsSum[i][8], &rootsSum[i][9], &rootsSum[i][10], &rootsSum[i][11], &rootsSum[i][12])
		i++
	}
	return rootsSum
}

func (sxw *SXWDebug) getRootsXphen(RGroupID int) [][]float32 {
	rows, err := sxw.db.Query(fmt.Sprintf("SELECT Layer, January, February, March, April, May, June, July, August, September, October, November, December FROM sxwrootsxphen WHERE RGroupID=%d;", RGroupID))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	XPhen := make([][]float32, sxw.TranspirationLayers)

	i := 0
	for rows.Next() {
		XPhen[i] = make([]float32, 13)
		rows.Scan(&XPhen[i][0], &XPhen[i][1], &XPhen[i][2], &XPhen[i][3], &XPhen[i][4], &XPhen[i][5], &XPhen[i][6], &XPhen[i][7], &XPhen[i][8], &XPhen[i][9], &XPhen[i][10], &XPhen[i][11], &XPhen[i][12])
		i++
	}
	return XPhen
}
