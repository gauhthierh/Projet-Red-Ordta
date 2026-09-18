// Command dev prepares the portable Windows toolchain, builds the game and
// optionally launches it. The command itself uses only the Go standard library.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	toolchainURL  = "https://github.com/skeeto/w64devkit/releases/download/v2.9.1/w64devkit-x64-2.9.1.7z.exe"
	toolchainHash = "9208c19755cd4964b7915b9afcf02c66d493a4c870c4b3e83f6c538d9c1237a5"
)

func main() {
	runAfterBuild := flag.Bool("run", false, "launch the game after a successful build")
	flag.Parse()

	root, err := projectRoot()
	must(err)
	if runtime.GOOS != "windows" {
		must(fmt.Errorf("this helper currently prepares the Windows build only"))
	}

	gcc, err := ensureToolchain(root)
	must(err)
	environment := buildEnvironment(root, gcc)
	must(run(root, environment, "go", "mod", "download"))

	binDir := filepath.Join(root, "bin")
	must(os.MkdirAll(binDir, 0o755))
	executable := filepath.Join(binDir, "red3d.exe")
	must(run(root, environment, "go", "build", "-o", executable, "./cmd/red3d"))
	must(copyRuntimeDLLs(root, binDir, environment))
	fmt.Printf("Jeu compile : %s\n", executable)

	if *runAfterBuild {
		must(run(binDir, append(environment, "PATH="+binDir+string(os.PathListSeparator)+os.Getenv("PATH")), executable))
	}
}

func projectRoot() (string, error) {
	directory, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(directory, "go.mod")); err == nil {
			return directory, nil
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			return "", fmt.Errorf("go.mod introuvable")
		}
		directory = parent
	}
}

func ensureToolchain(root string) (string, error) {
	toolsDir := filepath.Join(root, ".tools")
	gcc := filepath.Join(toolsDir, "w64devkit", "bin", "gcc.exe")
	if _, err := os.Stat(gcc); err == nil {
		return gcc, nil
	}
	if err := os.MkdirAll(toolsDir, 0o755); err != nil {
		return "", err
	}

	archive := filepath.Join(toolsDir, "w64devkit-x64-2.9.1.7z.exe")
	if !validHash(archive, toolchainHash) {
		fmt.Println("Telechargement de la chaine de compilation portable...")
		if err := download(toolchainURL, archive); err != nil {
			return "", err
		}
		if !validHash(archive, toolchainHash) {
			return "", fmt.Errorf("empreinte SHA-256 invalide pour %s", archive)
		}
	}

	extractor := exec.Command(archive, "-y", "-o"+toolsDir)
	extractor.Stdout = os.Stdout
	extractor.Stderr = os.Stderr
	if err := extractor.Run(); err != nil {
		return "", fmt.Errorf("extraction de w64devkit : %w", err)
	}
	if _, err := os.Stat(gcc); err != nil {
		return "", fmt.Errorf("gcc introuvable apres extraction : %w", err)
	}
	return gcc, nil
}

func download(url, destination string) error {
	client := &http.Client{Timeout: 15 * time.Minute}
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("telechargement : %s", response.Status)
	}

	temporary := destination + ".part"
	file, err := os.Create(temporary)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(file, response.Body)
	closeErr := file.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(temporary, destination)
}

func validHash(path, expected string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false
	}
	return strings.EqualFold(hex.EncodeToString(hash.Sum(nil)), expected)
}

func buildEnvironment(root, gcc string) []string {
	toolBin := filepath.Dir(gcc)
	environment := append([]string{}, os.Environ()...)
	environment = setEnv(environment, "GOARCH", "amd64")
	environment = setEnv(environment, "CGO_ENABLED", "1")
	environment = setEnv(environment, "CC", gcc)
	environment = setEnv(environment, "GOCACHE", filepath.Join(root, ".tools", "go-build-cache"))
	environment = setEnv(environment, "PATH", toolBin+string(os.PathListSeparator)+os.Getenv("PATH"))
	return environment
}

func copyRuntimeDLLs(root, binDir string, environment []string) error {
	command := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/g3n/engine")
	command.Dir = root
	command.Env = environment
	output, err := command.Output()
	if err != nil {
		return fmt.Errorf("localisation de G3N : %w", err)
	}
	sourceDir := filepath.Join(strings.TrimSpace(string(output)), "audio", "windows", "bin")
	for _, name := range []string{"OpenAL32.dll", "libogg.dll", "libvorbis.dll", "libvorbisfile.dll"} {
		if err := copyFile(filepath.Join(sourceDir, name), filepath.Join(binDir, name)); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(source, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	output, err := os.Create(destination)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(output, input)
	closeErr := output.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func run(directory string, environment []string, name string, arguments ...string) error {
	command := exec.Command(name, arguments...)
	command.Dir = directory
	command.Env = environment
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	fmt.Printf("> %s %s\n", name, strings.Join(arguments, " "))
	return command.Run()
}

func setEnv(environment []string, key, value string) []string {
	prefix := strings.ToUpper(key) + "="
	for index, entry := range environment {
		if strings.HasPrefix(strings.ToUpper(entry), prefix) {
			environment[index] = key + "=" + value
			return environment
		}
	}
	return append(environment, key+"="+value)
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur :", err)
		os.Exit(1)
	}
}
