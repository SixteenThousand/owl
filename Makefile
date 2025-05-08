SOURCES=*.go
LINKER_FLAGS=-X main.OwlVersion=`git describe --tags --dirty`
owl: ${SOURCES}
	go build -o owl -ldflags "${LINKER_FLAGS}" -- ${SOURCES}
owl.1: man.md
	pandoc -s --shift-heading-level-by=-1 --to=man man.md > owl.1
build: owl owl.1
.PHONY: build
