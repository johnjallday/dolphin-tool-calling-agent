# Name of build output folder
BUILD_DIR   ?= build

# Name of the binary
BINARY      ?= gui

# Path to your cmd package
CMD_PATH    ?= ./cmd/gui

.PHONY: all clean gui

# Default: clean then build
all: clean gui

# remove everything in BUILD_DIR (but keep the folder)
clean:
	@mkdir -p $(BUILD_DIR)
	@rm -rf $(BUILD_DIR)/*

# compile your cmd/gui into BUILD_DIR/gui
gui:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) $(CMD_PATH)
