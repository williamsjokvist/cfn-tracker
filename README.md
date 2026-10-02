<h1 align="center">CFN Tracker</h1>

<p align="center">
  Street Fighter 6 / Tekken 8 の対戦成績をリアルタイムで取得し、OBS に表示するツール<br>
  <sub>Realtime tracking of ranked matches in Street Fighter 6 and Tekken 8</sub>
</p>

---

## これは何か

[williamsjokvist/cfn-tracker](https://github.com/williamsjokvist/cfn-tracker) の**非公式フォーク**です。
本体の機能はそのままに、**配信中に監視が黙って止まってしまう不具合**を実機で特定して修正しています。

> [!IMPORTANT]
> 本家の作者とは無関係の有志フォークです。**不具合の報告は本家ではなく、このリポジトリへ**お願いします。
> 本家をお使いの方で「気づいたら成績が更新されていなかった」という経験がある方に向けた版です。

## 本家との違い

実際の配信環境で起きた障害を、ログ解析で原因を特定して直したものです。

| 起きていたこと | 本家 | このフォーク |
|---|---|---|
| 回線が一瞬切れると、以後ずっと更新されない | 止まったまま戻らない | 自動で再接続を繰り返す。回数制限なし |
| 開始から 30 分ほどで、何も言わず監視が止まる | 止まる（ログにも何も残らない） | タイムアウトで検知し、自動で再試行する |
| ログイン画面がループして先へ進めない | 進めない | ブラウザが開き、自分でログインできる |
| 不具合が起きてもログに原因が残らない | 握りつぶされる | 原因が記録される |

技術的な詳細は各コミットのメッセージに日本語で書いてあります。

## ダウンロード

**[最新版をダウンロード](https://github.com/layla-amagiri/cfn-tracker/releases/latest)**

`cfn-tracker-windows-amd64.zip` を展開してください。Windows 64bit 版のみです。インストール作業は要りません。

## 使い方

### 1. 設定ファイルを作る

展開したフォルダにある `example.env` を、同じ場所に `.env` という名前でコピーします。
中身をメモ帳などで開き、自分の Capcom ID の情報に書き換えてください。

```
CAP_ID_EMAIL="あなたのメールアドレス"
CAP_ID_PASSWORD="あなたのパスワード"

HEADLESS=true
BROWSER_SOURCE_PORT=4242
```

> [!WARNING]
> **パスワードはこのファイルに平文で保存されます。** 配信画面にこのファイルを映さないでください。
> 設定ファイルごと他人に渡さないでください。この仕組みは本家からそのまま引き継いでいます。

### 2. 起動する

`CFN Tracker.exe` を実行し、ゲームと自分の CFN ID を選んで開始します。

初回、および**月に一度ほど**、ログインを求められます。その場合はブラウザが自動で開くので、
画面の案内に従って Capcom ID でログインし、**ログインできたらそのブラウザを閉じてください**。
閉じた時点でアプリが自動的に続きを進めます。

### 3. OBS に表示する

OBS の「ブラウザ」ソースに、アプリ内の Output ページに表示される URL を設定します。
テキストファイル出力にも対応しているので、テキストソースでも表示できます。
表示の見た目は CSS で自由に変更できます。

## 注意事項

- **自己責任でお使いください。** 本ツールは Capcom の公式ツールではありません。
- ウイルス対策ソフトが誤検知することがあります。署名のない個人製のアプリのため、この種の警告が出ることがあります。気になる場合はソースから自分でビルドしてください。
- 対戦データの取得にブラウザを内部で使うため、起動中は Chrome のプロセスが常駐します。
- Tekken 8 の対戦取得は本家の実装のままで、このフォークでは手を入れていません。

## 開発者向け

### 必要なもの

- [Wails](https://wails.io/docs/gettingstarted/installation)
- [Bun](https://bun.sh)
- [Task](https://taskfile.dev) — 任意

環境変数は `app/example.env` を参照してください。`task --list` でコマンド一覧が出ます。

### ビルド

```sh
cd app
wails build -ldflags "-X main.isProduction=true"
```

> [!NOTE]
> `Taskfile.yml` 経由のビルドは、ldflags の値に引用符が含まれる形になっており
> `isProduction` の判定が外れます。上のコマンドを直接実行してください。

## ライセンス

本家と同じく [CC0 1.0 Universal + Commons Clause](LICENSE) です。改変・再配布は自由ですが、販売は禁止されています。

## クレジット

- 本家: [williamsjokvist/cfn-tracker](https://github.com/williamsjokvist/cfn-tracker) — 元となるアプリのすべての機能は本家の作者によるものです。

---

## English

This is an **unofficial fork** of [williamsjokvist/cfn-tracker](https://github.com/williamsjokvist/cfn-tracker).

The upstream app stops tracking silently under several conditions. This fork fixes them, each one diagnosed from real logs on a live streaming setup:

- **Network drops are no longer fatal.** Transient errors are classified and retried with exponential backoff, with no retry cap.
- **Polling can no longer hang forever.** `GetBattleLog` had no timeout, so a stalled element wait blocked the poll loop indefinitely — tracking died silently after roughly 30 minutes. It is now guarded by a 25s inner and 40s outer timeout.
- **Login no longer loops on Cloudflare Turnstile.** A browser driven over the DevTools Protocol cannot pass Turnstile. The fork stops fighting it: it closes the controlled browser, opens a plain Chrome with the same profile, lets you log in by hand, and picks the session back up when you close the window.
- **Errors reach the log.** Several failures were swallowed, including a panic in the error formatter that destroyed the original error.

Commit messages are written in Japanese. Downloads are on the [releases page](https://github.com/layla-amagiri/cfn-tracker/releases/latest); Windows x64 only. Copy `example.env` to `.env` and fill in your Capcom ID credentials before running. Note that credentials are stored in plain text, which is inherited from upstream.

Bug reports belong here, not upstream.
