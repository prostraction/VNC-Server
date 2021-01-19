#ifndef CONTROL_H
#define CONTROL_H

#include <Windows.h>
#include <WinUser.h>
#include <stdio.h>
#include <stdlib.h>
//#include <malloc.h>

void MouseMove(int x, int y);
void MouseClick(int button);
void KeyboardPress(int button);
void KeyboardRelease(int button);

#endif