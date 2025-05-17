package build

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"portal/internal/server/github"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

const BuildFolder = "build"

type BuildOptions struct {
	Verbose bool
}

// installPackages calls "npm install" on the cloned repo.
func installPackages() error {
	cmd := exec.Command("npm", "install")

	cmd.Dir = github.RepoFolderName

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func copyStatic() error {
	staticFolder, err := seekFiles([]string{"public", "static"})
	if err != nil {
		return err
	}

	return copyDir(staticFolder, fmt.Sprintf("%s/%s", github.RepoFolderName, BuildFolder))
}

func runPostcss() error {
	inputCss, err := seekFiles([]string{"src/index.css", "index.css"})
	if err != nil {
		return err
	}

	inputCss = strings.TrimPrefix(inputCss, github.RepoFolderName+"/")

	cmd := exec.Command("postcss", inputCss, "-o", fmt.Sprintf("%s/%s", BuildFolder, "styles.css"))

	cmd.Dir = github.RepoFolderName

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func generateIndex() error {
	return os.WriteFile(fmt.Sprintf("%s/%s/index.html", github.RepoFolderName, BuildFolder), []byte(indexTemplate), 0644)
}

func Build(options BuildOptions) error {
	env, err := parseEnvFile(fmt.Sprintf(".env%s", "."+"test"))
	if err != nil {
		return err
	}

	if options.Verbose {
		slog.Info("parsed env file")
	}

	entryPoint, err := seekFiles([]string{"src/index.jsx", "src/index.tsx"})
	if err != nil {
		return err
	}

	if options.Verbose {
		slog.Info(fmt.Sprintf("found entryPoint %s", entryPoint))
	}

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{entryPoint},
		Outfile:     fmt.Sprintf("%s/%s", github.RepoFolderName, "build/bundle.js"),
		Bundle:      true,
		Write:       true,
		Sourcemap:   api.SourceMapInline,
		Target:      api.ESNext,
		Format:      api.FormatESModule,
		Platform:    api.PlatformBrowser,
		Define:      env.toReactVite(),
		Loader: map[string]api.Loader{
			".js":    api.LoaderJSX,
			".ts":    api.LoaderTS,
			".tsx":   api.LoaderTSX,
			".css":   api.LoaderCSS,
			".png":   api.LoaderFile,
			".svg":   api.LoaderFile,
			".jpg":   api.LoaderFile,
			".gif":   api.LoaderFile,
			".eot":   api.LoaderFile,
			".woff":  api.LoaderFile,
			".woff2": api.LoaderFile,
			".otf":   api.LoaderFile,
			".ttf":   api.LoaderFile,
		},
	})

	if options.Verbose {
		for _, warning := range result.Warnings {
			slog.Warn(warning.Text)
		}
	}

	if len(result.Errors) > 0 {
		for _, err := range result.Errors {
			slog.Error(err.Text)
		}
		return errors.New("esbuild errors")
	}

	if options.Verbose {
		slog.Info("esbuild done")
	}

	return nil
}

func TotalBuild() error {
	slog.Info("performing npm instal...")

	err := installPackages()
	if err != nil {
		return err
	}

	slog.Info("performing esbuild...")
	err = Build(BuildOptions{Verbose: false})
	if err != nil {
		return err
	}

	slog.Info("generating index.html...")
	err = generateIndex()
	if err != nil {
		return err
	}

	slog.Info("running postcss...")
	err = runPostcss()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info("copying static files...")
	err = copyStatic()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}
