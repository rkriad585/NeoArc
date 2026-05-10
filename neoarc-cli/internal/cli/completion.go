package cli

import (
	"fmt"
	"os"
	"strings"
)

func GenerateCompletion(shell string, aliases []string) int {
	switch shell {
	case "bash":
		return genBash(aliases)
	case "zsh":
		return genZsh(aliases)
	case "powershell":
		return genPowerShell(aliases)
	default:
		fmt.Fprintf(os.Stderr, "Unknown shell: %s\nSupported shells: bash, zsh, powershell\n", shell)
		return 1
	}
}

func aliasQuoted(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return `"` + strings.Join(names, `" "`) + `"`
}

func genBash(aliases []string) int {
	aq := aliasQuoted(aliases)
	fmt.Printf(`_neoarc_completions() {
	local cur prev opts
	COMPREPLY=()
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[COMP_CWORD-1]}"

	opts="get run config config-token help completion update"

	if [[ ${prev} == "completion" ]]; then
		COMPREPLY=($(compgen -W "bash zsh powershell" -- ${cur}))
		return 0
	fi

	if [[ ${prev} == "config" ]]; then
		COMPREPLY=($(compgen -W "insecure secure" -- ${cur}))
		return 0
	fi

	if [[ ${prev} == "get" || ${prev} == "run" ]]; then
`)
	if aq != "" {
		fmt.Fprintf(os.Stdout, `		COMPREPLY=($(compgen -W "%s" -- ${cur}))
`, aq)
	} else {
		fmt.Fprintf(os.Stdout, "		COMPREPLY=()\n")
	}
	fmt.Printf(`		return 0
	fi

	if [[ ${cur} == -* ]]; then
		COMPREPLY=($(compgen -W "--dry-run --yes" -- ${cur}))
		return 0
	fi

	COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
	return 0
}

complete -F _neoarc_completions neoarc
`)
	return 0
}

func genZsh(aliases []string) int {
	aq := aliasQuoted(aliases)
	fmt.Printf(`#compdef neoarc

_neoarc_completions() {
	local -a opts
	opts=("get:Print the code for an alias" "run:Run an alias command" "config:Set server URL or options" "config-token:Set API token" "help:Show help" "completion:Generate completion script" "update:Self-update the binary")

	_arguments -C \
		'--dry-run[Print the alias code without executing]' \
		'--yes[Skip execution confirmation prompt]' \
		'1: :->cmds' \
		'*: :->args' \
		&& return 0

	case $state in
		cmds)
			_describe -t commands "neoarc subcommands" opts
			;;
		args)
			case ${words[1]} in
				get|run)
`)
	if aq != "" {
		fmt.Fprintf(os.Stdout, `					_values "aliases" %s
`, aq)
	}
	fmt.Printf(`					;;
				completion)
					_values "shells" bash zsh powershell
					;;
				config)
					_values "options" insecure secure
					;;
			esac
			;;
	esac
}

_neoarc_completions
`)
	return 0
}

func genPowerShell(aliases []string) int {
	aq := aliasQuoted(aliases)
	fmt.Printf(`using namespace System.Management.Automation
using namespace System.Collections.Generic

Register-ArgumentCompleter -Native -CommandName neoarc -ScriptBlock {
	param($wordToComplete, $commandAst, $cursorPosition)

	$commands = @("get", "run", "config", "config-token", "help", "completion", "update")
	$completionCmd = $commandAst.CommandElements | Where-Object { $_ -is [CommandElementAst] -and $_.Extent.Text -notlike "-*" }
	$currentIndex = $commandAst.CommandElements.IndexOf($commandAst.CommandElements | Where-Object { $_.Extent.EndOffset -eq $cursorPosition })
	$prevIndex = $currentIndex - 1

	if ($prevIndex -ge 0 -and $currentIndex -ge 0) {
		$prev = $commandAst.CommandElements[$prevIndex].Extent.Text
		if ($prev -eq "completion") {
			@("bash", "zsh", "powershell") | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
				[CompletionResult]::new($_)
			}
			return
		}
		if ($prev -eq "config") {
			@("insecure", "secure") | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
				[CompletionResult]::new($_)
			}
			return
		}
		if ($prev -eq "get" -or $prev -eq "run") {
`)
	if aq != "" {
		fmt.Fprintf(os.Stdout, `			@(%s) | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
`, aq)
		fmt.Fprintf(os.Stdout, `				[CompletionResult]::new($_)
`)
		fmt.Fprintf(os.Stdout, `			}
`)
	} else {
		fmt.Fprintf(os.Stdout, `			# no aliases available
`)
	}
	fmt.Printf(`			return
		}
	}

	if ($wordToComplete -like "-*") {
		@("--dry-run", "--yes") | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
			[CompletionResult]::new($_)
		}
		return
	}

	$commands | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
		[CompletionResult]::new($_)
	}
}
`)
	return 0
}
