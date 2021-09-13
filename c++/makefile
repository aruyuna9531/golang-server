MUDUO_DIRECTORY ?= /usr/local
MUDUO_INCLUDE = $(MUDUO_DIRECTORY)/include
MUDUO_LIBRARY = $(MUDUO_DIRECTORY)/lib
SRC = .

CXXFLAGS = -g -O0 -I$(MUDUO_INCLUDE)

LDFLAGS =  -L$(MUDUO_LIBRARY) -lmuduo_http -lmuduo_net -lmuduo_base -lpthread -lrt

server: $(SRC)/*.cpp
	g++ $(CXXFLAGS) -o $@ $^ $(LDFLAGS)

clean:
	rm -f echo core *.o

.PHONY: all clean