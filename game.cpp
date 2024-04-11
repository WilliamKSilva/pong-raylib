#include <iostream> 
#include "raylib.h"

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

int main() {
    Window::init();

    while (!WindowShouldClose()) {
        BeginDrawing();
            ClearBackground(WHITE);
            DrawText("Hello", Window::width / 2, Window::height / 2, 31, BLACK);
        EndDrawing();
    }

    return 0;
}