REF_DIRECTORY ?= /usr/local
REF_INCLUDE = $(REF_DIRECTORY)/include
REF_LIBRARY = $(REF_DIRECTORY)/lib

SRC = ./src
DST = ./bin/server
LOG = ./log

CXXFLAGS = -g -O0 -I$(REF_INCLUDE) -I./references/log4cplus

LDFLAGS =  -L$(REF_LIBRARY) -lmuduo_http -lmuduo_net -lmuduo_base -lpthread -lcurl -llog4cplus -lrt

$(DST): $(SRC)/*.cpp
	g++ -std=c++17 -O3 $(CXXFLAGS) -o $@ $^ $(LDFLAGS)
# log4cplus需要c++17支持

clean:
	rm -f $(SRC)/*.o $(DST) $(LOG)/*

.PHONY: all clean
