package main

type input struct {
	Type string `yaml:"type"`
	Spec struct {
		Args     []string `yaml:"args,omitempty"`
		ExecCmd  string   `yaml:"exec,omitempty"`
		Sudo     bool     `yaml:"sudo,omitempty"`
		Path     string   `yaml:"path,omitempty"`
		Encoding string   `yaml:"encoding,omitempty"`
	} `yaml:"spec"`
}

type inputs struct {
	Ins []input `yaml:"inputs"`
}
