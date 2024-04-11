#include "game_object.h"
#ifndef ENEMY_H 
#define ENEMY_H 

#pragma once

class Enemy : public GameObject {
    public:
        const int width = 30;
        const int height = 100;
        const float speed = 10.0;
        int score = 0;

        Enemy(Vector2 _pos, Color _color, const char* _name);
        void render();
        void move(float player_pos_y);
};

#endif