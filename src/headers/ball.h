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

typedef struct {
    bool playerScored;
    bool enemyScored;
} Scored;

class Ball : public GameObject {
    private:
        const int collision_point_top = 20.0; 
        const int collision_point_bottom = 70.0; 

    public:
        float radius;
        float speed = 12.0;
        // Always starting going to the player side
        BounceHorizontal bounceHorizontal = LEFT; 
        BounceVertical bounceVertical = NONE;

        Ball(Vector2 _pos, float _radius, Color _color);
        void render();
        void bounce();
        // Player and Enemy are rectangle shapes
        void check_collision(Rectangle rect, const char* object_name);
        void check_out_of_bounds();
        Scored check_scored();
};
#endif