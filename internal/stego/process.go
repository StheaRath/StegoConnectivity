package stego

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"unicode"
)

type AnalysisResult struct {
	RealData string `json:"realData"`
	ConnMap  string `json:"connMap"`
	BitGrid  string `json:"bitGrid"`
	DataPath string `json:"dataPath"`
}

type Point struct{ X, Y int }

// --- SPATIAL FEATURE DEFINITIONS ---

func isCoreFeature(x, y int) bool { return (x%2 == 0) && (y%2 == 0) }
func isEdgeFeature(x, y int) bool { return !((x%2 == 0) && (y%2 == 0)) }

func isRadialFeature(x, y, w, h int) bool {
	cx, cy := float64(w)/2, float64(h)/2
	dx, dy := float64(x)-cx, float64(y)-cy
	dist := math.Sqrt(dx*dx + dy*dy)
	return int(dist)%20 < 10
}

func isSkeletonFeature(x, y, w, h int) bool {
	cx, cy := w/2, h/2
	onAxis := (x == cx) || (y == cy)
	onDiagonal := (math.Abs(float64(x-cx)) == math.Abs(float64(y-cy)))
	return onAxis || onDiagonal
}

func isTextureFeature(x, y int) bool { return (x+y)%3 == 0 }

// Master Validator
func isValidPixel(x, y, w, h int, featureMode string) bool {
	switch featureMode {
	case "Core":
		return isCoreFeature(x, y)
	case "Edge":
		return isEdgeFeature(x, y)
	case "Radial":
		return isRadialFeature(x, y, w, h)
	case "Skeleton":
		return isSkeletonFeature(x, y, w, h)
	case "Texture":
		return isTextureFeature(x, y)
	case "Full":
		return true
	default:
		return isCoreFeature(x, y)
	}
}

// --- CONNECTIVITY TRAVERSAL ---
func getNeighbors(p Point, mode string) []Point {
	n4 := []Point{{0, -1}, {-1, 0}, {1, 0}, {0, 1}}
	diag := []Point{{-1, -1}, {1, -1}, {-1, 1}, {1, 1}}

	if mode == "4" {
		var res []Point
		for _, d := range n4 {
			res = append(res, Point{p.X + d.X, p.Y + d.Y})
		}
		return res
	}
	if mode == "m" {
		var res []Point
		for _, d := range n4 {
			res = append(res, Point{p.X + d.X, p.Y + d.Y})
		}
		for _, d := range diag {
			res = append(res, Point{p.X + d.X, p.Y + d.Y})
		}
		return res
	}
	var res []Point
	all8 := append(n4, diag...)
	for _, d := range all8 {
		res = append(res, Point{p.X + d.X, p.Y + d.Y})
	}
	return res
}

func Embed(payload []byte, connMode string, featureMode string) ([]byte, error) {
	data := make([]byte, 4+len(payload))
	binary.BigEndian.PutUint32(data, uint32(len(payload)))
	copy(data[4:], payload)

	totalBits := len(data) * 8

	// --- DYNAMIC DENSITY CALCULATION ---
	// Adjust canvas size based on how "wasteful" the feature is.
	var densityFactor float64

	switch featureMode {
	case "Full":
		// Uses every pixel -> Tight fit
		densityFactor = 1.1
	case "Core", "Edge", "Texture":
		// Uses ~50% of pixels (checkerboard) -> Need 2x space
		densityFactor = 2.1
	case "Radial":
		// Skips bands -> Needs ~2.5x space
		densityFactor = 2.5
	case "Skeleton":
		// Very sparse (only axes/diagonals) -> Needs lots of space
		densityFactor = 7.0
	default:
		densityFactor = 2.5
	}

	// Calculate Dimension (Square Root of Area needed)
	area := float64(totalBits) * densityFactor
	dim := int(math.Ceil(math.Sqrt(area)))

	// Add a small safety padding (10px) to prevent edge clipping on irregular shapes
	dim += 10

	// Minimum bounds (reduced to avoid huge empty images for small text)
	if dim < 64 {
		dim = 64
	}

	img := image.NewRGBA(image.Rect(0, 0, dim, dim))
	for i := 0; i < len(img.Pix); i += 4 {
		img.Pix[i+3] = 255
	} // Black Background

	queue := []Point{{dim / 2, dim / 2}}
	visited := make(map[Point]bool)
	visited[queue[0]] = true

	bitIdx := 0

	// BFS Loop
	for len(queue) > 0 && bitIdx < len(data)*8 {
		curr := queue[0]
		queue = queue[1:]

		// Check Spatial Feature
		if isValidPixel(curr.X, curr.Y, dim, dim, featureMode) {
			bytePos := bitIdx / 8
			bitPos := 7 - (bitIdx % 8)
			if (data[bytePos]>>bitPos)&1 == 1 {
				img.Set(curr.X, curr.Y, color.White)
			}
			bitIdx++
		}

		neighbors := getNeighbors(curr, connMode)
		for _, n := range neighbors {
			if n.X >= 0 && n.X < dim && n.Y >= 0 && n.Y < dim {
				if !visited[n] {
					visited[n] = true
					queue = append(queue, n)
				}
			}
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes(), nil
}

func Extract(imgData []byte, connMode string, featureMode string) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	queue := []Point{{w / 2, h / 2}}
	visited := make(map[Point]bool)
	visited[queue[0]] = true

	var bits []int
	maxBits := 800000

	for len(queue) > 0 && len(bits) < maxBits {
		curr := queue[0]
		queue = queue[1:]

		if isValidPixel(curr.X, curr.Y, w, h, featureMode) {
			r, _, _, _ := img.At(curr.X, curr.Y).RGBA()
			val := 0
			if (r >> 8) > 128 {
				val = 1
			}
			bits = append(bits, val)
		}

		neighbors := getNeighbors(curr, connMode)
		for _, n := range neighbors {
			if n.X >= 0 && n.X < w && n.Y >= 0 && n.Y < h {
				if !visited[n] {
					visited[n] = true
					queue = append(queue, n)
				}
			}
		}
	}

	byteLen := len(bits) / 8
	if byteLen == 0 {
		return nil, fmt.Errorf("no data")
	}
	raw := make([]byte, byteLen)
	for i := 0; i < byteLen; i++ {
		for b := 0; b < 8; b++ {
			if bits[i*8+b] == 1 {
				raw[i] |= 1 << (7 - b)
			}
		}
	}

	if len(raw) < 4 {
		return raw, nil
	}
	lenVal := binary.BigEndian.Uint32(raw[:4])
	if lenVal > uint32(len(raw)) || lenVal == 0 {
		return raw, nil
	}
	return raw[4 : 4+lenVal], nil
}

func Analyze(imgData []byte, connMode string, featureMode string) AnalysisResult {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return AnalysisResult{}
	}

	// 1. EXTRACT
	extractedBytes, _ := Extract(imgData, connMode, featureMode)

	safeData := ""
	for _, b := range extractedBytes {
		if unicode.IsPrint(rune(b)) {
			safeData += string(b)
		} else {
			safeData += "."
		}
	}
	if len(safeData) > 1000 {
		safeData = safeData[:1000] + "..."
	}

	// 2. MAP
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	queue := []Point{{w / 2, h / 2}}
	visited := make(map[Point]bool)
	visited[queue[0]] = true
	charMap := make(map[Point]string)

	bitIdx := 0
	headerBits := 32

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if isValidPixel(curr.X, curr.Y, w, h, featureMode) {
			if bitIdx >= headerBits {
				byteIdx := (bitIdx - headerBits) / 8
				if byteIdx < len(extractedBytes) {
					b := extractedBytes[byteIdx]
					char := "."
					if unicode.IsPrint(rune(b)) {
						char = string(b)
					}
					charMap[curr] = char
				}
			}
			bitIdx++
		}

		neighbors := getNeighbors(curr, connMode)
		for _, n := range neighbors {
			if n.X >= 0 && n.X < w && n.Y >= 0 && n.Y < h {
				if !visited[n] {
					visited[n] = true
					queue = append(queue, n)
				}
			}
		}
	}

	// 3. GENERATE GRIDS
	var sbConn, sbBit, sbPath bytes.Buffer
	limit := 256

	for y := 0; y < limit && y < h; y++ {
		for x := 0; x < limit && x < w; x++ {
			r, _, _, _ := img.At(x, y).RGBA()
			isWhite := (r >> 8) > 128
			p := Point{x, y}

			if isWhite {
				sbBit.WriteString("1 ")
			} else {
				sbBit.WriteString("0 ")
			}

			// For Connectivity Map
			if isWhite {
				if isCoreFeature(x, y) {
					sbConn.WriteString("C ")
				} else {
					sbConn.WriteString("E ")
				}
			} else {
				sbConn.WriteString(". ")
			}

			// Data Path
			if isWhite {
				if val, ok := charMap[p]; ok {
					sbPath.WriteString(val + " ")
				} else {
					sbPath.WriteString("? ")
				}
			} else {
				sbPath.WriteString(". ")
			}
		}
		sbBit.WriteString("\n")
		sbConn.WriteString("\n")
		sbPath.WriteString("\n")
	}

	return AnalysisResult{
		RealData: safeData,
		ConnMap:  sbConn.String(),
		BitGrid:  sbBit.String(),
		DataPath: sbPath.String(),
	}
}
