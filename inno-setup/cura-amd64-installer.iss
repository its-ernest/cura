; --- Cura v0.1.9 ---

[Setup]
AppId={{D2445A8C-ED2B-4273-A3C2-2EA74F5916D7}}
AppName=Cura System Utility
AppVersion=0.1.9
AppPublisher=its-ernest@github.com
DefaultDirName={autopf}\Cura
DefaultGroupName=Cura
AllowNoIcons=yes
; Save the installer in the current directory
OutputDir=.
OutputBaseFilename=cura-amd64-installer
Compression=lzma
SolidCompression=yes
WizardStyle=modern

; forces the installer to run as Admin
PrivilegesRequired=admin
; ensures the app is installed for all users
PrivilegesRequiredOverridesAllowed=dialog

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked
Name: "startup"; Description: "Run Cura automatically on system boot"; GroupDescription: "Windows Integration:"

[Files]
; IMPORTANT: Replace 'amd64_folder' with actal architecture folder
; amd64 as the default for the installer bundle
Source: "amd64_folder\*"; DestDir: "{app}"; Flags: ignoreversion recursesubdirs createallsubdirs

[Icons]
Name: "{group}\Cura"; Filename: "{app}\cura-amd64.exe"
Name: "{autodesktop}\Cura"; Filename: "{app}\cura-amd64.exe"; Tasks: desktopicon

[Registry]
; This handles the 'Auto-start on boot' task
Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\Run"; \
    ValueType: string; ValueName: "CuraEnforcer"; \
    ValueData: """{app}\launcher-amd64.exe"""; \
    Tasks: startup; Flags: uninsdeletevalue

[Run]
; Option to launch Cura immediately after installation
Filename: "{app}\launcher-amd64.exe"; Description: "{cm:LaunchProgram,Cura}"; Flags: nowait postinstall skipifsilent