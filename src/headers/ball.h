#include "raylib.h"
#include "game_object.h"
#ifndef BALL_H
#define BALL_H
class Ball : public GameObject {
    public:
        float radius;

        Ball(Vector2 _pos, float _radius, Color _color);
        void render();
};
#endif