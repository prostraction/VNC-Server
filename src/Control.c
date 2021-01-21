#include "Control.h"

void MouseMove(int x, int y) {
	SetCursorPos(x, y);
}
void MouseClick(int button, int x, int y, int up) {
	INPUT Inputs[1] = {0};
	Inputs[0].type = INPUT_MOUSE;
	Inputs[0].mi.dx = x;
	Inputs[0].mi.dy = y;
	
	if (button == 1) {
		// LCM
		if (up == 1) {
			// Hold
			Inputs[0].mi.dwFlags = MOUSEEVENTF_LEFTDOWN;
		} else {
			// Released
			Inputs[0].mi.dwFlags = MOUSEEVENTF_LEFTUP;
		}
	} else if (button == 4){
		// PCM
		if (up == 1) {
			// Hold
			Inputs[0].mi.dwFlags = MOUSEEVENTF_RIGHTDOWN;
		} else {
			// Released
			Inputs[0].mi.dwFlags = MOUSEEVENTF_RIGHTUP;
		}
	}
	SendInput(1, Inputs, sizeof(INPUT));
}
void KeyboardPress(int button) {
	
}
void KeyboardRelease(int button){
	
}