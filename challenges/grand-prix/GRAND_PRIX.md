# Multithreaded Grand Prix

This is an implementation of a a multithreaded grand prix simulator where every car will be an independent entity. 
It will be simulating an N laps race.


## Technical Features

- The race's lap is introduced by the user as well as the number of cars.
- Random max speed for each car.
- Emulate the increasing speed behaviour that any car has.
- If a car detect another car on his route it will slow down its speed simulating a crash.
- Each racer behaviour is a separated thread.
- Cars threads use the same map or city layout data structure resource.
- Displays each car's speed, position, racing time and lap.
- At the end, top 3 winners is displayed, if there is a tie for 3rd more cars will be shown.


## Build

### On __Linux__ or __MacOs__

```bash
# To compile the go proyect
make build
# This will install every package necessary for the build and create a binary file that can be run with ./car

# To remove the executable created after compilation.
make clean
```

### On __Windows__

```bash
# To get dependencies
go get -v

# To build the proyect
go build
```

## Run

### On __Linux__

```bash
make run
# or
make
```
### On __Windows__

```bash
go run game.go
```


## Info

### Language
- [Go](https://golang.org/)

### Packages
We used the external packages image and faiface/pixel.

The first 3 were used to edit the images and text on them. Then the last 2 were used to import font types and colors
 
-	"github.com/faiface/pixel/text"
-	"github.com/faiface/pixel"
-	"github.com/faiface/pixel/pixelgl"

- "golang.org/x/image/font/basicfont"
-	"golang.org/x/image/colornames"

### Structures
#### Car
-	id int - An identifier for the cars
-	currentLap int - Current lap of the cars
-	speed chan int - Channel for the goroutine to get the speed
-	boundPos chan int - Channel for the gourutine to get the position in the Y bounds of the car
-	sprite *pixel.Sprite - The sprite(image) of the car
-	mat pixel.Matrix  - Matrix of the sprite
-	seed rand.Source - 
-	botBound int
-	topBound int
-	crashing bool
-	finished bool
-	position int
-	finalPosition int
-	timeElapsed time.Duration
  
