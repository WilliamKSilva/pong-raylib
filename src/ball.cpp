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

void Ball::bounce() {
    // Bounce to the left or to the right
    if (bounceHorizontal == LEFT) {
        pos.x -= speed;
    } else {
        pos.x += speed;

    }

    if (bounceVertical == NONE) {
        return;
    }

    // Bounce to the top or to the bottom
    if (bounceVertical == TOP) {
        pos.y -= speed;
    } else {
        pos.y += speed;
    }
}