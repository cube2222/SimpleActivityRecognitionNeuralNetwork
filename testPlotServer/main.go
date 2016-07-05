package main

import (
	"encoding/json"
	"fmt"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"image/color"
	"io/ioutil"
	"net/http"
	"sync"
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

	http.HandleFunc("/training/orientation", handleOrientationTraining)
	http.HandleFunc("/debuginfo", printDebugInfo)
	http.HandleFunc("/plot", plotGyroData)
	http.ListenAndServe(":3000", nil)

}

func handleOrientationTraining(w http.ResponseWriter, r *http.Request) {
	// Read and parse request data.

	myData := &orientationTrainingData{}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &myData)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	gyroMutex.Lock()
	gyroData = append(gyroData, &myData.CurData)
	gyroMutex.Unlock()

	// Android app expects the Status Created code for responses signaling success.
	w.WriteHeader(http.StatusCreated)
}

func printDebugInfo(w http.ResponseWriter, r *http.Request) {
	gyroMutex.RLock()
	for _, item := range gyroData {
		fmt.Fprintln(w, item.Timestamp, ":", item.Pitch, ",", item.Roll, ",", item.Yaw)
	}
	gyroMutex.RUnlock()
}

func plotGyroData(w http.ResponseWriter, r *http.Request) {
	gyroMutex.RLock()
	if len(gyroData) == 0 {
		gyroMutex.RUnlock()
		return
	}
	fmt.Println(len(gyroData))
	gyroMutex.RUnlock()

	gyroMutex.Lock()
	smoothedGyroData := smooth(gyroData)
	gyroMutex.Unlock()

	rollData := make(plotter.XYs, len(smoothedGyroData))
	pitchData := make(plotter.XYs, len(smoothedGyroData))
	yawData := make(plotter.XYs, len(smoothedGyroData))
	testingData := make(plotter.XYs, len(smoothedGyroData))

	minX := smoothedGyroData[0].Timestamp

	for i := 0; i < len(smoothedGyroData); i++ {
		rollData[i].X = float64(smoothedGyroData[i].Timestamp - minX)
		rollData[i].Y = smoothedGyroData[i].Roll
	}
	for i := 0; i < len(smoothedGyroData); i++ {
		pitchData[i].X = float64(smoothedGyroData[i].Timestamp - minX)
		pitchData[i].Y = smoothedGyroData[i].Pitch
	}
	for i := 0; i < len(smoothedGyroData); i++ {
		yawData[i].X = float64(smoothedGyroData[i].Timestamp - minX)
		yawData[i].Y = smoothedGyroData[i].Yaw
	}
	for i := 0; i < len(smoothedGyroData); i++ {
		testingData[i].X = float64(smoothedGyroData[i].Timestamp - minX)
		testingData[i].Y = smoothedGyroData[i].Roll + smoothedGyroData[i].Yaw
	}

	rollLine, err := plotter.NewLine(rollData)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rollLine.LineStyle.Width = vg.Points(1)
	rollLine.LineStyle.Color = color.RGBA{B: 255, A: 255}

	pitchLine, err := plotter.NewLine(pitchData)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pitchLine.LineStyle.Width = vg.Points(1)
	pitchLine.LineStyle.Color = color.RGBA{G: 255, A: 255}

	yawLine, err := plotter.NewLine(yawData)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	yawLine.LineStyle.Width = vg.Points(1)
	yawLine.LineStyle.Color = color.RGBA{R: 255, A: 255}

	testingLine, err := plotter.NewLine(testingData)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	testingLine.LineStyle.Width = vg.Points(1)
	testingLine.LineStyle.Color = color.RGBA{R: 255, B: 255, A: 255}

	p, err := plot.New()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p.Title.Text = "Gyro Data"
	p.X.Label.Text = "t"
	p.Y.Label.Text = "Gyro"
	p.Add(plotter.NewGrid())
	p.Add(rollLine)
	p.Add(pitchLine)
	p.Add(yawLine)
	p.Add(testingLine)
	p.Legend.Add("Roll", rollLine)
	p.Legend.Add("Pitch", pitchLine)
	p.Legend.Add("Yaw", yawLine)
	p.Legend.Add("Testing", testingLine)

	wt, err := p.WriterTo(vg.Inch*16, vg.Inch*4, "jpg")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = wt.WriteTo(w)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func smooth(data []*orientation) []*orientation {
	startTime := data[0].Timestamp

	newData := make([]*orientation, data[len(data)-1].Timestamp-startTime)
	for i := 0; i < len(data)-1; i++ {
		fmt.Println("i:", i)
		unitCount := data[i+1].Timestamp - data[i].Timestamp
		//newData[data[i].Timestamp] = &rawOrientation{Roll: data[i].Roll, Pitch: data[i].Pitch, Yaw: data[i].Yaw}
		for j := data[i].Timestamp - startTime; j <= data[i].Timestamp-startTime+(unitCount/2); j++ {
			newData[j] = &orientation{Roll: data[i].Roll, Pitch: data[i].Pitch, Yaw: data[i].Yaw, Timestamp: j + startTime}
			fmt.Println("j:", j, ":", newData[j].Timestamp)
		}
		for j := data[i].Timestamp - startTime + (unitCount / 2) + 1; j < data[i].Timestamp-startTime+unitCount; j++ {
			newData[j] = &orientation{Roll: data[i+1].Roll, Pitch: data[i+1].Pitch, Yaw: data[i+1].Yaw, Timestamp: j + startTime}
			fmt.Println("j:", j, ":", newData[j].Timestamp)
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
