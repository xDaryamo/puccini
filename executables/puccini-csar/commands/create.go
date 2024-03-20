package commands

import (
	"archive/tar"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/zip"
	"github.com/klauspost/pgzip"
	"github.com/spf13/cobra"
	"github.com/tliron/exturl"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca/csar"
)

var (
	compressionLevel     int
	toscaMetaFileVersion string
	csarVersion          string
	createdBy            string
	entryDefinitions     string
	otherDefinitions     []string
)

func init() {
	rootCommand.AddCommand(createCommand)

	createCommand.Flags().IntVarP(&compressionLevel, "compression", "c", 6, "compression level (0 to 9, where 0 is no compression and 9 is maximum compression)")
	createCommand.Flags().StringVarP(&archiveFormat, "archive-format", "a", "", "force archive format (\"tar.gz\", \"tar\", or \"zip\"); leave empty to determine automatically from extension")
	createCommand.Flags().StringVar(&toscaMetaFileVersion, "tosca-meta-file-version", "1.1", "TOSCA-Meta-File-Version field")
	createCommand.Flags().StringVar(&csarVersion, "csar-version", "1.1", "CSAR-Version field")
	createCommand.Flags().StringVar(&createdBy, "created-by", toolName, "Created-By field")
	createCommand.Flags().StringVar(&entryDefinitions, "entry-definitions", "", "Entry-Definitions field; leave empty to use root YAML file; if more then one root YAML exists then must be set")
	createCommand.Flags().StringArrayVar(&otherDefinitions, "other-definitions", nil, "Other-Definitions field")
}

var createCommand = &cobra.Command{
	Use:   "create [CSAR PATH] [BASE DIRECTORY PATH]",
	Short: "Create CSAR",
	Long:  `Creates a CSAR from a directory.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		csarPath := args[0]
		dir := args[1]

		CreateCSAR(csarPath, dir)
	},
}

func CreateCSAR(csarPath string, dir string) {
	if (compressionLevel < 0) || (compressionLevel > 9) {
		util.Failf("invalid compression level, must be >=0 and <=9: %d", compressionLevel)
	}

	if archiveFormat == "" {
		archiveFormat = exturl.GetFormat(csarPath)
	}

	if !csar.IsValidFormat(archiveFormat) {
		util.Failf("unsupported CSAR archive format: %q", archiveFormat)
	}

	stat, err := os.Stat(dir)
	util.FailOnError(err)
	if !stat.IsDir() {
		util.Failf("not a directory: %s", dir)
	}

	file, err := os.OpenFile(csarPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	util.FailOnError(err)

	switch archiveFormat {
	case "tar":
		CreateTarCSAR(dir, file)

	case "tar.gz":
		CreateGzipTarCSAR(dir, file)

	case "zip", "csar":
		CreateZipCSAR(dir, file)
	}

	util.OnExit(func() {
		log.Noticef("created CSAR: %s", csarPath)
	})
}

func CreateTarCSAR(dir string, writer io.Writer) {
	tarWriter := tar.NewWriter(writer)
	util.OnExitError(tarWriter.Close)

	createCsar(dir, func(internalPath string, buffer []byte, file *os.File) error {
		internalPath = filepath.ToSlash(internalPath)

		header := tar.Header{
			Name: internalPath,
			Mode: 0600,
		}

		if buffer != nil {
			header.Size = int64(len(buffer))
		} else {
			if stat, err := file.Stat(); err == nil {
				header.Size = stat.Size()
			} else {
				return err
			}
		}

		if err := tarWriter.WriteHeader(&header); err == nil {
			if buffer != nil {
				_, err = tarWriter.Write(buffer)
			} else {
				_, err = io.Copy(tarWriter, file)
			}
			return err
		} else {
			return err
		}
	})
}

func CreateGzipTarCSAR(dir string, writer io.Writer) {
	gzipWriter, err := pgzip.NewWriterLevel(writer, compressionLevel)
	util.FailOnError(err)
	util.OnExitError(gzipWriter.Close)

	CreateTarCSAR(dir, gzipWriter)
}

func CreateZipCSAR(dir string, file *os.File) {
	zipWriter := zip.NewWriter(file)
	util.OnExitError(zipWriter.Close)

	zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, compressionLevel)
	})

	createCsar(dir, func(internalPath string, buffer []byte, file *os.File) error {
		internalPath = filepath.ToSlash(internalPath)
		if writer, err := zipWriter.Create(internalPath); err == nil {
			if buffer != nil {
				_, err = writer.Write(buffer)
			} else {
				_, err = io.Copy(writer, file)
			}
			return err
		} else {
			return err
		}
	})
}

func createCsar(dir string, writeEntry func(string, []byte, *os.File) error) {
	prefix := len(dir) + 1
	var hasMeta bool

	err := filepath.WalkDir(dir, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !dirEntry.IsDir() {
			internalPath := path[prefix:]
			if internalPath == csar.TOSCA_META_PATH {
				// Validate meta
				_, err = csar.ReadMetaFromPath(path)
				util.FailOnError(err)

				hasMeta = true
				log.Infof("using included %s", csar.TOSCA_META_PATH)
			}
			log.Infof("adding: %s", internalPath)
			file, err := os.Open(path)
			util.FailOnError(err)
			defer file.Close()
			return writeEntry(internalPath, nil, file)
		}

		return nil
	})
	util.FailOnError(err)

	if !hasMeta {
		if entryDefinitions == "" {
			dirEntries, err := os.ReadDir(dir)
			util.FailOnError(err)
			for _, dirEntry := range dirEntries {
				name := dirEntry.Name()
				if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
					if entryDefinitions != "" {
						util.Failf("dir has more than one potential service template at the root: %s", dir)
					} else {
						entryDefinitions = name
					}
				}
			}
		}

		log.Infof("generating new %s", csar.TOSCA_META_PATH)

		toscaMetaFileVersion_, err := csar.ParseVersion(toscaMetaFileVersion)
		util.FailOnError(err)
		csarVersion_, err := csar.ParseVersion(csarVersion)
		util.FailOnError(err)

		meta := csar.Meta{
			Version:          toscaMetaFileVersion_,
			CsarVersion:      csarVersion_,
			CreatedBy:        createdBy,
			EntryDefinitions: entryDefinitions,
			OtherDefinitions: otherDefinitions,
		}

		meta_, err := meta.ToBytes()
		util.FailOnError(err)

		err = writeEntry(csar.TOSCA_META_PATH, meta_, nil)
		util.FailOnError(err)
	}
}
