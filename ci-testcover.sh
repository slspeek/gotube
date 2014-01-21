#!/bin/bash
PACKAGE=github.com/slspeek/$1
SHORTNAME=$(basename $PACKAGE)
outfile=test-temp.txt
go test -v $PACKAGE| tee $outfile 
go test -v -bench='.*' -benchmem $PACKAGE
go2xunit -fail -input $outfile -output ${SHORTNAME}-tests.xml
gocov test $PACKAGE | gocov-xml > ${SHORTNAME}-coverage.xml
gocov test $PACKAGE | gocov-html > ${SHORTNAME}-coverage.html

