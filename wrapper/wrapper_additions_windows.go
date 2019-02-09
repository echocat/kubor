// +build windows

package wrapper

import "github.com/levertonai/kubor/log"

func noticeAfterCreation(unixScriptFile string) {
	log.Warn("You created successfully the kubor wrapper - this includes the file '%s'."+
		" This file is made for UNIX like operation systems like Linux and macOS."+
		" This systems you need an executable flag to execute it seamless."+
		" But you are currently running on Windows which is not able to create this flag for you."+
		" If your project is based on Git you can simply run following command to set it to the Git repository index: git update-index --chmod=+x %s", unixScriptFile, unixScriptFile)
}
