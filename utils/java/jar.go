package java

// Represents the configuration required
// for executing the command 'jar'
type JarCommand struct {
	JarPath        string
	OutputFile     string
	Sources        [][]string
	ReleaseVersion string
	ManifestFile   string
	MainClass      string
}

func NewJarCommand(jarPath, outputFile string, sources [][]string, mainClass string) JarCommand {
	return JarCommand{
		JarPath:    jarPath,
		OutputFile: outputFile,
		Sources:    sources,
		MainClass:  mainClass,
	}
}

// Returns the arguments for the 'jar' command
func (j JarCommand) Arguments() []string {
	args := make([]string, 0)
	args = append(args, "--create", "--file", j.OutputFile)

	if j.MainClass != "" {
		args = append(args, "--main-class", j.MainClass)
	}

	if j.ManifestFile != "" {
		args = append(args, "--manifest", j.ManifestFile)
	}

	for _, sourceBundle := range j.Sources {
		if len(sourceBundle) == 1 {
			args = append(args, "-C", sourceBundle[0], ".")
		} else if len(sourceBundle) > 1 {
			args = append(args, "-C")
			args = append(args, sourceBundle[0:]...)
		}
	}

	return args
}
