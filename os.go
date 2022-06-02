package goutil

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Exec(wd, stdoutPath, binPath string, args ...string) error {
	if wd == "" {
		wd = "."
	}

	cmd := exec.Command(binPath, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Dir = wd

	if stdoutPath != "" {
		stdoutFile, err := os.Create(stdoutPath)
		if err != nil {
			return err
		}
		defer stdoutFile.Close()

		cmd.Stdout = stdoutFile
	} else {
		cmd.Stdout = os.Stdout
	}

	fmt.Printf("CMD: %+v\n", cmd)
	fmt.Printf(" wd: %s\n", wd)
	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()

	/*go func() {
		err := cmd.Wait()
		if err != nil {
			fmt.Println(nfType, sid, "error:", err)
		} else {
			fmt.Println(nfType, sid, "finished")
		}
	}()

	return nil
	*/
}
