package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"os"
	"strconv"

	"StegoConnectivity/internal/crypto"
	"StegoConnectivity/internal/stego"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App                         { return &App{} }
func (a *App) startup(ctx context.Context) { a.ctx = ctx }

type GenResult struct {
	Image   string `json:"image"`
	PrivKey string `json:"privKey"`
	Log     string `json:"log"`
}

type ExtResult struct {
	Payload   string `json:"payload"`
	PublicKey string `json:"publicKey"`
	Log       string `json:"log"`
}

func (a *App) Generate(mode string, text string, keyType string, keySizeOrCurve string, conn string, feature string, encrypt bool, pass string) GenResult {
	var payload []byte
	var pubKey, privKey string
	var err error
	meta := make(map[string]string)

	if mode == "Asymmetric" {
		if keyType == "RSA" {
			bits, _ := strconv.Atoi(keySizeOrCurve)
			privKey, pubKey, err = crypto.GenerateRSA(bits)
		} else if keyType == "DH" {
			privKey, pubKey, err = crypto.GenerateDHGroup(keySizeOrCurve)
		} else {
			privKey, pubKey, err = crypto.GenerateECC(keySizeOrCurve)
		}
		if err != nil {
			return GenResult{Log: "Key Error: " + err.Error()}
		}
		payload = []byte(privKey)
		meta["PublicKey"] = pubKey
	} else {
		payload = []byte(text)
	}

	if encrypt {
		cipher, salt, nonce, e := crypto.Encrypt(payload, pass)
		if e != nil {
			return GenResult{Log: "Encrypt Error"}
		}
		payload = cipher
		meta["Salt"] = hex.EncodeToString(salt)
		meta["Nonce"] = hex.EncodeToString(nonce)
		meta["Algo"] = "AES-GCM"
	}

	// 1. Generate the Blob (stego/process.go)
	imgBytes, err := stego.Embed(payload, conn, feature)
	if err != nil {
		return GenResult{Log: "Embed Error"}
	}

	// 2. Inject Metadata (stego/metadata.go)
	meta["Feature"] = feature
	imgBytes = stego.InjectMeta(imgBytes, meta)

	return GenResult{
		Image:   base64.StdEncoding.EncodeToString(imgBytes),
		PrivKey: privKey,
		Log:     "Success: Image Generated (" + feature + ")",
	}
}

func (a *App) AnalyzeImage(imgB64 string, connMode string, featureMode string) stego.AnalysisResult {
	b, _ := base64.StdEncoding.DecodeString(imgB64)
	return stego.Analyze(b, connMode, featureMode)
}

func (a *App) AnalyzeText(text string, connMode string, featureMode string) stego.AnalysisResult {
	imgBytes, _ := stego.Embed([]byte(text), connMode, featureMode)
	return stego.Analyze(imgBytes, connMode, featureMode)
}

func (a *App) Extract(imgB64 string, conn string, feature string, decrypt bool, pass string) ExtResult {
	data, err := base64.StdEncoding.DecodeString(imgB64)
	if err != nil {
		return ExtResult{Log: "Base64 Error"}
	}

	// 1. Extract Logic (stego/process.go)
	raw, _ := stego.Extract(data, conn, feature)

	// 2. Read Metadata (stego/metadata.go)
	meta := stego.ReadMeta(data)

	finalData := raw
	statusLog := "Success: Extracted"

	if decrypt {
		salt, _ := hex.DecodeString(meta["Salt"])
		nonce, _ := hex.DecodeString(meta["Nonce"])
		algo := meta["Algo"]
		if len(salt) == 0 {
			statusLog = "Metadata Missing"
		} else {
			dec, err := crypto.Decrypt(raw, salt, nonce, pass)
			if err != nil {
				statusLog = "Decryption Failed (" + algo + ")"
			} else {
				finalData = dec
				statusLog = "Success: Decrypted (" + algo + ")"
			}
		}
	}

	return ExtResult{
		Payload:   string(finalData),
		PublicKey: meta["PublicKey"],
		Log:       statusLog,
	}
}

func (a *App) SaveFile(name, dataB64 string) {
	b, _ := base64.StdEncoding.DecodeString(dataB64)
	path, _ := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{DefaultFilename: name})
	if path != "" {
		os.WriteFile(path, b, 0644)
	}
}
