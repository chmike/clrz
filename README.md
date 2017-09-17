# Colorize

Package clrz colorize text like code for display in HTML or other encoding. 

Note: this is a Work In Progress (WIP).

Go is still lacking a package equivalent to Python's Pygments library.
For fun and programming with Go experience, I spent some time in my 
holidays to study how Pygments and equivalent libraries work, and 
implemented this experimental package. 

Most package of this type use regex expression to identify token to
colorize. But regexp are also known to be slow lexers. 

This package provides thus an abstraction allowing to provide efficient
lexers and regex based lexers. It is expected that we could port, in a 
first step, regex lexers of pygments to Colorize, and later, in a 
second step, optimize them as needed. 

Progress on this package was suspended by the end of my holidays. 
Development may resume later, when I have time to work on it again.

Feedback is welcome. 