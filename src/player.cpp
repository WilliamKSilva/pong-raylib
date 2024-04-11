#include "player.h"
#include "raylib.h"

Player::Player(Vector2 _pos, Color _color, const char* _name) {
    pos = _pos;
    color = _color;
    name = _name;
}

void Player::render() {
    DrawRectangle(pos.x, pos.y, width, height, color);
}

void Player::move() {
    if (IsKeyDown(KEY_W)) {
        pos.y -= speed;
    }

    if (IsKeyDown(KEY_S)) {
        pos.y += speed;
    }
}