[Setup]
AppName=Go2Web
AppVersion=1.0.0
DefaultDirName={autopf}\Go2Web
; Outputs to the Output/ folder at the root of the repo
OutputDir=../Output
OutputBaseFilename=go2web
Compression=lzma
SolidCompression=yes
ArchitecturesInstallIn64BitMode=x64
; Tells Windows to refresh the PATH variable immediately
ChangesEnvironment=yes

[Files]
; Looks for the executable in the same scripts/ directory
Source: "go2web.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{autoprograms}\Go2Web"; Filename: "{app}\go2web.exe"

[Registry]
; Appends the installation directory to the user's PATH variable
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}"; Check: NeedsAddPath(ExpandConstant('{app}'))

[Code]
function NeedsAddPath(Param: string): boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKEY_CURRENT_USER, 'Environment', 'Path', OrigPath) then
  begin
    Result := True;
    exit;
  end;
  Result := Pos(';' + UpperCase(Param) + ';', ';' + UpperCase(OrigPath) + ';') = 0;
end;