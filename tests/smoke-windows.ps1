$ErrorActionPreference = "Stop"

$archivePath = if ($args.Count -eq 1) { (Resolve-Path $args[0]).Path } else { $null }
$repoRoot = Split-Path -Parent $PSScriptRoot
$tempRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("codex-next-prompt-smoke-" + [System.Guid]::NewGuid().ToString("N"))

try {
	if ($null -ne $archivePath) {
		$archive = [System.IO.Compression.ZipFile]::OpenRead($archivePath)
		try {
			foreach ($entry in $archive.Entries) {
				$segments = $entry.FullName.Replace("\", "/").Split("/")
				$unixMode = ($entry.ExternalAttributes -shr 16) -band 0xF000
				if ([System.IO.Path]::IsPathRooted($entry.FullName) -or $segments -contains ".." -or $unixMode -eq 0xA000) {
					throw "smoke: unsafe archive entry: $($entry.FullName)"
				}
			}
		}
		finally {
			$archive.Dispose()
		}
		Expand-Archive -Path $archivePath -DestinationPath $tempRoot
		$releaseRoots = @(Get-ChildItem -Path $tempRoot -Directory)
		if ($releaseRoots.Count -ne 1) {
			throw "smoke: archive must contain one release root"
		}
		$pluginRoot = $releaseRoots[0].FullName
	}
	else {
		$pluginRoot = Join-Path $tempRoot "plugin"
		$binaryDir = Join-Path $pluginRoot "bin\windows-amd64"
		$hookDir = Join-Path $pluginRoot "hooks"
		New-Item -ItemType Directory -Path $binaryDir, $hookDir | Out-Null
		Copy-Item (Join-Path $repoRoot "hooks\run.cmd") (Join-Path $hookDir "run.cmd")
		$env:GOTOOLCHAIN = "auto"
		$env:CGO_ENABLED = "0"
		$env:GOOS = "windows"
		$env:GOARCH = "amd64"
		& go build -buildvcs=false -trimpath -o (Join-Path $binaryDir "codex-next-prompt.exe") (Join-Path $repoRoot "cmd\codex-next-prompt")
		if ($LASTEXITCODE -ne 0) {
			throw "smoke: native Windows build failed"
		}
	}
	$hookDir = Join-Path $pluginRoot "hooks"
	$targetDirs = @(Get-ChildItem -Path (Join-Path $pluginRoot "bin") -Directory)
	if ($targetDirs.Count -ne 1) {
		throw "smoke: archive must contain one Windows target"
	}
	$targetArch = $targetDirs[0].Name.Replace("windows-", "")
	$nativeArchitecture = if ($env:PROCESSOR_ARCHITEW6432) { $env:PROCESSOR_ARCHITEW6432 } else { $env:PROCESSOR_ARCHITECTURE }
	$nativeArch = switch -Regex ($nativeArchitecture) {
		"^(AMD64|x86_64)$" { "amd64"; break }
		"^(ARM64|aarch64)$" { "arm64"; break }
		default { "unsupported" }
	}
	if ($targetArch -ne $nativeArch) {
		if (-not (Test-Path (Join-Path $targetDirs[0].FullName "codex-next-prompt.exe"))) {
			throw "smoke: Windows executable is missing"
		}
		Write-Output "smoke: extracted Windows $targetArch archive structure passed"
		return
	}

    $sessionOutput = '{"hook_event_name":"SessionStart","source":"startup"}' | & (Join-Path $hookDir "run.cmd") session-start
    if (($sessionOutput -notmatch '"hookEventName":"SessionStart"') -or ($sessionOutput -notmatch 'Suggested next prompt:')) {
        throw "smoke: invalid SessionStart output: $sessionOutput"
    }

    $stopOutput = '{"hook_event_name":"Stop","last_assistant_message":"Suggested next prompt:"}' | & (Join-Path $hookDir "run.cmd") stop
    if ($stopOutput -ne '{"systemMessage":"codex-next-prompt: invalid Suggested next prompt line"}') {
        throw "smoke: invalid Stop output: $stopOutput"
    }

    $savedArchitecture = $env:PROCESSOR_ARCHITECTURE
    $savedWowArchitecture = $env:PROCESSOR_ARCHITEW6432
    try {
        $env:PROCESSOR_ARCHITECTURE = "MIPS"
        Remove-Item Env:PROCESSOR_ARCHITEW6432 -ErrorAction SilentlyContinue
        $unsupportedStdout = Join-Path $tempRoot "unsupported.stdout"
        $unsupportedStderr = Join-Path $tempRoot "unsupported.stderr"
        & (Join-Path $hookDir "run.cmd") stop 1> $unsupportedStdout 2> $unsupportedStderr
        if ($LASTEXITCODE -ne 0) {
            throw "smoke: unsupported architecture must fail open with exit 0"
        }
        if ((Get-Item $unsupportedStdout).Length -ne 0) {
            throw "smoke: unsupported architecture wrote to stdout"
        }
        $unsupportedMessage = (Get-Content -Raw $unsupportedStderr).Trim()
        if ($unsupportedMessage -ne "codex-next-prompt: unsupported architecture: MIPS") {
            throw "smoke: invalid unsupported-architecture diagnostic: $unsupportedMessage"
        }
    }
    finally {
        $env:PROCESSOR_ARCHITECTURE = $savedArchitecture
        if ($null -eq $savedWowArchitecture) {
            Remove-Item Env:PROCESSOR_ARCHITEW6432 -ErrorAction SilentlyContinue
        }
        else {
            $env:PROCESSOR_ARCHITEW6432 = $savedWowArchitecture
        }
    }

    Write-Output "smoke: Windows launcher passed"
}
finally {
    Remove-Item -Recurse -Force $tempRoot -ErrorAction SilentlyContinue
}
