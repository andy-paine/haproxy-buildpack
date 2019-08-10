package supply

import (
  "os"
  "os/exec"
)

type CompileCommand struct {}

func (c *CompileCommand) Run(dir string) error {
  cmd := exec.Command(
    "make",
    "TARGET=linux-glibc",
    "EXTRA_OBJS=contrib/prometheus-exporter/service-prometheus.o",
    "USE_OPENSSL=1",
  )
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	return cmd.Run()
}
