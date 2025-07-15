package helpers

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"github.com/dev-kas/xel/shared"
)

func ExtractTarGz(tarGzPath, targetDir string) error {
	f, err := os.Open(tarGzPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // done
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Join(targetDir, header.Name)
		cleanPath := filepath.Clean(targetPath)

		if !strings.HasPrefix(cleanPath, filepath.Clean(targetDir)) {
			return fmt.Errorf("possible path traversal attempt: %s", targetPath)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			err := os.Mkdir(targetPath, 0755)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			err := os.MkdirAll(filepath.Dir(targetPath), 0755)
			if err != nil {
				return err
			}

			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, tr)
			outFile.Close()
			if err != nil {
				return err
			}
		default:
			if header.Name == "pax_global_header" {
				break
			}
			shared.ColorPalette.Warning.Printf("Skipping unknown file type %c for %s\n", header.Typeflag, header.Name)
		}
	}

	return nil
}

var integrityHashers = map[string]func() hash.Hash{
	"sha256": sha256.New,
	"sha512": sha512.New,
	"md5":    md5.New,
}

func VerifyTarballIntegrity(tarballPath, algo, expectedHash string) error {
	hasherFn, ok := integrityHashers[algo]
	if !ok {
		return fmt.Errorf("unsupported integrity algorithm: %s", algo)
	}

	f, err := os.Open(tarballPath)
	if err != nil {
		return fmt.Errorf("failed to open tarball: %w", err)
	}
	defer f.Close()

	hasher := hasherFn()
	if _, err := io.Copy(hasher, f); err != nil {
		return fmt.Errorf("failed to hash file: %w", err)
	}

	actualHash := hex.EncodeToString(hasher.Sum(nil))
	if actualHash != expectedHash {
		return fmt.Errorf(
			"integrity check failed for %s:\nexpected: %s\ngot:      %s",
			algo, expectedHash, actualHash,
		)
	}

	return nil
}
