libs=-lmuduo_net -lmuduo_base -lpthread

obj=echo.o echo_server.o

server : $(obj)
	g++ -o server $(obj) $(libs)

echo.o : echo.cpp
	g++ -c echo.cpp

echo_server.o : echo_server.cpp
	g++ -c echo_server.cpp

clean:
	rm -rf server* $(obj)