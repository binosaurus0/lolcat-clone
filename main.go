package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	// "strconv"
	"strings"
	"time"
)

// ColorMode represents different coloring algorithms
type ColorMode int

const (
	Rainbow ColorMode = iota
	Fire
	Ocean
	Matrix
	Pastel
)

// Config holds the application configuration
type Config struct {
	mode       ColorMode
	speed      float64
	spread     float64
	frequency  float64
	animate    bool
	delay      time.Duration
	force      bool
	help       bool
}

// rgb generates RGB values based on position and mode
func rgb(i int, config *Config) (int, int, int) {
	f := config.frequency
	s := config.speed
	offset := float64(i) * config.spread

	switch config.mode {
	case Rainbow:
		return int(math.Sin(f*(offset+s)+0)*127 + 128),
			int(math.Sin(f*(offset+s)+2*math.Pi/3)*127 + 128),
			int(math.Sin(f*(offset+s)+4*math.Pi/3)*127 + 128)
	case Fire:
		r := int(math.Sin(f*(offset+s))*127 + 128)
		g := int(math.Sin(f*(offset+s)+math.Pi/2)*64 + 64)
		b := int(math.Sin(f*(offset+s)+math.Pi)*32 + 32)
		return r, g, b
	case Ocean:
		r := int(math.Sin(f*(offset+s)+math.Pi)*64 + 64)
		g := int(math.Sin(f*(offset+s)+math.Pi/2)*127 + 128)
		b := int(math.Sin(f*(offset+s))*127 + 128)
		return r, g, b
	case Matrix:
		g := int(math.Sin(f*(offset+s))*127 + 128)
		return 0, g, int(float64(g) * 0.3)
	case Pastel:
		r := int(math.Sin(f*(offset+s)+0)*64 + 192)
		g := int(math.Sin(f*(offset+s)+2*math.Pi/3)*64 + 192)
		b := int(math.Sin(f*(offset+s)+4*math.Pi/3)*64 + 192)
		return r, g, b
	default:
		return 255, 255, 255
	}
}

// printColored prints a rune with the specified RGB color
func printColored(r rune, red, green, blue int) {
	fmt.Printf("\033[38;2;%d;%d;%dm%c\033[0m", red, green, blue, r)
}

// printLine processes and prints a line of text
func printLine(line string, config *Config, lineOffset int) {
	runes := []rune(line)
	for i, r := range runes {
		red, green, blue := rgb(i+lineOffset, config)
		printColored(r, red, green, blue)
	}
	fmt.Println()
}

// animateLine prints a line with animation effect
func animateLine(line string, config *Config, lineOffset int) {
	if !config.animate {
		printLine(line, config, lineOffset)
		return
	}

	runes := []rune(line)
	for frame := 0; frame < 20; frame++ {
		fmt.Print("\r")
		tempConfig := *config
		tempConfig.speed = config.speed + float64(frame)*0.1
		
		for i, r := range runes {
			red, green, blue := rgb(i+lineOffset, &tempConfig)
			printColored(r, red, green, blue)
		}
		
		time.Sleep(config.delay)
	}
	fmt.Println()
}

// processInput reads from input source and applies coloring
func processInput(input io.Reader, config *Config) error {
	scanner := bufio.NewScanner(input)
	lineOffset := 0
	
	for scanner.Scan() {
		line := scanner.Text()
		if config.animate {
			animateLine(line, config, lineOffset)
		} else {
			printLine(line, config, lineOffset)
		}
		lineOffset += len([]rune(line))
	}
	
	return scanner.Err()
}

// parseColorMode converts string to ColorMode
func parseColorMode(mode string) (ColorMode, error) {
	switch strings.ToLower(mode) {
	case "rainbow", "r":
		return Rainbow, nil
	case "fire", "f":
		return Fire, nil
	case "ocean", "o":
		return Ocean, nil
	case "matrix", "m":
		return Matrix, nil
	case "pastel", "p":
		return Pastel, nil
	default:
		return Rainbow, fmt.Errorf("unknown color mode: %s", mode)
	}
}

// printUsage displays help information
func printUsage() {
	fmt.Println("gololcat - Colorize text with rainbow and other effects")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("    gololcat [OPTIONS] [FILE...]")
	fmt.Println("    command | gololcat [OPTIONS]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("    -m, --mode MODE      Color mode: rainbow(r), fire(f), ocean(o), matrix(m), pastel(p) [default: rainbow]")
	fmt.Println("    -s, --speed FLOAT    Animation speed [default: 0.0]")
	fmt.Println("    -p, --spread FLOAT   Color spread factor [default: 1.0]")
	fmt.Println("    -f, --frequency FLOAT Frequency of color changes [default: 0.1]")
	fmt.Println("    -a, --animate        Enable animation effect")
	fmt.Println("    -d, --delay DURATION Animation delay [default: 100ms]")
	fmt.Println("    --force              Force color output even when not to a terminal")
	fmt.Println("    -h, --help           Show this help")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("    echo 'Hello World' | gololcat")
	fmt.Println("    gololcat --mode fire --animate file.txt")
	fmt.Println("    fortune | gololcat -m ocean -s 2.0")
	fmt.Println("    ls -la | gololcat --mode matrix")
}

// parseFlags handles command line argument parsing
func parseFlags() (*Config, []string, error) {
	config := &Config{
		mode:      Rainbow,
		speed:     0.0,
		spread:    1.0,
		frequency: 0.1,
		animate:   false,
		delay:     100 * time.Millisecond,
		force:     false,
		help:      false,
	}

	var modeStr string
	var delayStr string

	flag.StringVar(&modeStr, "mode", "rainbow", "Color mode")
	flag.StringVar(&modeStr, "m", "rainbow", "Color mode (short)")
	flag.Float64Var(&config.speed, "speed", 0.0, "Animation speed")
	flag.Float64Var(&config.speed, "s", 0.0, "Animation speed (short)")
	flag.Float64Var(&config.spread, "spread", 1.0, "Color spread factor")
	flag.Float64Var(&config.spread, "p", 1.0, "Color spread factor (short)")
	flag.Float64Var(&config.frequency, "frequency", 0.1, "Frequency of color changes")
	flag.Float64Var(&config.frequency, "f", 0.1, "Frequency of color changes (short)")
	flag.BoolVar(&config.animate, "animate", false, "Enable animation")
	flag.BoolVar(&config.animate, "a", false, "Enable animation (short)")
	flag.StringVar(&delayStr, "delay", "100ms", "Animation delay")
	flag.StringVar(&delayStr, "d", "100ms", "Animation delay (short)")
	flag.BoolVar(&config.force, "force", false, "Force color output")
	flag.BoolVar(&config.help, "help", false, "Show help")
	flag.BoolVar(&config.help, "h", false, "Show help (short)")

	flag.Parse()

	if config.help {
		return config, nil, nil
	}

	// Parse color mode
	mode, err := parseColorMode(modeStr)
	if err != nil {
		return nil, nil, err
	}
	config.mode = mode

	// Parse delay
	delay, err := time.ParseDuration(delayStr)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid delay format: %s", delayStr)
	}
	config.delay = delay

	return config, flag.Args(), nil
}

// isTerminal checks if we're outputting to a terminal
func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func main() {
	config, files, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if config.help {
		printUsage()
		return
	}

	// Check if output is to terminal (unless forced)
	if !config.force && !isTerminal() {
		fmt.Fprintf(os.Stderr, "Output is not a terminal. Use --force to override.\n")
		os.Exit(1)
	}

	// If no files specified, read from stdin
	if len(files) == 0 {
		// Check if stdin is a terminal (no piped input)
		stdinInfo, _ := os.Stdin.Stat()
		if (stdinInfo.Mode() & os.ModeCharDevice) != 0 {
			fmt.Fprintf(os.Stderr, "Error: No input provided. Use pipes or specify files.\n")
			fmt.Fprintf(os.Stderr, "Usage: echo 'text' | gololcat or gololcat file.txt\n")
			os.Exit(1)
		}
		
		err := processInput(os.Stdin, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Process each file
	for _, filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening %s: %v\n", filename, err)
			continue
		}

		err = processInput(file, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", filename, err)
		}

		file.Close()
	}
}