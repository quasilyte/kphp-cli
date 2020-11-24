package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cespare/subcmd"
)

var kphpRoot = findKPHPRoot()

func main() {
	cmds := []subcmd.Command{
		{
			Name:        "run",
			Description: "compile and run KPHP script",
			Do:          cmdRun,
		},
		{
			Name:        "env",
			Description: "print KPHP environment information",
			Do:          cmdEnv,
		},
		{
			Name:        "version",
			Description: "print KPHP version information",
			Do:          cmdVersion,
		},
	}

	subcmd.Run(cmds)
}

func findKPHPRoot() string {
	// TODO(quasilyte): try to infer it automatically?
	// ~/kphp may be a good guess
	return os.Getenv("KPHP_ROOT")
}

func requireKPHPRoot() {
	if kphpRoot == "" {
		log.Fatal("$KPHP_ROOT is unset")
	}
}

func cmdEnv(args []string) {
	kphpVars := []string{
		"KPHP_ROOT",
		"KPHP_TESTS_POLYFILLS_REPO",
	}

	for _, name := range kphpVars {
		v := os.Getenv(name)
		fmt.Printf("%s=%q\n", name, v)
	}
}

func cmdRun(args []string) {
	log.SetFlags(0)

	fs := flag.NewFlagSet("kphp run", flag.ExitOnError)
	composerRoot := fs.String("composer-root", "",
		"A folder that contains the root composer.json file")
	serverMode := fs.Bool("server", false,
		"Whether to compile a script with --mode=server")
	fs.Parse(args)

	requireKPHPRoot()

	// TODO(quasilyte): support binary KPHP installations?
	kphp2cpp := filepath.Join(kphpRoot, "objs", "bin", "kphp2cpp")
	if !fileExists(kphp2cpp) {
		log.Print("$KPHP_ROOT/objs/bin/kphp2cpp does not exist; have you compiled the KPHP?")
		return
	}

	if *composerRoot == "" {
		// Maybe use current dir as a composer root?
		wd, err := os.Getwd()
		if err != nil {
			log.Panicf("getwd: %v", err)
		}
		if fileExists(filepath.Join(wd, "composer.json")) {
			*composerRoot = wd
		}
	}

	tempDir, err := ioutil.TempDir("", "kphp-run")
	if err != nil {
		log.Panicf("create temp build dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			log.Printf("remove temp build dir: %v", err)
		}
	}()

	kphpArgs := []string{
		"--destination-directory", tempDir,
	}
	if !*serverMode {
		kphpArgs = append(kphpArgs, "--mode", "cli")
	}
	if *composerRoot != "" {
		kphpArgs = append(kphpArgs, "--composer-root", *composerRoot)
	}
	kphpArgs = append(kphpArgs, fs.Args()...)

	buildOut, err := exec.Command(kphp2cpp, kphpArgs...).CombinedOutput()
	if err != nil {
		log.Printf("kphp2cpp error: %v: %s", err, buildOut)
		return
	}

	if *serverMode {
		scriptOut, err := exec.Command(filepath.Join(tempDir, "server"), "-o").CombinedOutput()
		if err != nil {
			log.Printf("server error: %v: %s", err, scriptOut)
		}
		fmt.Print(string(scriptOut))
	} else {
		scriptOut, err := exec.Command(filepath.Join(tempDir, "cli")).CombinedOutput()
		if err != nil {
			log.Printf("script error: %v: %s", err, scriptOut)
		}
		fmt.Print(string(scriptOut))
	}
}

func cmdVersion(args []string) {
	requireKPHPRoot()

	// TODO(quasilyte): support binary KPHP installations?
	kphp2cpp := filepath.Join(kphpRoot, "objs", "bin", "kphp2cpp")
	if !fileExists(kphp2cpp) {
		log.Print("$KPHP_ROOT/objs/bin/kphp2cpp does not exist; have you compiled the KPHP?")
		return
	}

	out, err := exec.Command(kphp2cpp, "--version").CombinedOutput()
	if err != nil {
		log.Fatalf("run kphp2cpp: %v: %s", err, out)
	}
	fmt.Print(string(out))
}
