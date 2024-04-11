#include "raylib.h"
#include "game_object.h"
#ifndef BALL_H
#define BALL_H

enum BounceHorizontal {
    LEFT,
    RIGHT
};

enum BounceVertical {
    NONE,
    TOP,
    BOTTOM
};

class Ball : public GameObject {
    public:
        float radius;
        float speed = 12.0;
        // Always starting going to the player side
        BounceHorizontal bounceHorizontal = LEFT; 
        BounceVertical bounceVertical = NONE;

        Ball(Vector2 _pos, float _radius, Color _color);
        void render();
        void bounce();
};
#endif