default: run

build:
	go get github.com/faiface/pixel
	go get github.com/faiface/mainthread
	go get -u github.com/go-gl/glfw/v3.3/glfw
	go build -v -o ./car

run: build
	./car

clean:
	$(RM) car