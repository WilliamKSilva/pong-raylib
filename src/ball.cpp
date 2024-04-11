#include "raylib.h"
#include "ball.h"

Ball::Ball(Vector2 _pos, float _radius, Color _color) {
    pos = _pos;
    radius = _radius;
    color = _color;
}

void Ball::render() {
    DrawCircle(pos.x, pos.y, radius, color);
}