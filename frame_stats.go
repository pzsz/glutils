package glutils

import (
	"fmt"
)

const FPS_ARRAY_SIZE = 60
const FRAME_TIME_ARRAY_SIZE = 60

type FrameStats struct {
	Frames int64
	FPS    int

	CurrentFrameTime int
	FrameTime    [FRAME_TIME_ARRAY_SIZE]int64
	AvgFrameTime int64
	MaxFrameTime int64
	MinFrameTime int64

	FPSCounter int
	SecondCounter int64
}

func (self *FrameStats) FrameFinished(timeStep int64) (SecondPassed bool) {
	self.Frames++
	self.FPSCounter++
	SecondPassed = false

	self.SecondCounter += timeStep
	if self.SecondCounter > 1000 {
		self.FPS = self.FPSCounter

		self.SecondCounter = 0
		self.FPSCounter = 0
		SecondPassed = true
	}

	self.CurrentFrameTime = (self.CurrentFrameTime+1) % FRAME_TIME_ARRAY_SIZE
	self.FrameTime[self.CurrentFrameTime] = timeStep

	self.CalcFrameTimeStats()
	return
}

func (self *FrameStats) CalcFrameTimeStats() {
	var frameTimeAcu int64 = 0
	var frames int64 = 0
	self.MaxFrameTime = 0
	self.MinFrameTime = 1000000

	for i:=0 ; i < FRAME_TIME_ARRAY_SIZE; i++ {
		cur := self.FrameTime[i]
		if self.MaxFrameTime < cur {self.MaxFrameTime = cur}
		if cur < self.MinFrameTime {self.MinFrameTime = cur}

		if cur != 0 {
			frameTimeAcu += cur
			frames++
		}
	}
	if frames != 0 {
		self.AvgFrameTime = frameTimeAcu / frames
	} else {
		self.AvgFrameTime = 0
	}
}

func (self *FrameStats) String() string {
	return fmt.Sprintf("FPS: %d, Frame time (avg,max,min): %d %d %d",
		self.FPS,
		self.AvgFrameTime,
		self.MaxFrameTime,
		self.MinFrameTime)
}