; Windows のセットアップ fixmarge_<版>_windows_setup.exe を作る。
; iscc /DAppVersion=0.2.0 installer\fixmarge.iss で動かす。BinDir の amd64\ と arm64\ に fixmarge.exe を置いておく。
; 1つの exe に両方の CPU 向けを入れ、入れるときに選ぶ。

#ifndef AppVersion
  #error AppVersion を /DAppVersion=0.2.0 のように渡してください
#endif
#ifndef BinDir
  #define BinDir "..\dist\windows"
#endif
; Explorer の「ファイル バージョン」(VersionInfoVersion) は数字4つだけなので、プレリリースかどうかを4つ目で表す。
; 正式版は最大の 65535 にし、プレリリースは段階×1000+番号にする (alpha.N は 1000+N、beta.N は 2000+N、rc.N は 3000+N)。
; こうすると 0.2.0-rc.1 (0.2.0.3001) が 0.2.0 (0.2.0.65535) より古く並ぶ。識別子は scripts/release/next-version.sh がこの3つに限っている。
; それ以外 (CI の 0.0.0-ci) は 0。製品バージョンの文字列は AppVersion のまま (-rc.1 付き) になる
#define BaseVersion AppVersion
#define Stage "65535"
; 前処理は通らない #if の中の式も読むので、正式版でも PreId と PreNum を定義しておく
#define PreId ""
#define PreNum 0
#if Pos("-", AppVersion) > 0
  #define BaseVersion Copy(AppVersion, 1, Pos("-", AppVersion) - 1)
  #define Pre Copy(AppVersion, Pos("-", AppVersion) + 1, Len(AppVersion))
  #if Pos(".", Pre) > 0
    #define PreId Copy(Pre, 1, Pos(".", Pre) - 1)
    #define PreNum Int(Copy(Pre, Pos(".", Pre) + 1, Len(Pre)), 0)
  #else
    #define PreId Pre
    #define PreNum 0
  #endif
  #if PreId == "alpha"
    #define Stage Str(1000 + PreNum)
  #elif PreId == "beta"
    #define Stage Str(2000 + PreNum)
  #elif PreId == "rc"
    #define Stage Str(3000 + PreNum)
  #else
    #define Stage "0"
  #endif
#endif
#define NumericVersion BaseVersion + "." + Stage

[Setup]
; 「設定 → アプリ」の登録キー (fiXmArgeScript_is1) になる。上書きとアンインストールがこれで見分けるので変えない
AppId=fiXmArgeScript
AppName=fiXmArge Script
AppVersion={#AppVersion}
AppPublisher=otoneko.
AppPublisherURL=https://github.com/otnc/fiXmArgeScript
AppSupportURL=https://github.com/otnc/fiXmArgeScript/issues
VersionInfoVersion={#NumericVersion}
DefaultDirName={autopf}\fixmarge
DisableProgramGroupPage=yes
; 既定は管理者権限なしで自分だけ ({localappdata}\Programs\fixmarge)。起動時に全員向けも選べる
PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog
ArchitecturesAllowed=x64compatible or arm64
ArchitecturesInstallIn64BitMode=x64compatible or arm64
ChangesEnvironment=yes
LicenseFile=license.txt
SetupIconFile=fixmarge.ico
UninstallDisplayIcon={app}\fixmarge.ico
OutputDir=..\dist
OutputBaseFilename=fixmarge_{#AppVersion}_windows_setup
SolidCompression=yes
WizardStyle=modern

[Languages]
Name: "japanese"; MessagesFile: "compiler:Languages\Japanese.isl"

[Files]
Source: "{#BinDir}\amd64\fixmarge.exe"; DestDir: "{app}"; Check: not IsArm64; Flags: ignoreversion
Source: "{#BinDir}\arm64\fixmarge.exe"; DestDir: "{app}"; Check: IsArm64; Flags: ignoreversion
Source: "fixmarge.ico"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\README.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\LICENSE-MIT"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\LICENSE-APACHE"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\THIRD_PARTY_LICENSES"; DestDir: "{app}"; Flags: ignoreversion

[Code]
const
  UserEnvironment = 'Environment';
  SystemEnvironment = 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment';

// 自分だけなら自分の PATH、全員向けならシステムの PATH を書き換える
function EnvironmentRoot: Integer;
begin
  if IsAdminInstallMode then
    Result := HKEY_LOCAL_MACHINE
  else
    Result := HKEY_CURRENT_USER;
end;

function EnvironmentKey: String;
begin
  if IsAdminInstallMode then
    Result := SystemEnvironment
  else
    Result := UserEnvironment;
end;

function ReadPath: String;
begin
  if not RegQueryStringValue(EnvironmentRoot, EnvironmentKey, 'Path', Result) then
    Result := '';
end;

// 前後に ; を付けて比べ、C:\a が C:\ab に一致しないようにする
function PathHas(Paths, Dir: String): Boolean;
begin
  Result := Pos(';' + Uppercase(Dir) + ';', ';' + Uppercase(Paths) + ';') > 0;
end;

procedure AddToPath(Dir: String);
var
  Paths: String;
begin
  Paths := ReadPath;
  if PathHas(Paths, Dir) then
    Exit;
  if (Paths <> '') and (Paths[Length(Paths)] <> ';') then
    Paths := Paths + ';';
  // %USERPROFILE% などを含む既存の値を壊さないように、展開しない形で書く
  RegWriteExpandStringValue(EnvironmentRoot, EnvironmentKey, 'Path', Paths + Dir);
end;

procedure RemoveFromPath(Dir: String);
var
  Paths: String;
begin
  Paths := ';' + ReadPath + ';';
  if StringChangeEx(Paths, ';' + Dir + ';', ';', True) = 0 then
    Exit;
  Paths := Copy(Paths, 2, Length(Paths) - 2);
  RegWriteExpandStringValue(EnvironmentRoot, EnvironmentKey, 'Path', Paths);
end;

procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall then
    AddToPath(ExpandConstant('{app}'));
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall then
    RemoveFromPath(ExpandConstant('{app}'));
end;
