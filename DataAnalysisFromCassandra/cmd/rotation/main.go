package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gonum/matrix/mat64"
	"os"
	"strconv"
	"time"
)

type acceleration struct {
	Time int64
	X    float64
	Y    float64
	Z    float64
}

type rotation struct {
	Time   int64
	Matrix []float64
}

func main() {
	accelerationData := make([]*acceleration, 0, 10000)
	rotationData := make([]*rotation, 0, 10000)

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

	stamp, _ := strconv.Atoi(os.Args[2])
	var time int64
	var x float64
	var y float64
	var z float64
	iter := session.Query(`SELECT time, x, y, z FROM trainingacceleration WHERE userid='AndrzejS' AND starttime=?`, int64(stamp)).Iter()
	for iter.Scan(&time, &x, &y, &z) {
		accelerationData = append(accelerationData, &acceleration{time, x, y, z})
	}

	tab := make([]float64, 9, 9)

	iter = session.Query(`SELECT time, a0, a1, a2, b0, b1, b2, c0, c1, c2 FROM trainingrotation WHERE userid='AndrzejS' AND starttime=?`, int64(stamp)).Iter()

	for iter.Scan(&time, &tab[0], &tab[1], &tab[2], &tab[3], &tab[4], &tab[5], &tab[6], &tab[7], &tab[8]) {
		newTab := make([]float64, 9, 9)
		for i := 0; i < 9; i++ {
			newTab[i] = tab[i]
		}
		rotationData = append(rotationData, &rotation{time, newTab})
	}

	starttime := rotationData[0].Time
	for i := 0; i < len(rotationData); i++ {
		rotationData[i].Time = rotationData[i].Time - starttime
		accelerationData[i].Time = accelerationData[i].Time - starttime
	}

	accelerationGlobal := make([]*mat64.Vector, 0, len(rotationData))
	for i := 0; i < len(rotationData); i++ {
		temp := mat64.NewVector(3, nil)
		accel := make([]float64, 3, 3)
		accel[0] = accelerationData[i].X
		accel[1] = accelerationData[i].Y
		accel[2] = accelerationData[i].Z
		rotMatrix := mat64.NewDense(3, 3, rotationData[i].Matrix)
		temp.MulVec(rotMatrix, mat64.NewVector(3, accel))
		accelerationGlobal = append(accelerationGlobal, temp)
		fmt.Println("Before:", accel[0], accel[1], accel[2])
		fmt.Println("After:", temp.At(0, 0), temp.At(1, 0), temp.At(2, 0))
	}
}
