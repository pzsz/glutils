package glutils

import (
	"container/list"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/pzsz/gl"
	"flag"
	"os"
	"runtime/pprof"
)

var FLAG_profile *bool = flag.Bool("profile", false, "start profiling to gowar.prof")

type AppState interface {
	Setup(man *AppStateManager)
	Destroy()
	Pause()
	Resume()
	Process(time_step float32)

	OnKeyDown(key *sdl.Keysym)
	OnKeyUp(key *sdl.Keysym)

	OnMouseMove(x, y, dx, dy float32)
	OnMouseClick(x, y float32, button int, down bool)

	OnSdlEvent(event *sdl.Event)

	OnViewportResize(x, y float32)
}

type AppStateManager struct {
	StateStack              *list.List
	Screen                  *sdl.Surface
	LastMouseX, LastMouseY  float32    
	MouseSampleTaken        bool

	FPSMouseModeEnabled     bool
}

var AppStateManagerInstance *AppStateManager = &AppStateManager{StateStack: list.New()}

func GetManager() *AppStateManager {
	return AppStateManagerInstance
}

func (self *AppStateManager) Setup(state AppState, caption string) {
	if *FLAG_profile {
		pfile, _ := os.Create("gowar.prof")
		pprof.StartCPUProfile(pfile)
	}

	// Lock, so we got no GL-in-wrong-thread foolery
	sdl.Init(sdl.INIT_VIDEO)

	self.Screen = sdl.SetVideoMode(800, 600, 32, sdl.OPENGL|sdl.RESIZABLE)

	gl.Init()

	if self.Screen == nil {
		sdl.Quit()
		panic("Couldn't set 300x300 GL video mode: " + sdl.GetError() + "\n")
	}

	sdl.WM_SetCaption(caption, caption)

	Setup()

	GetViewport().SetScreenSize(float32(self.Screen.W), float32(self.Screen.H))

	self.Push(state)
}

func (self *AppStateManager) FPSMouseMode(on bool) {
	self.FPSMouseModeEnabled = on
	if on {
		sdl.ShowCursor(0)
	} else {
		sdl.ShowCursor(1)
	}
}

func (self *AppStateManager) Push(state AppState) {
	e := self.StateStack.Back()
	if e != nil {
		prev := e.Value.(AppState)
		prev.Pause()
	}

	self.StateStack.PushBack(state)
	state.Setup(self)
}

func (self *AppStateManager) Pop() (ret AppState) {
	e := self.StateStack.Back()
	if e != nil {
		ret = e.Value.(AppState)
		ret.Destroy()
		self.StateStack.Remove(e)
	}

	ne := self.StateStack.Back()
	if ne != nil {
		nstate := ne.Value.(AppState)
		nstate.Resume()
	}

	return
}

func (self *AppStateManager) Replace(state AppState) {
	e := self.StateStack.Back()
	if e != nil {
		ret := e.Value.(AppState)
		ret.Destroy()
		self.StateStack.Remove(e)
	}

	self.StateStack.PushBack(state)
	state.Setup(self)
}

func (self *AppStateManager) GetRunningState() AppState {
	e := self.StateStack.Back()
	if e != nil {
		return e.Value.(AppState)
	}
	return nil
}

func (self *AppStateManager) Process(time_step float32) {
	e := self.StateStack.Back()
	if e != nil {
		state := e.Value.(AppState)
		state.Process(time_step)
	}
}

func (self *AppStateManager) HandleEvents() (done bool) {
	done = false

	running_state := self.GetRunningState()

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			done = true
			break
		case *sdl.ResizeEvent:
			re := event.(*sdl.ResizeEvent)
			self.Screen = sdl.SetVideoMode(int(re.W), int(re.H), 16,
				sdl.OPENGL|sdl.RESIZABLE)
			if self.Screen != nil {
				GetViewport().SetScreenSize(float32(self.Screen.W), float32(self.Screen.H))

				if running_state != nil {
					running_state.OnViewportResize(float32(self.Screen.W), float32(self.Screen.H))
				}

			} else {
				panic("we couldn't set the new video mode??")
			}
			break

		case *sdl.KeyboardEvent:
			if running_state != nil {
				kevent := event.(*sdl.KeyboardEvent)
				if kevent.State == 1 {
					running_state.OnKeyDown(&kevent.Keysym)
				} else {
					running_state.OnKeyUp(&kevent.Keysym)
				}
			}
			break
		case *sdl.MouseMotionEvent:
			if running_state != nil {
				mevent := event.(*sdl.MouseMotionEvent)
				dx, dy := float32(0), float32(0)
				fx, fy := float32(mevent.X), float32(mevent.Y)

				if self.MouseSampleTaken {
					dx = fx - self.LastMouseX
					dy = fy - self.LastMouseY
				} else {
					self.MouseSampleTaken = true
				}

				running_state.OnMouseMove(
					float32(mevent.X),
					float32(mevent.Y), dx, dy)

				if self.FPSMouseModeEnabled {
					sdl.EventState(sdl.MOUSEMOTION, sdl.IGNORE)
					sdl.WarpMouse(int(self.Screen.W/2), int(self.Screen.H/2))
					sdl.EventState(sdl.MOUSEMOTION, sdl.ENABLE)
					self.LastMouseX = float32(self.Screen.W/2)
					self.LastMouseY = float32(self.Screen.H/2)
				} else {
					self.LastMouseX = fx
					self.LastMouseY = fy
				}
			}
			break
		case *sdl.MouseButtonEvent:
			if running_state != nil {
				mevent := event.(*sdl.MouseButtonEvent)
				running_state.OnMouseClick(float32(mevent.X),
					float32(mevent.Y),
					int(mevent.Button),
					mevent.State == 1)
			}
			break
		default:
			if running_state != nil {
				running_state.OnSdlEvent(&event)
			}
			break
		}
	}
	return
}

func (self *AppStateManager) RunLoop() {
	done := false
	last_ticks := sdl.GetTicks()
	for !done {
		//		for ;sdl.GetTicks()-last_ticks < 10;time.Sleep(1) {}
		current_ticks := sdl.GetTicks()
		timeStepMS := current_ticks - last_ticks
		time_step := float32(timeStepMS) / 1000.0

		done = self.HandleEvents()
		self.Process(time_step)
		last_ticks = current_ticks
	}
}

func (self *AppStateManager) Destroy() {
	sdl.Quit()

	if *FLAG_profile {
		pprof.StopCPUProfile()
	}
}
