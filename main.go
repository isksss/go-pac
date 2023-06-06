package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	tsize "github.com/kopoli/go-terminal-size"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	delay = time.Second / 10 // 進む速度
)

var (
	pacmanPos       = 1
	pacmanRow       = 1
	lastWidth       = -1
	lastHeight      = -1
	totalDistance   = 0
	totalStarsEaten = 0
	stars           = make(map[[2]int]struct{})
)

func clearLine() {
	cmd := exec.Command("tput", "el")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func moveCursor(x, y int) {
	fmt.Printf("\033[%d;%dH", y+1, x+1)
}

func hideCursor() {
	cmd := exec.Command("tput", "civis")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func showCursor() {
	cmd := exec.Command("tput", "cnorm")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func clearScreen(height int) {
	for i := 0; i < height; i++ {
		moveCursor(0, i)
		clearLine()
	}
}

func initializePacman(width, height int) {
	pacmanPos = (width - 2) / 2
	pacmanRow = (height - 3) / 2
}

func drawHorizontalBorders(width int) {
	fmt.Print("+")
	for i := 0; i < width-2; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")
}

func drawVerticalBorders(height, width int) {
	for i := 3; i < height; i++ {
		moveCursor(0, i)
		fmt.Print("|")
		moveCursor(width-1, i)
		fmt.Println("|")
	}
}

func drawBorder(width, height int) {
	if lastWidth != width || lastHeight != height {
		clearScreen(height)

		moveCursor(0, 2) // 描画位置を1段下にする
		drawHorizontalBorders(width)

		drawVerticalBorders(height, width)

		moveCursor(0, height)
		drawHorizontalBorders(width)

		lastWidth = width
		lastHeight = height
	}
}

func drawPacman() {
	moveCursor(pacmanPos+1, pacmanRow+1) // 描画位置を1段下にする
	fmt.Print("@")
}

func erasePacman() {
	moveCursor(pacmanPos+1, pacmanRow+1) // 描画位置を1段下にする
	fmt.Print(" ")
}

func eraseStar(x, y int) {
	_, exists := stars[[2]int{x, y}]
	if exists {
		moveCursor(x+1, y+1)
		fmt.Print(" ")
		delete(stars, [2]int{x, y})
		totalStarsEaten++
	}
}

func updatePosition(width, height, direction int) {
	distance := rand.Intn(3) + 1 // 1から3の間でランダムに移動距離を決定
	switch direction {
	case 0: // 上
		for i := 0; i < distance; i++ {
			if pacmanRow > 0 {
				eraseStar(pacmanPos, pacmanRow)
				erasePacman()
				pacmanRow--
				totalDistance++
			}
		}
	case 1: // 下
		for i := 0; i < distance; i++ {
			if pacmanRow < height-4 {
				eraseStar(pacmanPos, pacmanRow+1)
				erasePacman()
				pacmanRow++
				totalDistance++
			}
		}
	case 2: // 左
		for i := 0; i < distance; i++ {
			if pacmanPos > 0 {
				eraseStar(pacmanPos-1, pacmanRow)
				erasePacman()
				pacmanPos--
				totalDistance++
			}
		}
	case 3: // 右
		for i := 0; i < distance; i++ {
			if pacmanPos < width-2 {
				eraseStar(pacmanPos+1, pacmanRow)
				erasePacman()
				pacmanPos++
				totalDistance++
			}
		}
	}
	drawPacman()
}

func update() {
	width, height, _ := terminal.GetSize(0)
	width -= 2                // 左右の境界線分を除外
	height -= 3               // 上下の境界線と最終行分を除外
	direction := rand.Intn(4) // 0: 上, 1: 下, 2: 左, 3: 右
	updatePosition(width, height, direction)
}

func drawDistance() {
	moveCursor(0, 0)
	fmt.Printf("Total Distance: %d", totalDistance)
}

func spawnStars(width, height int) {
	numStars := rand.Intn(8) + 3 // 3から10の間で星の数をランダムに決定
	for i := 0; i < numStars; i++ {
		starX := rand.Intn(width-3) + 1  // 左右の境界を除外
		starY := rand.Intn(height-4) + 2 // 上下の境界と最終行を除外
		moveCursor(starX, starY)
		fmt.Print("*")
		stars[[2]int{starX, starY}] = struct{}{}
	}
}

func drawStarsEaten() {
	moveCursor(0, 0)
	fmt.Printf("Total Stars Eaten: %d", totalStarsEaten)
}

func main() {
	rand.Seed(time.Now().UnixNano()) // ランダムシードを設定
	hideCursor()                     // カーソル非表示
	defer showCursor()               // プログラム終了時にカーソルを再表示

	starTicker := time.NewTicker(5 * time.Second) // 10秒ごとのタイマーを作成

	for {
		var s tsize.Size
		var err error
		s, err = tsize.GetSize()

		w := s.Width
		h := s.Height

		if err != nil {
			panic(err)
		}

		if lastWidth != w || lastHeight != h {
			initializePacman(w, h)
		}
		drawBorder(w, h)
		updatePosition(w, h, rand.Intn(4)) // 0: 上, 1: 下, 2: 左, 3: 右
		drawStarsEaten()

		select {
		case <-starTicker.C: // 10秒ごとに星を生成
			spawnStars(w, h)
		default:
		}

		time.Sleep(delay)
	}
}
