#ifndef IMAGE_H
#define IMAGE_H

#include <Windows.h>
#include <WinUser.h>
#include <stdio.h>
#include <stdlib.h>
#include <malloc.h>

void clean_custom_mem(unsigned char* screen, int length);

unsigned char* screenCapture(int x, int y, int w, int h);
struct CopyRectStruct CopyRectangle(int x, int y, int w, int h);
#endif