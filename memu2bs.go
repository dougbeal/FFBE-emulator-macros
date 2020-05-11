package main

import (
	"bufio"
	"encoding/json"
	//	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const memu6 = "MULTI2"
const yd = 1280.0
const xd = 720.0
const Timestamp = "\"Timestamp\": %d,\n"
const Delta = "\"Delta\": 0,"
const EventMouseUp = "\"EventType\": \"MouseUp\","
const EventMouseDown = "\"EventType\": \"MouseDown\","
const Xc = "\"X\": %.2f,\n"
const Yc = "\"Y\": %.2f,\n"

type Event struct {
	Timestamp int64
	Delta int64
	EventType string
	X float64
	Y float64
}

type Macro struct {
	TimeCreated string
	Name string
	Events []Event
	LoopType string
	LoopNumber int
	LoopTime int
	LoopInterval int
	Acceleration int
	PlayOnStart bool
	DonotShowWindowOnFinish bool
	RestartPlayer bool
	RestartPlayerAfterMinutes int
	ShortCut string
	UserName string
	MacroID string `json:"MacroId"`
}

func main() {
	memufile := os.Args[1]

	file, err := os.Open(memufile)
	check(err)

	defer file.Close()

	base := path.Base(memufile)
	name := strings.TrimSuffix(base, path.Ext(memufile))
	destFilename := name + ".json"

	macro := Macro{}
	macro.Events = make([]Event, 0)

	macro.TimeCreated = time.Now().Format("20060102T150405")
	macro.Name = name
	macro.LoopType = "TillLoopNumber"
	macro.LoopNumber = 1
	macro.LoopTime = 0
	macro.LoopInterval = 0
	macro.Acceleration = 1.0
	macro.PlayOnStart = false
	macro.DonotShowWindowOnFinish = false
	macro.RestartPlayer = false
	macro.RestartPlayerAfterMinutes = 60

	
	scanner := bufio.NewScanner(file)
	x := 0.0
	y := 0.0
	for scanner.Scan() {
		event := Event{}
		s := scanner.Text()
		sp := strings.Split(s, "--")
		if len(sp) == 3 {

			muls := sp[2]
			mulsp := strings.Split(muls, ":")
			if memu6 == mulsp[0] {
				if len(mulsp) == 7 {
					//fmt.Println("{")
					t, et := strconv.ParseInt(sp[0], 10, 64)
					check(et)
					// convert nanoseconds to milliseconds
					event.Timestamp = t/1000
					event.Delta = 0
					//fmt.Printf(Timestamp, t)
					//fmt.Println(Delta)
					if mulsp[6] == "0" {
						//fmt.Println(EventMouseDown)
						event.EventType = "MouseDown"
						// only mousedown events have valid x,y
						fx, ex := strconv.ParseFloat(mulsp[4], 64)
						fy, ey := strconv.ParseFloat(mulsp[5], 64)
						check(ex)
						x = fx
						check(ey)
						y = fy
					}
					if mulsp[6] == "2" {
						//fmt.Println(EventMouseUp)
						event.EventType = "MouseUp"
					}
					xp := (x / xd) * 100
					yp := (y / yd) * 100
					// output most recent x,y
					event.X = math.Round(xp*100)/100
					event.Y = math.Round(yp*100)/100
					
					//fmt.Printf(Xc, math.Round(xp*100)/100)
					//fmt.Printf(Yc, math.Round(yp*100)/100)
					//fmt.Println("},")
					macro.Events = append(macro.Events, event)
				}
			}

		}

	}

	
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(macro, " ", "  "); 
	
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(destFilename, b, 0777)
	check(err)
	//fmt.Println(string(b))	

	// fmt.Println(memufile)
}
