#include "game_object.h"
#ifndef PLAYER_H
#define PLAYER_H

#pragma once

class Player : public GameObject {
    public:
        const int width = 30;
        const int height = 100;
        const float speed = 10.0;
        int score = 0;

        Player(Vector2 _pos, Color _color, const char* _name);
        void render();
        void move();
        void render_score();
};

#endif