// note: this code is only used to check whether `cifs.utils` is installed
package storage

import (
	"fmt"
	"os/exec"
)

func CheckDependencies() error {
	_, err := exec.LookPath("mount.cifs")
	if err != nil {
		return fmt.Errorf("机器中未找到`cifs.utils`，请先安装。")
	}
	return nil
}
