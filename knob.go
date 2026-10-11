package main

// a model of the knob that controls one axis of stylus movement
type knob struct {
	leftSpinKey      string
	rightSpinKey     string
	leftSpinKeyDown  bool
	rightSpinKeyDown bool
	loc              int
	size             int
	knobDirection    int // anticlockwise == -1, clockwise == 1, none == 0
}

var KNOB_FRAMES = []string{"▓▓░░▓▓\n▓▓  ▓▓\n▓▓▓▓▓▓", "▓▓▓▓░░\n▓▓  ▓▓\n▓▓▓▓▓▓", "▓▓▓▓▓▓\n▓▓  ░░\n▓▓▓▓▓▓", "▓▓▓▓▓▓\n▓▓  ▓▓\n▓▓▓▓░░", "▓▓▓▓▓▓\n▓▓  ▓▓\n▓▓░░▓▓", "▓▓▓▓▓▓\n▓▓  ▓▓\n░░▓▓▓▓", "▓▓▓▓▓▓\n░░  ▓▓\n▓▓▓▓▓▓", "░░▓▓▓▓\n▓▓  ▓▓\n▓▓▓▓▓▓"}

// update the knob state and return it's location delta.
// It's position value will only ever change by 1 per call of this func.
func (k *knob) gameStep() int {
	oldLoc := k.loc
	modification := k.knobDirection
	k.loc = max(0, min(k.loc+modification, k.size-1))

	return k.loc - oldLoc
}

func (k *knob) handleKeyPress(key string) {
	if key == k.leftSpinKey {
		k.leftSpinKeyDown = true
	} else if key == k.rightSpinKey {
		k.rightSpinKeyDown = true
	}

	// pressing both keys deadlocks the knob
	if k.leftSpinKeyDown && k.rightSpinKeyDown {
		k.knobDirection = NO_SPIN
	} else if k.leftSpinKeyDown {
		k.knobDirection = ANTICLOCKWISE_SPIN
	} else if k.rightSpinKeyDown {
		k.knobDirection = CLOCKWISE_SPIN
	} else {
		k.knobDirection = NO_SPIN
	}
}

func (k *knob) handleKeyRelease(key string) {
	if key != k.leftSpinKey && key != k.rightSpinKey {
		return
	}

	if key == k.leftSpinKey {
		k.leftSpinKeyDown = false
	} else if key == k.rightSpinKey {
		k.rightSpinKeyDown = false
	}

	if k.leftSpinKeyDown {
		k.knobDirection = ANTICLOCKWISE_SPIN
	} else if k.rightSpinKeyDown {
		k.knobDirection = CLOCKWISE_SPIN
	} else {
		k.knobDirection = NO_SPIN
	}
}

func (k *knob) getView() string {
	return KNOB_FRAMES[k.loc%len(KNOB_FRAMES)]
}


// ▓▓░░▓▓
// ▓▓  ▓▓
// ▓▓▓▓▓▓

// ▓▓▓▓░░
// ▓▓  ▓▓
// ▓▓▓▓▓▓

// ▓▓▓▓▓▓
// ▓▓  ░░
// ▓▓▓▓▓▓

// ▓▓▓▓▓▓
// ▓▓  ▓▓
// ▓▓▓▓░░

// ▓▓▓▓▓▓
// ▓▓  ▓▓
// ▓▓░░▓▓

// ▓▓▓▓▓▓
// ▓▓  ▓▓
// ░░▓▓▓▓

// ▓▓▓▓▓▓
// ░░  ▓▓
// ▓▓▓▓▓▓

// ░░▓▓▓▓
// ▓▓  ▓▓
// ▓▓▓▓▓▓
