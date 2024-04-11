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
        // Check position of stuff
        ball.check_collision({
            x: player.pos.x,
            y: player.pos.y,
            width: (float)player.width,
            height: (float)player.height
        }, player.name);

        ball.check_collision({
            x: enemy.pos.x,
            y: enemy.pos.y,
            width: (float)enemy.width,
            height: (float)enemy.height
        }, enemy.name);

        ball.check_out_of_bounds();

        Scored scored = ball.check_scored();
        if (scored.playerScored) {
            player.score++;

            ball.reset_state();
        }

        if (scored.enemyScored) {
            enemy.score++;

            ball.reset_state();
        }

        // Move stuff
        player.move();
        ball.bounce();
        enemy.move(ball.pos.y);

        // Render stuff
        BeginDrawing();
            ClearBackground(WHITE);
            player.render();
            ball.render();
            enemy.render();

            player.render_score();
            enemy.render_score();
        EndDrawing();
    }

    return 0;
}