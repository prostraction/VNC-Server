#include "Image.h"

void clean_custom_mem(unsigned char* screen, int length)
{
	memset(screen, 0, length);
	free(screen);
	
	screen = NULL;
}

struct CopyRectStruct {
	int countRect;
	unsigned char* dib;
};


unsigned char* screenCapture(int x, int y, int w, int h)
{
    HDC hdcSource, hdcMemory;
    HBITMAP hBitMap;
    HWND hwnd;

    hwnd        = GetDesktopWindow();
    hdcSource   = GetDCEx(hwnd, NULL, 0);
    hdcMemory   = CreateCompatibleDC(hdcSource);
    hBitMap     = CreateCompatibleBitmap(hdcSource, w, h);
	
    unsigned char* dib = NULL;
    if (hBitMap) {

        SelectObject(hdcMemory, hBitMap);
		BitBlt(hdcMemory, x, y, w, h, hdcSource, 0, 0, SRCCOPY);
		
		//BITMAP bm;
        //GetObject(hBitMap, sizeof(BITMAP), (LPSTR)&bm);
        BITMAPINFO bmf = { 0 };
        bmf.bmiHeader.biSize = sizeof(bmf.bmiHeader);
		
		
        HDC dc = GetDC(NULL);
		if (!GetDIBits(dc, hBitMap, 0, 0, NULL, &bmf, DIB_RGB_COLORS))
		{
			printf ("GetDIBits err\n");
			DeleteDC(hdcSource);
			DeleteDC(hdcMemory);
			ReleaseDC(NULL, hdcSource);
			ReleaseDC(hwnd, hdcMemory);
			DeleteObject(hBitMap);
			return dib;
		}
		

        dib = malloc(bmf.bmiHeader.biSizeImage);
		
        bmf.bmiHeader.biSize = sizeof(BITMAPINFOHEADER);
        bmf.bmiHeader.biBitCount = 32;
        bmf.bmiHeader.biCompression = BI_RGB;
		
        // if (!GetDIBits(dc, hBitMap, 0, bmf.bmiHeader.biHeight, (LPVOID)dib, &bmf, DIB_RGB_COLORS))
		if (!GetBitmapBits(hBitMap, w * h  * 4, dib)) 
		{
			free(dib);
			DeleteDC(hdcSource);
			DeleteDC(hdcMemory);
			ReleaseDC(NULL, hdcSource);
			ReleaseDC(hwnd, hdcMemory);
			DeleteObject(hBitMap);
			return dib;
		}
        ReleaseDC(NULL, dc);
    }
    else
    {
        printf("BitMap err\n");
    }
		
    DeleteDC(hdcSource);
	DeleteDC(hdcMemory);
    ReleaseDC(NULL, hdcSource);
	ReleaseDC(hwnd, hdcMemory);
	DeleteObject(hBitMap);
	
	return dib;
}
