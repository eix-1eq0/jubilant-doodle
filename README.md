Simple command line stopwatch, written in GO
Currently works properly only under Linux. Under BSD (and perhaps on macOS) the keyboard input reading doesn't work, so the program can't take lap time or exit without ctrl-C. The problem seems to be related with GO channels.
