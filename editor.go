package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	term "github.com/nsf/termbox-go"
)

const (
	CLOSE = iota // C-q
	SAFE  = iota // C-s
	YANK  = iota // C-y
	CUT   = iota // C-k
	DEL   = iota // C-d
	BACK  = iota // BACK
	UP    = iota
	DOWN  = iota
	RIGHT = iota
	LEFT  = iota
	ENTER = iota
	CHAR  = iota // any non control key
)

type Action struct {
	Action int
	Ch     rune
}

func PollKeys(c chan Action) {
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyCtrlQ:
				//term.Sync()
				c <- Action{CLOSE, ' '}
			case term.KeyCtrlS:
				//term.Sync()
				c <- Action{SAFE, ' '}
			case term.KeyCtrlY:
				//term.Sync()
				c <- Action{YANK, ' '}
			case term.KeyCtrlK:
				//term.Sync()
				c <- Action{CUT, ' '}
			case term.KeyCtrlD:
				//term.Sync()
				c <- Action{DEL, ' '}
			case term.KeyBackspace:
			case term.KeyBackspace2:
				//term.Sync()
				c <- Action{BACK, ' '}
			case term.KeyArrowDown:
				//term.Sync()
				c <- Action{DOWN, ' '}
			case term.KeyArrowUp:
				//term.Sync()
				c <- Action{UP, ' '}
			case term.KeyArrowLeft:
				//term.Sync()
				c <- Action{LEFT, ' '}
			case term.KeyArrowRight:
				//term.Sync()
				c <- Action{RIGHT, ' '}
			case term.KeyEnter:
				//term.Sync()
				c <- Action{ENTER, ' '}
			case term.KeySpace:
				//term.Sync()
				c <- Action{CHAR, ' '}
			default:
				//term.Sync()
				c <- Action{CHAR, ev.Ch}
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}

type Editor struct {
	col   int
	row   int
	index int
	buf   []rune
}

type Stat struct {
	sync.Mutex
	editor Editor
}

func main() {
	must(term.Init())
	defer term.Close()

	keys := make(chan Action)
	go PollKeys(keys)

	editor := Editor{}
	editor.buf = make([]rune, 0)
	editor.index = -1
	s := &Stat{editor: editor}
	go render(s)

	for {
		key := <-keys
		s.Lock()
		switch key.Action {
		case CLOSE:
			term.Close()
			os.Exit(0)
		case ENTER:
			s.editor.buf = append(s.editor.buf, '\r')
			s.editor.buf = append(s.editor.buf, '\n')
			s.editor.index += 2
			s.editor.row++
			s.editor.col = 0
			break
		case CHAR:
			s.editor.buf = append(s.editor.buf, key.Ch)
			s.editor.index++
			s.editor.col++
		case BACK:
			if s.editor.index > 0 {
				if s.editor.buf[s.editor.index] == '\n' {
					//line break
					//remove \r \n
					s.editor.buf = s.editor.buf[:len(s.editor.buf)-2]
					s.editor.index -= 2

					// dertermine col of new line
					i := 0
					for ; checkI(i, s.editor.index); i++ {
						if s.editor.buf[s.editor.index-i] == '\n' {
							break
						}
					}
					s.editor.col = i
					s.editor.row--
					//s.editor.buf = s.editor.buf[:s.editor.index]
					//s.editor.buf = append(s.editor.buf, ' ')

				} else {
					//non newline break
					s.editor.buf = s.editor.buf[:len(s.editor.buf)-1]
					s.editor.col--
					s.editor.index--
				}

			} else if s.editor.index == 0 {
				s.editor.buf = make([]rune, 0)
				s.editor.index = -1
			}
		case RIGHT:
			s.editor.col++
		case LEFT:
			s.editor.col--
			if s.editor.col < 0 {
				s.editor.col = 0
			}
		}
		s.Unlock()
	}
}

func checkI(i, index int) bool {
	return index-i > 0
}

func render(s *Stat) {
	for {
		s.Lock()
		term.Sync()

		cursorIndex := s.editor.index + 1

		if len(s.editor.buf) < 1 || s.editor.index < 0 {
			fmt.Println("row:", s.editor.row, " col:", s.editor.col, " index:", s.editor.index, " cursor:", cursorIndex)
			s.Unlock()
			time.Sleep(50 * time.Millisecond)
			continue
		}

		out := make([]rune, cursorIndex)
		copy(out, s.editor.buf)
		for len(out) <= cursorIndex {
			out = append(out, ' ')
		}
		out[cursorIndex] = '▁'

		fmt.Println("row:", s.editor.row, " col:", s.editor.col, " index:", s.editor.index, " cursor:", cursorIndex, " buf[i]:", s.editor.buf[s.editor.index])
		fmt.Print(string(out))
		// fmt.Println("\n---")
		// fmt.Print(string(s.editor.buf))
		// fmt.Println("\n---")
		// fmt.Println(s.editor.buf)
		s.Unlock()
		time.Sleep(50 * time.Millisecond)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
