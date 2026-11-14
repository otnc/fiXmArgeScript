# Releases の Windows 向けの zip からセットアップを作り、checksums.txt に足して、どちらも Releases に上げる。
# タグは環境変数 TAG (v0.2.0 など)。gh を使うので GH_TOKEN も要る。Inno Setup は Windows でしか動かないため、GoReleaser とは別に作る。
$ErrorActionPreference = "Stop"

$version = $env:TAG.Substring(1)
$setup = "fixmarge_${version}_windows_setup.exe"

foreach ($arch in "amd64", "arm64") {
  gh release download $env:TAG -p "fixmarge_${version}_windows_${arch}.zip" -D dist
  if ($LASTEXITCODE -ne 0) { throw "download failed: $arch" }
  Expand-Archive "dist\fixmarge_${version}_windows_${arch}.zip" "dist\windows\$arch"
}
Copy-Item dist\windows\amd64\THIRD_PARTY_LICENSES .
& "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe" "/DAppVersion=$version" installer\fixmarge.iss
if ($LASTEXITCODE -ne 0) { throw "ISCC failed" }

gh release download $env:TAG -p checksums.txt -D dist --clobber
if ($LASTEXITCODE -ne 0) { throw "download failed: checksums.txt" }
$hash = (Get-FileHash "dist\$setup" -Algorithm SHA256).Hash.ToLower()
# やり直しで前の行が残っていれば外してから足す。GoReleaser と同じ形 (ハッシュ、空白2つ、名前) で、改行は LF
$lines = [IO.File]::ReadAllLines("dist\checksums.txt") | Where-Object { $_ -and -not $_.EndsWith("  $setup") }
[IO.File]::WriteAllText("dist\checksums.txt", ((@($lines) + "$hash  $setup") -join "`n") + "`n")
gh release upload $env:TAG "dist\$setup" dist\checksums.txt --clobber
if ($LASTEXITCODE -ne 0) { throw "upload failed" }
