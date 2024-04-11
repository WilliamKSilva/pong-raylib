#include <stdio.h>
#include "raylib.h"
#include "enemy.h"
#include "game.h"

Enemy::Enemy(Vector2 _pos, Color _color, const char *_name)
{
    pos = _pos;
    color = _color;
    name = _name;
}

void Enemy::render()
{
    DrawRectangle(pos.x, pos.y, width, height, color);
}

void Enemy::move(float ball_pos_y)
{
    if (pos.y <= ball_pos_y && pos.y < Window::height - height)
    {
        pos.y += speed;
    }

    if (pos.y >= ball_pos_y && pos.y > 0)
    {
        pos.y -= speed;
    }
}

void Enemy::render_score() {
    char buf[256];
    snprintf(buf, sizeof(buf), "Score: %d", score);
    DrawText(buf, Window::width - 200, 10, 31, BLACK);
}