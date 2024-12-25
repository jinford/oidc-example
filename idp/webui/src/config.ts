export const config = {
  apiUrl: import.meta.env.VITE_API_URL || '', // 開発環境では空文字列にして、Viteのプロキシ機能を使ってAPIサーバーにリクエストを転送する
};
