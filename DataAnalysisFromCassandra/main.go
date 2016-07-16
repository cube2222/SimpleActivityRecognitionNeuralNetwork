package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/mjibson/go-dsp/fft"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
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

	stamp, _ := strconv.Atoi(os.Args[2])
	iter := session.Query(`SELECT time, pitch, roll, yaw FROM traininggyro WHERE userid='JonatanB' AND starttime=?`, int64(stamp)).Iter()
	gyroMutex.Lock()
	for iter.Scan(&time, &pitch, &roll, &yaw) {
		gyroData = append(gyroData, &orientation{time, pitch, roll, yaw})
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
		return
	}
	/*
		// Begin json
		data, err := json.Marshal(func() *struct {
			Data []orientation
		} {
			dataStruct := struct {
				Data []orientation
			}{}
			dataStruct.Data = make([]orientation, 0, len(gyroData))
			for _, item := range gyroData {
				dataStruct.Data = append(dataStruct.Data, *item)
			}
			return &dataStruct
		}())

		file, _ := os.Create("/tmp/myData.json")
		io.Copy(file, bytes.NewBuffer(data))
		file.Close()
		// End json
	*/

	//gyroData = smooth(gyroData, 1, 10.0) // 2 - 400.0
	gyroMutex.Unlock()

	myData := make([]float64, 0, len(gyroData))
	for _, item := range gyroData {
		myData = append(myData, math.Sqrt(item.Pitch*item.Pitch+item.Yaw*item.Yaw))
	}

	fourierData := fft.FFTReal(myData)
	/*
		for i := 0; i < len(fourierData); i++ {
			r, θ := cmplx.Polar(fourierData[i])
			θ *= 360.0 / (2 * math.Pi)
			if dsputils.Float64Equal(r, 0) {
				θ = 0 // (When the magnitude is close to 0, the angle is meaningless)
			}
			fmt.Printf("X(%d) = %.1f ∠ %.1f°\n", i, r, θ)
		}
	*/

	chartData := make([]float64, 0, len(fourierData))
	for _, item := range fourierData {
		chartData = append(chartData, real(item))
	}
	chartData = chartData[:len(chartData)/2]
	chartData = chartData[1:51]
	for i, item := range chartData {
		chartData[i] = math.Abs(item)
	}
	/*for i, _ := range fourierData {
		if i%2 == 0 {

			chartData = append(chartData[:i], chartData[i+1:]...)
		}
	}*/

	plotToFile(os.Args[3], chartData)
}

func smooth(data []*orientation, iterations int, averageSize float64) []*orientation {
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

	for i := 0; i < iterations; i++ {
		averageSize := averageSize

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

func plotToFile(name string, data []float64) error {
	file, err := os.Create("/tmp/" + name + ".jpg")
	defer file.Close()
	if err != nil {
		return err
	}
	dataXYs := make(plotter.XYs, len(data))

	for i := 0; i < len(data); i++ {
		dataXYs[i].X = float64(i)
		dataXYs[i].Y = data[i]
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

/*
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
		dataXYs[i].Y = (data[i].Pitch * data[i].Pitch + data[i].Yaw * data[i].Yaw)
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

	wt, err := p.WriterTo(vg.Inch * 16, vg.Inch * 16, "jpg")
	if err != nil {
		return err
	}

	_, err = wt.WriteTo(file)

	return err
}
*/

//func chartSinusoids(name string, )
