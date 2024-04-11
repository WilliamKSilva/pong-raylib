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
    Player player = Player({x: 100, y: 30}, BLACK, "Player");
    Ball ball = Ball({x: Window::width / 2, y: Window::height / 2}, 20.0, BLACK);

    while (!WindowShouldClose()) {
        ball.check_collision({
            x: player.pos.x,
            y: player.pos.y,
            width: (float)player.width,
            height: (float)player.height
        }, player.name);

        ball.check_out_of_bounds();

        Scored scored = ball.check_scored();
        if (scored.playerScored) {
            player.score++;
        }

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