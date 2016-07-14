package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"image/color"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

type orientation struct {
	Timestamp int64   `json:"timestamp"`
	Roll      float64 `json:"roll"`
	Pitch     float64 `json:"pitch"`
	Yaw       float64 `json:"yaw"`
}

type orientationTrainingData struct {
	UserId    string      `json:"userID"`
	Activity  string      `json:"activity"`
	StartTime int64       `json:"starttime"`
	CurData   orientation `json:"orientation"`
}

var gyroData []*orientation
var gyroMutex sync.RWMutex

func main() {

	gyroData = make([]*orientation, 0, 10000)
	gyroMutex = sync.RWMutex{}

	var err error
	add := make([]string, 0, 5)

	//add, err = net.LookupHost("cassandra")
	add = append(add, os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	credentials := gocql.PasswordAuthenticator{Username: os.Getenv("CASSANDRA_USERNAME"), Password: os.Getenv("CASSANDRA_PASSWORD")}

	cluster := gocql.NewCluster(add[0])
	if len(credentials.Username) > 0 {
		cluster.Authenticator = credentials
	}
	cluster.Timeout = time.Second * 4
	cluster.ProtoVersion = 4
	cluster.Keyspace = "activitytracking"
	session, err := cluster.CreateSession()
	for err != nil {
		fmt.Println("Error when connecting for active use. Trying again in 2 seconds.")
		fmt.Println(err)
		err = nil
		session, err = cluster.CreateSession()
		time.Sleep(time.Second * 2)
	}

	var time int64
	var pitch float64
	var roll float64
	var yaw float64

	iter := session.Query(`SELECT time, pitch, roll, yaw FROM traininggyro WHERE userid='JonatanB' AND starttime=1468518081922`).Iter()
	gyroMutex.Lock()
	for iter.Scan(&time, &pitch, &roll, &yaw) {
		gyroData = append(gyroData, &orientation{time, pitch, roll, yaw})
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
		return
	}

	gyroData = smooth(gyroData)
	gyroMutex.Unlock()

	plotToFile("dataBasic", gyroData)
}

func smooth(data []*orientation) []*orientation {
	startTime := data[0].Timestamp

	newData := make([]*orientation, data[len(data)-1].Timestamp-startTime)
	for i := 0; i < len(data)-1; i++ {
		unitCount := data[i+1].Timestamp - data[i].Timestamp
		//newData[data[i].Timestamp] = &rawOrientation{Roll: data[i].Roll, Pitch: data[i].Pitch, Yaw: data[i].Yaw}
		for j := data[i].Timestamp - startTime; j <= data[i].Timestamp-startTime+(unitCount/2); j++ {
			newData[j] = &orientation{Roll: data[i].Roll, Pitch: data[i].Pitch, Yaw: data[i].Yaw, Timestamp: j + startTime}
		}
		for j := data[i].Timestamp - startTime + (unitCount / 2) + 1; j < data[i].Timestamp-startTime+unitCount; j++ {
			newData[j] = &orientation{Roll: data[i+1].Roll, Pitch: data[i+1].Pitch, Yaw: data[i+1].Yaw, Timestamp: j + startTime}
		}
	}

	for i := 0; i < 2; i++ {
		averageSize := 400.0

		averageData := make([]*orientation, 0, int(averageSize))
		for i := 0; i < int(averageSize); i++ {
			averageData = append(averageData, newData[0])
		}
		curAverage := newData[0]

		for i := 0; i < len(newData); i++ {
			curAverage = &orientation{
				Roll:      (curAverage.Roll*(averageSize-1) + newData[i].Roll) / averageSize,
				Pitch:     (curAverage.Pitch*(averageSize-1) + newData[i].Pitch) / averageSize,
				Yaw:       (curAverage.Yaw*(averageSize-1) + newData[i].Yaw) / averageSize,
				Timestamp: newData[i].Timestamp,
			}
			newData[i] = curAverage
		}
	}

	return newData
}

func plotToFile(name string, data []*orientation) error {
	file, err := os.Create("/tmp/" + name + ".jpg")
	defer file.Close()
	if err != nil {
		return err
	}
	dataXYs := make(plotter.XYs, len(data))

	minX := data[0].Timestamp

	for i := 0; i < len(data); i++ {
		dataXYs[i].X = float64(data[i].Timestamp - minX)
		dataXYs[i].Y = math.Sqrt(data[i].Pitch*data[i].Pitch + data[i].Yaw*data[i].Yaw)
		fmt.Println(dataXYs[i].X, dataXYs[i].Y)
	}

	line, err := plotter.NewLine(dataXYs)
	if err != nil {
		return err
	}

	line.LineStyle.Width = vg.Points(1)
	line.LineStyle.Color = color.RGBA{R: 255, B: 255, A: 255}

	p, err := plot.New()
	if err != nil {
		return err
	}

	p.Title.Text = name
	p.X.Label.Text = "t"
	p.Y.Label.Text = "Data"
	p.Add(plotter.NewGrid())
	p.Add(line)

	wt, err := p.WriterTo(vg.Inch*16, vg.Inch*16, "jpg")
	if err != nil {
		return err
	}

	_, err = wt.WriteTo(file)

	return err
}
