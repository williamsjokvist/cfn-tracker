#requires -Version 5.1
<#
検証用: cfn-tracker の通信を遮断／復旧する。

SF6 のポーリングは CFN Tracker.exe 本体ではなく rod が起動した専用 Chromium が行うため、
遮断対象はその Chromium である。普段使いの Google Chrome とは別物なので、
通常のブラウジングには影響しない。

管理者権限が要るので、非管理者で起動された場合は自分を昇格して起動し直す。
#>
param(
    [ValidateSet('Block', 'Unblock')]
    [string]$Action = 'Unblock'
)

$ruleName = 'BLOCK cfn-tracker (test)'

$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).
    IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Host '管理者権限が必要なため昇格します。UAC の確認が出たら「はい」を選んでください...' -ForegroundColor Cyan
    try {
        Start-Process powershell -Verb RunAs -ErrorAction Stop -ArgumentList @(
            '-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', "`"$PSCommandPath`"", '-Action', $Action
        )
        Write-Host '昇格したウィンドウを起動しました。そちらの表示を確認してください。' -ForegroundColor Cyan
    }
    catch {
        Write-Host '昇格に失敗しました: ' -ForegroundColor Red -NoNewline
        Write-Host $_.Exception.Message -ForegroundColor Red
        Write-Host ''
        Write-Host '手動で行う場合は、管理者権限の PowerShell を開いて次を実行してください:' -ForegroundColor Yellow
        Write-Host '  Remove-NetFirewallRule -DisplayName BLOCK*' -ForegroundColor Yellow
    }
    return
}

if ($Action -eq 'Block') {
    $rod = Get-Process chrome -ErrorAction SilentlyContinue |
           Where-Object { $_.Path -like '*\rod\browser\*' } |
           Select-Object -First 1 -ExpandProperty Path

    if (-not $rod) {
        Write-Host 'rod の Chromium が見つかりません。' -ForegroundColor Red
        Write-Host 'アプリを起動し、トラッキングを開始した状態で実行してください。' -ForegroundColor Red
        Read-Host 'Enter で閉じる'
        return
    }

    if (Get-NetFirewallRule -DisplayName $ruleName -ErrorAction SilentlyContinue) {
        Remove-NetFirewallRule -DisplayName $ruleName
    }
    New-NetFirewallRule -DisplayName $ruleName -Direction Outbound -Action Block -Program $rod | Out-Null

    Write-Host '通信を遮断しました。' -ForegroundColor Yellow
    Write-Host "対象: $rod"
    Write-Host '最大30秒（ポーリング周期）待つと再試行が始まります。'
}
else {
    if (Get-NetFirewallRule -DisplayName $ruleName -ErrorAction SilentlyContinue) {
        Remove-NetFirewallRule -DisplayName $ruleName
        Write-Host '通信を復旧しました。' -ForegroundColor Green
        Write-Host '次の再試行タイミングで自動的に追跡へ戻るはずです。'
    }
    else {
        Write-Host '規則は存在しません（既に復旧済み）。' -ForegroundColor Green
    }
}

Read-Host 'Enter で閉じる'
