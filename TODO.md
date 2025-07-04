# TODO - 開発ロードマップ

## 🎯 プロダクトゴールと戦略

**最終目標**: 実用的な画像圧縮ツールとして広く採用される  
**戦略**: 既存実装活用 → 品質基盤 → 新価値 → UX向上 → 完成品質

## 📋 マイルストーン別開発計画

### 🎯 Milestone 1: PNG対応完了 (v0.2.0) - 基本機能の完成
*既存PNG実装を活用した最速価値提供*

- ✅ **m1_1**: PNG圧縮をCLIに統合（TDD Phase 1-3） - **完了**
- ✅ **m1_2**: PNG圧縮の統合テスト追加（CLI動作確認） - **完了**

**期待成果**: 実用的な画像圧縮ツールとして機能（JPEG+PNG対応）

### 🎯 Milestone 2: 品質基盤構築 (v0.3.0前) - 開発効率と信頼性
*長期開発効率と信頼性の確保*

- [ ] **m2_1**: GitHub Actions CI/CD設定（テスト自動化）
- [ ] **m2_2**: GoReleaser設定（バイナリ自動配布）
- [ ] **m2_3**: PNG圧縮のユニットテスト強化

**期待成果**: 安定した開発サイクル確立

### 🎯 Milestone 3: WebP対応 (v0.3.0) - 新価値提供
*WebP対応で差別化*

- [ ] **m3_1**: WebP圧縮エンジン実装（golang.org/x/image/webp使用）
- [ ] **m3_2**: WebPコンプレッサー登録とCLI統合
- [ ] **m3_3**: WebP圧縮テスト実装

**期待成果**: モダンな画像形式対応ツール（主要3形式対応）

### 🎯 Milestone 4: UX大幅改善 (v0.4.0) - 実用性向上
*大幅なUX改善*

- [ ] **m4_1**: 複数ファイル一括圧縮（ディレクトリ指定対応）
- [ ] **m4_2**: プログレスバー実装（大容量ファイル処理UX）
- [ ] **m4_3**: 設定ファイル対応（.shuku.yaml）

**期待成果**: 企業・プロジェクト採用レベルの実用性

### 🎯 Milestone 5: Production Ready (v1.0.0) - 完成品質
*完成品質の提供*

- [ ] **m5_1**: 包括的ドキュメント整備
- [ ] **m5_2**: ベンチマークテスト・パフォーマンス最適化
- [ ] **m5_3**: 詳細ログ機能実装

**期待成果**: エンタープライズ品質

## ⚡ 効率化戦略

### 🔄 並行作業戦略
- PNG統合中にCI/CD設定を並行実施
- 国際化は隙間時間で対応
- ドキュメント整備は機能開発と並行

### 📦 既存資産最大活用
- PNG内部実装済み→CLI統合のみ
- JPEG実装パターン→WebP実装に応用
- テスト基盤→新機能テストに流用

### 🎪 段階的価値提供
- 各マイルストーンでユーザー価値を提供
- 早期フィードバック獲得で品質向上

## 🚧 現在の作業状況

### 進行中タスク: m1_1 - PNG圧縮CLI統合（TDD開発）

**TDD開発フロー**:
- ✅ **Phase 1 (RED)**: 失敗するテストを作成 - 完了
  - PNG用テスト画像生成
  - PNG圧縮の失敗テスト作成
  - 期待通りの失敗を確認済み
- ✅ **Phase 2 (GREEN)**: 最小限実装でテスト通過 - **完了**
  - ✅ compress.goの形式制限削除 - 完了
  - ✅ PNG拡張子の許可追加 - 完了
  - ✅ PNGコンプレッサー実装と登録 - 完了
  - ✅ 全テスト通過確認 - 完了（UnsupportedFormatテスト修正済み）
- ✅ **Phase 3 (REFACTOR)**: コード品質向上 - **完了**
  - ✅ サポート形式リストの外部化 - 完了
  - ✅ 形式判定ロジックの関数化 - 完了

**次の具体的アクション**: Milestone 1完了！次はMilestone 2開始 - 品質基盤構築

## 🚀 次のアクション

Milestone 2: 品質基盤構築 (v0.3.0前)  
→ m2_1: GitHub Actions CI/CD設定（テスト自動化）

## 📅 更新履歴

- 2025-07-03: Milestone 1完了！ - m1_2完了、PNG圧縮統合テスト追加、次はMilestone 2
- 2025-07-03: m1_1完了 - TDD Phase 3(REFACTOR)完了、コード品質向上、次はm1_2
- 2025-07-03: Phase 2完了 - 全テスト通過確認、UnsupportedFormatテスト修正、Phase 3準備完了
- 2025-07-03: Phase 2実装開始 - compress.goの制限削除完了、PNGコンプレッサー登録作業中
- 2025-07-01: m1_1 TDD Phase 1完了 - PNGテストの失敗確認済み
- 2025-07-01: TODO.mdを作成、開発ロードマップをCLAUDE.mdから移行
- 2025-07-01: TDD開発フロー管理セクションを追加、現在の作業状況を明確化
