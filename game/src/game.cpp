#include "raylib.h"
#include <emscripten/emscripten.h>

// Application level imports
#include "ball.h"
#include "enemy.h"
#include "game.h"
#include "player.h"

// To run on the browser
void UpdateDrawFrame(void);

Player player = Player({.x = 100, .y = 30}, BLACK, "Player");
Enemy enemy = Enemy({.x = Window::width - 100, .y = 30}, BLACK, "Enemy");
Ball ball =
    Ball({.x = (float)Window::width / 2.0, .y = (float)Window::height / 2.0},
         20.0, BLACK);

int main() {
  Window::init();
  GameMode gameMode = OFFLINE;

  emscripten_set_main_loop(UpdateDrawFrame, 60, 1);

  return 0;
}

// Update and draw game frame
void UpdateDrawFrame(void) {

  ball.check_collision({.x = player.pos.x,
                        .y = player.pos.y,
                        .width = (float)player.width,
                        .height = (float)player.height},
                       player.name);

  ball.check_collision({.x = enemy.pos.x,
                        .y = enemy.pos.y,
                        .width = (float)enemy.width,
                        .height = (float)enemy.height},
                       enemy.name);

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

  emscripten_log(0, "%f", player.pos.y);

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
