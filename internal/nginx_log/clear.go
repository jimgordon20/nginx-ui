package nginx_log

import (
	"fmt"
	"os"

	"github.com/0xJacky/Nginx-UI/internal/nginx"
	"github.com/0xJacky/Nginx-UI/internal/nginx_log/utils"
)

// ClearDomainLog truncates a single access/error log file in place so Nginx can
// keep writing to the same inode. The path must pass the log whitelist.
func ClearDomainLog(logPath string) error {
	if logPath == "" {
		return fmt.Errorf("log path is required")
	}
	if !utils.IsValidLogPath(logPath) {
		return fmt.Errorf("log path is not under whitelist")
	}

	info, err := nginx.Stat(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("log file does not exist: %s", logPath)
		}
		return fmt.Errorf("stat log file: %w", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("log path is not a regular file")
	}

	f, err := nginx.OpenFile(logPath, os.O_WRONLY|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return fmt.Errorf("truncate log file: %w", err)
	}
	if cerr := f.Close(); cerr != nil {
		return fmt.Errorf("close truncated log file: %w", cerr)
	}
	return nil
}
