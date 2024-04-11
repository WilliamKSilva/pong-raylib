#include "player.h"
#include "raylib.h"

Player::Player(Vector2 _pos, Color _color) {
    pos = _pos;
    color = _color;
}

void Player::render() {
    DrawRectangle(pos.x, pos.y, width, height, color);
}