#include "raylib.h"
#include "enemy.h"

Enemy::Enemy(Vector2 _pos, Color _color, const char* _name) {
    pos = _pos;
    color = _color;
    name = _name;
}

void Enemy::render() {
    DrawRectangle(pos.x, pos.y, width, height, color);
}

void Enemy::move(float ball_pos_y) {
   if (pos.y <= ball_pos_y) {
    pos.y += speed;
   } 

   if (pos.y >= ball_pos_y) {
    pos.y -= speed;
   }
}