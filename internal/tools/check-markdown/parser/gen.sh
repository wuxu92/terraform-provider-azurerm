#!/bin/bash
curdir=$(cd $(dirname $0) && pwd)
java -jar $HOME/dev/antlr/antlr-4.11.1-complete.jar -Dlanguage=Go -o ${curdir}/ ${curdir}/markdown.g4
