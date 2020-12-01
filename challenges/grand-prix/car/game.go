package main

import (
	"bufio"
	"log"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"image"
	_ "image/png"
	"math/rand"

	"fmt"
	"os"
	"sort"
	"time"
)

type Car struct {
	id            int
	currentLap    int
	speed         chan int
	boundPos      chan int
	sprite        *pixel.Sprite
	mat           pixel.Matrix
	seed          rand.Source
	botBound      int
	topBound      int
	crashing      bool
	finished      bool
	position      int
	finalPosition int
	timeElapsed   time.Duration
}

// Code taken from "github.com/faiface/pixel" Documentation
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {

	// set the actual window
	cfg := pixelgl.WindowConfig{
		Title:  "Grand Prix",
		Bounds: pixel.R(0, 0, 1200, 900),
		VSync:  true,
	}

	// create an actual window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	//Just print "Score:""
	scoreAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	scoreTxt := text.New(pixel.V(250, 750), scoreAtlas)
	fmt.Fprintln(scoreTxt, "Score: ")

	// race info
	infoAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	infoTxt := text.New(pixel.V(500, 850), infoAtlas)

	// winners info
	winnerAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	winnerTxt := text.New(pixel.V(250, 500), winnerAtlas)

	win.SetSmooth(true)

	// loads the images
	carImg, err := loadPicture("car.png")
	circuit, err := loadPicture("race.png")

	// checks the error
	if err != nil {
		panic(err)
	}

	carLocation := 600 / totalCars
	botSpace := 80
	initSpace := carLocation

	// initialize all the cars
	for i := 0; i < totalCars; i++ {
		cars = append(cars, Car{i, 0, make(chan int), make(chan int), pixel.NewSprite(carImg, carImg.Bounds()), pixel.IM.Moved(pixel.V(float64(initSpace), 0)), rand.NewSource(time.Now().UnixNano() + int64(botSpace)), initSpace - 10, initSpace + 10, false, false, 0, 0, 100 * time.Millisecond})
		totalSpeed = append(totalSpeed, 0)
		tmpBoundsPos = append(tmpBoundsPos, 0)
		botSpace += 1 / totalCars
		initSpace += carLocation
	}

	win.Clear(colornames.Black)

	for index := range cars {
		go getCarSpeed(cars[index].speed, cars[index].seed)
		go getBoundsPos(cars[index].boundPos, cars[index].seed)
	}

	start := time.Now()

	for !win.Closed() {

		if winners >= 3 {

			win.Clear(colornames.Black)
			winnerTxt.Clear()
			infoTxt.Clear()
			scoreTxt.Clear()

			for i := 0; i < winners; i++ {
				fmt.Fprintf(winnerTxt, "Place: %d		ID Car: %d		Time: %.4v\n", finalThree[i].finalPosition, finalThree[i].id+1, finalThree[i].timeElapsed)
			}

			winnerTxt.Draw(win, pixel.IM.Scaled(winnerTxt.Orig, 3))

		} else {

			// clear the window everytime
			carPositions = carPositions[:0]
			infoTxt.Clear()

			for i := 0; i < totalCars; i++ {
				win.Clear(colornames.Black)
				lapSpeed := <-cars[i].speed
				if !cars[i].crashing {
					tmpBoundsPos[i] += <-cars[i].boundPos
				} else {
					tmpBoundsPos[i] = 5
					lapSpeed -= 5
					go decreaseSpeed(i)
				}
				totalSpeed[i] += lapSpeed

				// print winners
				if !cars[i].finished {
					fmt.Fprintf(infoTxt, "Car: %d		Lap: %d		  pos: %d		%d mp/h\n", cars[i].id+1, cars[i].currentLap, cars[i].position, lapSpeed*20)
				} else {
					fmt.Fprintf(infoTxt, "Car: %d		Place: %d		Elapsed time: %.4v\n", cars[i].id+1, cars[i].finalPosition, cars[i].timeElapsed)
				}
			}

			// check the car positions
			for i := 0; i < totalCars; i++ {
				carPositions = append(carPositions, totalSpeed[i])
			}

			sort.Ints(carPositions)

			tmpIndex := totalCars
			for i := 0; i < totalCars; i++ {
				for j := 0; j < totalCars; j++ {
					if carPositions[i] == totalSpeed[j] {
						cars[j].position = tmpIndex - i + winners
					}
				}
			}

			// print the circuit each time
			circuitSprite := pixel.NewSprite(circuit, circuit.Bounds())
			mat := pixel.IM
			mat = mat.Moved(pixel.V(400, 500))
			mat = mat.ScaledXY(pixel.V(200, 540), pixel.V(2, 5))
			circuitSprite.Draw(win, mat)

			tmpInitPos := carLocation
			for i := 0; i < totalCars; i++ {

				tmpBoundPos := tmpInitPos + tmpBoundsPos[i]
				if tmpBoundPos < 60 {
					tmpBoundPos = 70
				}else if tmpBoundPos > 650 {
					tmpBoundPos = 600
				}

				newVector := pixel.V(float64(totalSpeed[i]), float64(tmpBoundPos))
				cars[i].mat = pixel.IM.Moved(newVector)

				cars[i].botBound = tmpInitPos + tmpBoundsPos[i] - 12
				cars[i].topBound = tmpInitPos + tmpBoundsPos[i] + 12

				tmpInitPos += carLocation

				// Checks for crashing
				for j := 0; j < totalCars; j++ {
					if i != j {
						if cars[i].botBound <= cars[j].topBound &&
							cars[i].topBound >= cars[j].botBound {
							cars[i].crashing = true
						}
					}
				}

				// redraw all cars
				cars[i].sprite.Draw(win, cars[i].mat)

				// 1450
				if totalSpeed[i] > 1100 {
					totalSpeed[i] = 0
					cars[i].currentLap++

					if cars[i].currentLap == laps && !cars[i].finished {
						cars[i].finished = true
						cars[i].timeElapsed = time.Since(start)
						winners++
						cars[i].finalPosition = winners
						finalThree = append(finalThree, cars[i])
						if winners >=3{
							close(cars[i].boundPos)
							close(cars[i].speed)
						}

					}

				}
			}

			scoreTxt.Draw(win, pixel.IM.Scaled(scoreTxt.Orig, 4))
			infoTxt.Draw(win, pixel.IM.Scaled(infoTxt.Orig, 1.5))
		}

		// update the frame
		win.Update()
	}

}

var totalCars int
var cars []Car

var totalSpeed []int
var tmpBoundsPos []int

var laps int
var carPositions []int

var winners int
var finalThree []Car

func main() {

	totalCars, laps = cars_laps()
	winners = 0

	// Render 
	pixelgl.Run(run)
}

// gorutine to obtain car speed
func getCarSpeed(rnd chan int, source rand.Source) {
	r := rand.New(source)
	a := 8
	b := 12
	for {
		time.Sleep(50 * time.Millisecond)
		select {
		case <-rnd:
			return
		default:
			rnd <- a + r.Intn(b-a+1)
		}
	}
}

// gorutine to get variation at car Y bounds
func getBoundsPos(rnd chan int, source rand.Source) {
	a := -5
	b := 5
	r := rand.New(source)
	for {
		time.Sleep(50 * time.Millisecond)
		select {
		case <-rnd:
			return
		default:
			rnd <- a + r.Intn(b-a+1)
		}
	}
}

func decreaseSpeed(i int) {
	time.Sleep(800 * time.Millisecond)
	cars[i].crashing = false
}

func scanner() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()

}

func ask() (string, string) {

	fmt.Println("Welcome to Grand Prix!!!\n")
	fmt.Println("--------------------------------------")
	fmt.Print("Please enter the number of cars 3-9 ----> ")

	cars := scanner()

	fmt.Println("\n-------------------------------------")
	fmt.Print("Please enter the number of laps 2-10 ----> ")

	laps := scanner()

	return cars, laps
}

func cars_laps() (int, int) {

	carstr, lapstr := (ask())

	cars, err := strconv.Atoi(carstr)

	if cars < 3 || cars > 10 {
		log.Fatal("Invalid number of cars")
	}

	laps, err := strconv.Atoi(lapstr)

	if laps < 2 || laps > 15 {
		log.Fatal("Invalid number of laps")
	}

	if err != nil {
		log.Fatal(err, " Please enter integers")
	}

	return cars, laps
}
