SOURCES=*.go
LINKER_FLAGS=-X main.OwlVersion=`git describe --tags --dirty`
ifeq (${XDG_DATA_HOME},)
	XDG_DATA_HOME=${HOME}/.local/share
endif
MAN_DIR=${XDG_DATA_HOME}/man/man1
BIN_DIR=${HOME}/.local/bin

owl: ${SOURCES}
	go build -o owl -ldflags "${LINKER_FLAGS}" -- ${SOURCES}
owl.1: man.md
	pandoc -s --shift-heading-level-by=-1 --to=man man.md > owl.1
build: owl owl.1
test: *.go
	go test -v .
install: build
	install -m 0777 owl ${BIN_DIR}
	mkdir -p ${MAN_DIR}
	install owl.1 ${MAN_DIR}
uninstall:
	rm -I ${BIN_DIR}/owl ${MAN_DIR}/owl.1

.PHONY: build test install uninstall
