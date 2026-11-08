export interface RunOptions {
  /** zanfession で読む入力。ブラウザでは対話的に読めないので、まとめて渡す。 */
  stdin?: string;
  /** werer と fxma:memoid から見えるファイル。名前 → 中身。 */
  files?: Record<string, string>;
  /** dv no "conwa" になるコマンドライン引数。 */
  args?: string[];
}

export interface RunResult {
  stdout: string;
  stderr: string;
  /** 正常なら 0、エラーなら 1、dv no "shid" で終われば渡した番号。 */
  exitCode: number;
  /** 実行した後のファイル。名前 → 中身。 */
  files: Record<string, string>;
}

export interface BuildResult {
  /** 清書。エラーなら null。 */
  output: string | null;
  /** 方言の警告とエラー。 */
  stderr: string;
}

export interface Word {
  /** 下書きに書く形。 */
  name: string;
  /** 清書に書く形。 */
  spelling: string;
  uses: Array<"keyword" | "builtin" | "package" | "function" | "key" | "error">;
}

/** 下書きを実行する。 */
export function run(source: string, options?: RunOptions): Promise<RunResult>;

/** 下書きを清書にする。 */
export function build(source: string): Promise<BuildResult>;

/** 下書きに書ける名前の一覧。エディタの色分けに使う。 */
export function words(): Promise<Word[]>;
