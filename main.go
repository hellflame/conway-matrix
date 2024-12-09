package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func main() {
	go func() {
		window := new(app.Window)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	args := parseArgs()
	if args == nil {
		return nil
	}
	ops := new(op.Ops)
	w := args.style.width
	h := args.style.height
	var seed int64
	if args.seed == 0 {
		seed = time.Now().UnixMilli()
	} else {
		seed = int64(args.seed)
	}
	fmt.Printf("seed: %d\nwidth: %d\nheight: %d\n", seed, w, h)
	boxes := NewBoxes(w, h)
	RandomizeBoxes(boxes, seed)
	cnt := 0
	total := args.iterations
	s := args.style
	for {
		switch e := window.Event().(type) {

		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			ops.Reset()
			// for _, l := range boxes {
			// 	fmt.Println(l)
			// }
			boxes = evolveLife(boxes)
			renderBoxes(boxes, s.boxSize, s.boxGap, s.boxRound, ops)
			e.Frame(ops)
			go func() {
				cnt += 1
				if cnt < total || total == 0 {
					e.Source.Execute(op.InvalidateCmd{})
					time.Sleep(time.Millisecond * 200)
				}
			}()
		}
	}
}

func NewBoxes(width, height int) [][]int {
	boxes := [][]int{}
	for h := 0; h < height; h++ {
		boxes = append(boxes, make([]int, width))
	}
	return boxes
}

func RandomizeBoxes(boxes [][]int, seed int64) {
	rnd := rand.New(rand.NewSource(seed))
	width := len(boxes[0])
	height := len(boxes)

	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			boxes[r][c] = rnd.Intn(2)
		}
	}
}

func GliderBoxes(boxes [][]int) {
	boxes[10][10] = 1
	boxes[10][11] = 1
	boxes[10][12] = 1
	boxes[11][10] = 1
	boxes[12][11] = 1
}

func evolveLife(boxes [][]int) [][]int {
	// println("evolve")
	width := len(boxes[0])
	height := len(boxes)
	tmp := NewBoxes(width, height)

	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			neighbor :=
				boxes[(height+r-1)%height][(width+c-1)%width] +
					boxes[(height+r-1)%height][c] +
					boxes[(height+r-1)%height][(c+1)%width] +

					boxes[r][(width+c-1)%width] +
					boxes[r][(c+1)%width] +

					boxes[(r+1)%height][(width+c-1)%width] +
					boxes[(r+1)%height][c] +
					boxes[(r+1)%height][(c+1)%width]
			origin := boxes[r][c]

			if origin == 1 {
				if neighbor < 2 || neighbor > 3 {
					origin = 0
				}
			} else {
				if neighbor == 3 {
					origin = 1
				}
			}
			tmp[r][c] = origin
		}
	}
	return tmp
}

func renderBoxes(boxes [][]int, blockSize, space, round int, ops *op.Ops) {
	width := len(boxes[0])
	height := len(boxes)

	area := clip.Rect(image.Rect(0, 0, width*(blockSize+space), height*(blockSize+space)))
	defer area.Push(ops).Pop()

	for r := 0; r < height; r++ {
		rowOffset := op.Offset(image.Pt(0, r*(space+blockSize))).Push(ops)
		for c := 0; c < width; c++ {
			colOffset := op.Offset(image.Pt(c*(space+blockSize), 0)).Push(ops)

			b := clip.RRect{Rect: image.Rect(0, 0, blockSize, blockSize),
				SE: round, SW: round, NW: round, NE: round}.Push(ops)

			var dye color.NRGBA
			if boxes[r][c] > 0 {
				dye = color.NRGBA{R: 0, G: 125, B: 200, A: 255}
			} else {
				dye = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			}
			paint.ColorOp{Color: dye}.Add(ops)
			paint.PaintOp{}.Add(ops)

			b.Pop()
			colOffset.Pop()
		}
		rowOffset.Pop()
	}
}
