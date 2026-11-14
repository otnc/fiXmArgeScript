# セットアップのファイル バージョンの4つ目が、プレリリースの段階と番号の決まり (installer/fixmarge.iss) どおりかを確かめる。
# 0.0.0-ci のセットアップは先に作ってあるものを使い、ほかの版はここで作る。
$ErrorActionPreference = "Stop"

$iscc = "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe"
$want = [ordered]@{
  "0.0.0-ci"      = "0.0.0.0"
  "1.2.3-alpha.0" = "1.2.3.1000"
  "1.2.3-beta.12" = "1.2.3.2012"
  "1.2.3-rc.4"    = "1.2.3.3004"
  "1.2.3"         = "1.2.3.65535"
}

foreach ($v in $want.Keys) {
  if ($v -ne "0.0.0-ci") {
    & $iscc /Q "/DAppVersion=$v" installer\fixmarge.iss
    if ($LASTEXITCODE -ne 0) { throw "ISCC failed for $v" }
  }
  $info = (Get-Item "dist\fixmarge_${v}_windows_setup.exe").VersionInfo
  $file = "$($info.FileMajorPart).$($info.FileMinorPart).$($info.FileBuildPart).$($info.FilePrivatePart)"
  if ($file -ne $want[$v]) { throw "${v}: file version $file, want $($want[$v])" }
  # Inno Setup は文字列の後ろを空白で埋めて書くので、外して比べる
  if ($info.ProductVersion.Trim() -ne $v) { throw "${v}: product version '$($info.ProductVersion)', want $v" }
}
