#include "game_object.h"
#ifndef PLAYER_H
#define PLAYER_H

#pragma once

class Player : public GameObject {
    public:
        const int width = 30;
        const int height = 100;
        const float speed = 10.0;

        Player(Vector2 _pos, Color _color);
        void render();
        void move();
};

#endif