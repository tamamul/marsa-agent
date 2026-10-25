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
Source: "dist\tray.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "dist\SumatraPDF.exe"; DestDir: "{app}"; Flags: ignoreversion
; config.json cuma disalin kalau belum ada, supaya re-install tidak menimpa
; setting yang sudah diisi admin sekolah
Source: "dist\config.json"; DestDir: "{app}"; Flags: onlyifdoesntexist

[Icons]
Name: "{group}\Buka Dashboard Agent"; Filename: "http://127.0.0.1:9191"
Name: "{group}\Uninstall MARSA Print Agent"; Filename: "{uninstallexe}"
; Tray icon ikut nyala otomatis tiap kali ada yang login Windows
Name: "{userstartup}\MARSA Print Agent Tray"; Filename: "{app}\tray.exe"

[Run]
; Daftarkan & jalankan sebagai Windows Service otomatis setelah install
Filename: "{app}\agent.exe"; Parameters: "install"; Flags: runhidden waituntilterminated
Filename: "{app}\agent.exe"; Parameters: "start"; Flags: runhidden waituntilterminated
; Jalankan tray app juga saat ini (tidak perlu tunggu login ulang)
Filename: "{app}\tray.exe"; Flags: nowait postinstall skipifsilent runasoriginaluser; Description: "Jalankan ikon tray sekarang"
; Langsung buka dashboard setting di browser biar admin bisa isi URL server + token
Filename: "http://127.0.0.1:9191"; Description: "Buka dashboard untuk isi setting"; Flags: postinstall shellexec skipifsilent

[UninstallRun]
; Hentikan tray app, lalu hentikan & hapus service saat di-uninstall
Filename: "{cmd}"; Parameters: "/C taskkill /IM tray.exe /F"; Flags: runhidden; RunOnceId: "KillTray"
Filename: "{app}\agent.exe"; Parameters: "stop"; Flags: runhidden waituntilterminated; RunOnceId: "StopAgent"
Filename: "{app}\agent.exe"; Parameters: "uninstall"; Flags: runhidden waituntilterminated; RunOnceId: "UninstallAgent"

[UninstallDelete]
Type: files; Name: "{userstartup}\MARSA Print Agent Tray.lnk"
