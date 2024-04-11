CC := g++ 
CFLAGS = -Wall

BUILD = ./build/

OBJS = $(BUILD)game.o

RAYLIB_PATH=/usr/local/lib/raylib/src/

asteroids:
	$(CC) -I$(RAYLIB_PATH) -g -o $(BUILD)pong game.cpp $(RAYLIB_PATH)libraylib.a -lGL -lm -lpthread -ldl -lrt -lX11