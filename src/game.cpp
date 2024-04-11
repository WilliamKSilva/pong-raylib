#include <iostream> 
#include "raylib.h"

// Application level imports
#include "player.h"

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
    Player* player = new Player({x: 100, y: 30}, BLACK);

    while (!WindowShouldClose()) {
        player->move();

        BeginDrawing();
            ClearBackground(WHITE);
            player->render();
        EndDrawing();
    }

    return 0;
}