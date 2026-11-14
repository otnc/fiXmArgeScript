# CI で作ったセットアップ (dist\fixmarge_0.0.0-ci_windows_setup.exe) を黙って入れて消し、
# fixmarge が動くこと、PATH が足されて元に戻ること、ファイルが消えることを確かめる。
$ErrorActionPreference = "Stop"

$app = "$env:LOCALAPPDATA\Programs\fixmarge"
function UserPath { [Environment]::GetEnvironmentVariable("Path", "User") }

Start-Process dist\fixmarge_0.0.0-ci_windows_setup.exe -ArgumentList "/VERYSILENT", "/SUPPRESSMSGBOXES", "/CURRENTUSER" -Wait
& "$app\fixmarge.exe" gun
if ($LASTEXITCODE -ne 0) { throw "fixmarge gun failed" }
if ((UserPath) -split ";" -notcontains $app) { throw "PATH does not contain $app" }

Start-Process "$app\unins000.exe" -ArgumentList "/VERYSILENT", "/SUPPRESSMSGBOXES" -Wait
# アンインストーラは自分を消すために別のプロセスを残して先に終わるので、少し待つ
Start-Sleep -Seconds 3
if ((UserPath) -split ";" -contains $app) { throw "PATH still contains $app" }
if (Test-Path "$app\fixmarge.exe") { throw "fixmarge.exe was not removed" }
