# build-images.ps1
# Requirements:
# - Go toolchain installed and on PATH
# - Docker installed and running
# - Run from the Go module root (main package)

param(
  [string]$AppName = "qs-backend",
  [string]$ImageRepo = "chetani007/quantum-scholar-backend",
  [string]$VersionTag = "latest",
  [string]$ContextDir = ".",
  [string]$Port = "8000",
  [string]$AlpineTag = "latest"
)

$ErrorActionPreference = "Stop"

# 1) Output directory
$OutDir = Join-Path -Path $PSScriptRoot -ChildPath "dist"
if (Test-Path $OutDir) { Remove-Item -Recurse -Force $OutDir }
New-Item -ItemType Directory -Path $OutDir | Out-Null

Write-Host "Building Go binaries for linux/amd64 and linux/arm64..."

# 2) Cross-compile Go binaries
# If CGO is used, the Alpine runtime (musl) is a good match.
# For static builds, set $env:CGO_ENABLED="0".
$env:GOOS = "linux"

$env:GOARCH = "amd64"
go build -o (Join-Path $OutDir "$AppName-linux-amd64.out") $ContextDir
Write-Host "Built $AppName-linux-amd64.out"

$env:GOARCH = "arm64"
go build -o (Join-Path $OutDir "$AppName-linux-arm64.out") $ContextDir
Write-Host "Built $AppName-linux-arm64.out"

Remove-Item Env:GOOS -ErrorAction Ignore
Remove-Item Env:GOARCH -ErrorAction Ignore

$MailSrc1 = Join-Path $PSScriptRoot "mail"
Copy-Item -Recurse -Force $MailSrc1 (Join-Path $OutDir "mail")

# 3) Create Alpine-based Dockerfiles per arch
# - Adds ca-certificates
# - Creates non-root user
# - Copies the matching binary
$DockerfileAmd64 = @"
FROM alpine:$AlpineTag
RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S -G app -u 10000 app
COPY $AppName-linux-amd64.out /$AppName
COPY mail/ /mail/
EXPOSE $Port
USER app
ENTRYPOINT ["/$AppName"]
"@

$DockerfileArm64 = @"
FROM alpine:$AlpineTag
RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S -G app -u 10000 app
COPY $AppName-linux-arm64.out /$AppName
COPY mail/ /mail/
EXPOSE $Port
USER app
ENTRYPOINT ["/$AppName"]
"@

$DfDir = Join-Path $OutDir "docker"
New-Item -ItemType Directory -Path $DfDir | Out-Null
Set-Content -Path (Join-Path $DfDir "Dockerfile.amd64") -Value $DockerfileAmd64 -NoNewline
Set-Content -Path (Join-Path $DfDir "Dockerfile.arm64") -Value $DockerfileArm64 -NoNewline

# 4) Build images per architecture
$tagAmd64 = '{0}:{1}-{2}' -f $ImageRepo, 'amd64', $VersionTag
$tagArm64 = '{0}:{1}-{2}' -f $ImageRepo, 'arm64', $VersionTag

Write-Host "Building Docker image for amd64: $tagAmd64"
docker build `
  --file (Join-Path $DfDir "Dockerfile.amd64") `
  --tag $tagAmd64 `
  --platform linux/amd64 `
  $OutDir

Write-Host "Building Docker image for arm64: $tagArm64"
docker build `
  --file (Join-Path $DfDir "Dockerfile.arm64") `
  --tag $tagArm64 `
  --platform linux/arm64 `
  $OutDir

Write-Host "Done. Built images:"
Write-Host " - $tagAmd64"
Write-Host " - $tagArm64"
