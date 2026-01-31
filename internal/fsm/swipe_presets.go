package fsm

// Swipe presets (screen 1080x2400, offset 300px)
var (
	SwipeRight300 = &Swipe{
		X1: 540, Y1: 1200,
		X2: 240, Y2: 1200,
	}
	SwipeLeft300 = &Swipe{
		X1: 540, Y1: 1200,
		X2: 840, Y2: 1200,
	}
	SwipeUp300 = &Swipe{
		X1: 540, Y1: 1200,
		X2: 540, Y2: 1500,
	}
	SwipeDown300 = &Swipe{
		X1: 540, Y1: 1200,
		X2: 540, Y2: 900,
	}
)
