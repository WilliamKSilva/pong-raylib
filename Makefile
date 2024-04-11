CC := g++ 
CFLAGS = -Wall

BUILD = ./build/

OBJS = $(BUILD)game.o \
		$(BUILD)player.o

RAYLIB_PATH=/usr/local/lib/raylib/src/

HEADER_FILES = ./src/headers/
SRC_FILES = ./src/game.cpp \
			./src/player.cpp \
			./src/ball.cpp

pong:
	$(CC) -I$(RAYLIB_PATH) -I$(HEADER_FILES) -o $(BUILD)pong $(SRC_FILES) $(RAYLIB_PATH)libraylib.a -lGL -lm -lpthread -ldl -lrt -lX11