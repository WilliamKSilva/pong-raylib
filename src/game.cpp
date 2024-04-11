#include <iostream> 
#include "raylib.h"

// Application level imports
#include "player.h"
#include "ball.h"
#include "game.h"
#include "enemy.h"

int main() {
    Window::init();
    Player player = Player({x: 100, y: 30}, BLACK, "Player");
    Enemy enemy = Enemy({x: Window::width - 100, y: 30}, BLACK, "Enemy");
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

            ball.reset_state();
        }

        player.move();
        ball.bounce();
        enemy.move(ball.pos.y);

        BeginDrawing();
            ClearBackground(WHITE);
            player.render();
            ball.render();
            enemy.render();
        EndDrawing();
    }

    return 0;
}