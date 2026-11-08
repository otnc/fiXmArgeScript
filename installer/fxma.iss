; Windows のセットアップ fxma_<版>_windows_setup.exe を作る。
; iscc /DAppVersion=0.2.0 installer\fxma.iss で動かす。BinDir の amd64\ と arm64\ に fxma.exe を置いておく。
; 1つの exe に両方の CPU 向けを入れ、入れるときに選ぶ。

#ifndef AppVersion
  #error AppVersion を /DAppVersion=0.2.0 のように渡してください
#endif
#ifndef BinDir
  #define BinDir "..\dist\windows"
#endif
; ファイルの版情報は数字しか書けないので、0.2.0-rc.0 の -rc.0 を外す
#if Pos("-", AppVersion) > 0
  #define NumericVersion Copy(AppVersion, 1, Pos("-", AppVersion) - 1)
#else
  #define NumericVersion AppVersion
#endif

[Setup]
AppId={{6C8863CF-DD6A-40B8-8436-3D5D380F15ED}
AppName=fiXmArge Script
AppVersion={#AppVersion}
AppPublisher=otnc (非公式のファンメイド作品)
AppPublisherURL=https://github.com/otnc/fiXmArgeScript
AppSupportURL=https://github.com/otnc/fiXmArgeScript/issues
VersionInfoVersion={#NumericVersion}
DefaultDirName={autopf}\fxma
DisableProgramGroupPage=yes
; 既定は管理者権限なしで自分だけ ({localappdata}\Programs\fxma)。起動時に全員向けも選べる
PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog
ArchitecturesAllowed=x64compatible or arm64
ArchitecturesInstallIn64BitMode=x64compatible or arm64
ChangesEnvironment=yes
LicenseFile=license.txt
UninstallDisplayIcon={app}\fxma.exe
OutputDir=..\dist
OutputBaseFilename=fxma_{#AppVersion}_windows_setup
SolidCompression=yes
WizardStyle=modern

[Languages]
Name: "japanese"; MessagesFile: "compiler:Languages\Japanese.isl"

[Files]
Source: "{#BinDir}\amd64\fxma.exe"; DestDir: "{app}"; Check: not IsArm64; Flags: ignoreversion
Source: "{#BinDir}\arm64\fxma.exe"; DestDir: "{app}"; Check: IsArm64; Flags: ignoreversion
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
