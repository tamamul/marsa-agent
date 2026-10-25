; Script Inno Setup untuk MARSA Print Agent
; Cara pakai: install Inno Setup Compiler (gratis, sekali di laptop dev),
; buka file ini, klik Compile -> hasilnya MarsaPrintAgentSetup.exe
;
; Folder "dist" harus berisi: agent.exe, SumatraPDF.exe, config.json

#define MyAppName "MARSA Print Agent"
#define MyAppVersion "1.0"
#define MyAppPublisher "SMK Ma'arif 9 Kebumen"

[Setup]
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
DefaultDirName={autopf}\MarsaPrintAgent
DefaultGroupName=MARSA Print Agent
DisableProgramGroupPage=yes
OutputDir=output
OutputBaseFilename=MarsaPrintAgentSetup
Compression=lzma
SolidCompression=yes
PrivilegesRequired=admin
ArchitecturesInstallIn64BitMode=x64

[Files]
Source: "dist\agent.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "dist\SumatraPDF.exe"; DestDir: "{app}"; Flags: ignoreversion
; config.json cuma disalin kalau belum ada, supaya re-install tidak menimpa
; setting yang sudah diisi admin sekolah
Source: "dist\config.json"; DestDir: "{app}"; Flags: onlyifdoesntexist

[Icons]
Name: "{group}\Buka Dashboard Agent"; Filename: "http://127.0.0.1:9191"
Name: "{group}\Uninstall MARSA Print Agent"; Filename: "{uninstallexe}"

[Run]
; Daftarkan & jalankan sebagai Windows Service otomatis setelah install
Filename: "{app}\agent.exe"; Parameters: "install"; Flags: runhidden waituntilterminated
Filename: "{app}\agent.exe"; Parameters: "start"; Flags: runhidden waituntilterminated
; Langsung buka dashboard setting di browser biar admin bisa isi URL server + token
Filename: "http://127.0.0.1:9191"; Description: "Buka dashboard untuk isi setting"; Flags: postinstall shellexec skipifsilent

[UninstallRun]
; Hentikan & hapus service saat di-uninstall
Filename: "{app}\agent.exe"; Parameters: "stop"; Flags: runhidden waituntilterminated; RunOnceId: "StopAgent"
Filename: "{app}\agent.exe"; Parameters: "uninstall"; Flags: runhidden waituntilterminated; RunOnceId: "UninstallAgent"
