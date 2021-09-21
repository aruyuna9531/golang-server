MUDUO_DIRECTORY ?= /usr/local
MUDUO_INCLUDE = $(MUDUO_DIRECTORY)/include
MUDUO_LIBRARY = $(MUDUO_DIRECTORY)/lib
SRC = ./src
DST = ./bin/server

CXXFLAGS = -g -O0 -I$(MUDUO_INCLUDE)

LDFLAGS =  -L$(MUDUO_LIBRARY) -lmuduo_http -lmuduo_net -lmuduo_base -lpthread -lrt

$(DST): $(SRC)/*.cpp
	g++ $(CXXFLAGS) -o $@ $^ $(LDFLAGS)

clean:
	rm -f $(SRC)/*.o $(DST)

.PHONY: all clean
