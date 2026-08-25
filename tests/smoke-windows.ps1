$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$tempRoot = Join-Path ([System.IO.Path]::GetTempPath()) ("codex-next-prompt-smoke-" + [System.Guid]::NewGuid().ToString("N"))
$pluginRoot = Join-Path $tempRoot "plugin"
$binaryDir = Join-Path $pluginRoot "bin\windows-amd64"
$hookDir = Join-Path $pluginRoot "hooks"

try {
    New-Item -ItemType Directory -Path $binaryDir, $hookDir | Out-Null
    Copy-Item (Join-Path $repoRoot "hooks\run.cmd") (Join-Path $hookDir "run.cmd")

    $env:GOTOOLCHAIN = "auto"
    $env:CGO_ENABLED = "0"
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    & go build -trimpath -o (Join-Path $binaryDir "codex-next-prompt.exe") (Join-Path $repoRoot "cmd\codex-next-prompt")
    if ($LASTEXITCODE -ne 0) {
        throw "smoke: native Windows build failed"
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
