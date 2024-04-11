#include <iostream> 
#include "raylib.h"

// Application level imports
#include "player.h"
#include "ball.h"

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
    Player player = Player({x: 100, y: 30}, BLACK);
    Ball ball = Ball({x: Window::width / 2, y: Window::height / 2}, 20.0, BLACK);

    while (!WindowShouldClose()) {
        player.move();
        ball.bounce();

        BeginDrawing();
            ClearBackground(WHITE);
            player.render();
            ball.render();
        EndDrawing();
    }

    return 0;
}