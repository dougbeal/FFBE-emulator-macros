#!/bin/sh
find . -path \*DONT_USE_THIS\* -prune -o -name LICENSE.txt -prune -o -type f -name \*.txt -print -exec ./memu2bs {} \;
