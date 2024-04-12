#include <iostream>
#include "raylib.h"
#include "ball.h"
#include "game.h"
#include "string.h"

Ball::Ball(Vector2 _pos, float _radius, Color _color)
{
    pos = _pos;
    radius = _radius;
    color = _color;
}

void Ball::render()
{
    DrawCircle(pos.x, pos.y, radius, color);
}

void Ball::bounce()
{
    // Bounce to the left or to the right
    if (bounceHorizontal == LEFT)
    {
        pos.x -= speed;
    }
    else
    {
        pos.x += speed;
    }

    if (bounceVertical == NONE)
    {
        return;
    }

    // Bounce to the top or to the bottom
    if (bounceVertical == TOP)
    {
        pos.y -= speed;
    }
    else
    {
        pos.y += speed;
    }
}

void Ball::check_collision(Rectangle rect, const char* object_name)
{
    bool collided = CheckCollisionCircleRec(pos, radius, rect);

    if (collided)
    {
        int collision_point_y = pos.y - rect.y;

        if (collision_point_y <= collision_point_top)
        {
            bounceVertical = TOP;
        }

        if (collision_point_y >= collision_point_bottom)
        {
            bounceVertical = BOTTOM;
        }


        if (strcmp(object_name, "Player") == 0) {
            bounceHorizontal = RIGHT;
        } else {
            bounceHorizontal = LEFT;
        }
    }
}

void Ball::check_out_of_bounds()
{
    if (pos.y < 0)
    {
        bounceVertical = BOTTOM;
    }

    // TODO: add Window class height instead of magic number
    if (pos.y > 1080)
    {
        bounceVertical = TOP;
    }
}

Scored Ball::check_scored()
{
    // If ball position is off the window to the right, the player has scored
    if (pos.x >= 1920)
    {
        return {playerScored : true, enemyScored : false};
    }

    // If ball position is off the window to the left, the enemy has scored
    if (pos.x <= 0)
    {
        return {playerScored : false, enemyScored : true};
    }

    return {playerScored : false, enemyScored : false};
}

void Ball::reset_state() {
    pos.x = Window::width / 2;
    pos.y = Window::height / 2;

    bounceHorizontal = LEFT;
    bounceVertical = NONE;
}