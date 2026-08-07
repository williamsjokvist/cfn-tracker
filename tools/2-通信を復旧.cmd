@echo off
chcp 65001 >nul
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0net-rule.ps1" -Action Unblock
echo.
echo ---- 終了しました。何かキーを押すと閉じます ----
pause >nul
