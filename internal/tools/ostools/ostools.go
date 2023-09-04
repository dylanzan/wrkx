package ostools

import (
	"sshx/internal/cores/sshx"
	"strings"
)

const (
	_fedora = iota
	_centos
	_suse
	_redhat
	_debian
	_ubuntu
)

const (
	_CHECK_VERSION_COMMAND = "sudo cat /etc/os-release"
	_TAG_NAME              = "NAME"
	_TAG_VERSION           = "VERSION"
	_TAG_VERSION_ID_LIKE   = "ID_LIKE"
)

func _getOsVersionTagList() []string {
	return []string{_TAG_NAME, _TAG_VERSION, _TAG_VERSION_ID_LIKE}
}

/*
centos:
NAME="CentOS Linux"
VERSION="7 (Core)"
ID="centos"
ID_LIKE="rhel fedora"
VERSION_ID="7"
PRETTY_NAME="CentOS Linux 7 (Core)"
ANSI_COLOR="0;31"
CPE_NAME="cpe:/o:centos:centos:7"
HOME_URL="https://www.centos.org/"
BUG_REPORT_URL="https://bugs.centos.org/"

CENTOS_MANTISBT_PROJECT="CentOS-7"
CENTOS_MANTISBT_PROJECT_VERSION="7"
REDHAT_SUPPORT_PRODUCT="centos"
REDHAT_SUPPORT_PRODUCT_VERSION="7"

ubuntu:
PRETTY_NAME="Ubuntu 22.04.2 LTS"
NAME="Ubuntu"
VERSION_ID="22.04"
VERSION="22.04.2 LTS (Jammy Jellyfish)"
VERSION_CODENAME=jammy
ID=ubuntu
ID_LIKE=debian
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
UBUNTU_CODENAME=jammy
*/

// GetOSVersionNmae get os version name
func GetOSVersionName(cli *sshx.Cli) (osversion string, err error) {

	stdout, err := cli.Run(_CHECK_VERSION_COMMAND)
	if err != nil {
		return "", err
	}
	stdouts := strings.Split(stdout, "\n")
	//fmt.Sprintln(splits)
	for _, stdoutVal := range stdouts {
		stdoutVals := strings.Split(stdoutVal, "=")

		for _, tag := range _getOsVersionTagList() {
			if stdoutVals[0] == tag {
				osversion += strings.Trim(stdoutVals[1], "\"") + "_"
			}
		}
	}

	return strings.Trim(osversion, "_"), nil
}
