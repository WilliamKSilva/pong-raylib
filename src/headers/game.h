#include "raylib.h"
#ifndef GAME_H 
#define GAME_H 

#pragma once

class Window {
    public:
        static const int width = 1920;
        static const int height = 1080;
        static const int fps = 60;

        static void init() {
            InitWindow(width, height, "Pong");
            SetTargetFPS(fps);
        };
};

#endif