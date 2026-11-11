// ロゴとアイコンを作り、使う場所へ書き出す。図形はすべてここで描き、公式の素材は使わない。
// 文字は Almendra Display (SIL Open Font License 1.1) の字形を線に変えて埋め込むので、見る側にフォントは要らない。
import { execFileSync } from "node:child_process";
import { mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { Resvg } from "@resvg/resvg-js";
import opentype from "opentype.js";
import pngToIco from "png-to-ico";

const root = join(import.meta.dirname, "..");

const color = {
  night: "#0d0a0c",
  ember: "#3a0712",
  crimson: "#b0102a",
  blood: "#6e0818",
  bone: "#f2ecec",
  ink: "#1a1215",
};

// 格子は柵のように先を尖らせる。蝶が前に重なるので、真ん中の2本は蝶の羽に隠れる位置に置く。
// 小さく表示する版 (32px 以下) は、つぶれないように格子を太くし、蝶を大きくして翅脈と斑点を省く。
function bars(small) {
  const xs = small ? [92, 204, 308, 420] : [118, 210, 302, 394];
  const width = small ? 36 : 24;
  const top = 92;
  const bottom = 452;
  const spikes = xs
    .map((x) => {
      const l = x - width / 2;
      const r = x + width / 2;
      return `M${l} ${bottom}V${top + 26}L${x} ${top}L${r} ${top + 26}V${bottom}Z`;
    })
    .join("");
  const rails = [168, 402].map((y) => `M70 ${y}H442V${y + 16}H70Z`).join("");
  return `<path d="${spikes}${rails}" fill="${color.bone}" opacity="0.88"/>`;
}

// 蝶は右半分を描いて左右に写す。羽の先を尖らせ、下の羽に尾を付けて、葉に見えないようにする。
// 羽は暗い色で描いてから、付け根を中心に縮めた同じ形を深紅で重ねる。外側ほど太い縁取りになる。
function butterfly(small) {
  const upper = "M6 -6C22 -60 70 -122 150 -140L196 -130C188 -96 172 -60 150 -40C128 -20 98 -6 60 0C36 3 16 2 6 -6Z";
  const lower = "M6 4C44 8 100 24 112 60C118 80 110 96 92 104L100 158L76 108C56 104 34 82 22 56C14 38 8 20 6 4Z";
  const veins = "M8 -4L150 -126M8 -2L172 -96M8 0L130 -30M8 4L96 96M8 6L104 62";
  const spots = [[176, -120, 8], [184, -100, 6], [178, -80, 5], [160, -128, 6]]
    .map(([x, y, r]) => `M${x} ${y - r}L${x + r} ${y}L${x} ${y + r}L${x - r} ${y}Z`)
    .join("");
  const half = `
    <path d="${upper}${lower}" fill="${color.blood}" stroke="${color.night}" stroke-width="10" stroke-linejoin="round" paint-order="stroke"/>
    <path d="${upper}${lower}" fill="${color.crimson}" transform="scale(0.8)"/>
    ${small ? "" : `<path d="${veins}" fill="none" stroke="${color.night}" stroke-width="3" stroke-linecap="round" opacity="0.55"/>
    <path d="${spots}" fill="${color.bone}"/>`}`;
  const body = `
    <path d="M0 -40C10 -40 12 -20 10 10C8 40 4 70 0 84C-4 70 -8 40 -10 10C-12 -20 -10 -40 0 -40Z" fill="${color.bone}" stroke="${color.night}" stroke-width="6" paint-order="stroke"/>
    <path d="M-4 -38C-14 -70 -30 -92 -50 -102M4 -38C14 -70 30 -92 50 -102" fill="none" stroke="${color.bone}" stroke-width="5" stroke-linecap="round"/>
    <path d="M-50 -110L-42 -102L-50 -94L-58 -102ZM50 -110L58 -102L50 -94L42 -102Z" fill="${color.bone}"/>`;
  return `<g transform="translate(256 ${small ? 270 : 286}) scale(${small ? 1.1 : 1})">
    <g>${half}</g>
    <g transform="scale(-1 1)">${half}</g>
    ${body}
  </g>`;
}

function iconBody(small = false) {
  return `<defs>
    <radialGradient id="glow" cx="50%" cy="56%" r="50%">
      <stop offset="0" stop-color="${color.ember}"/>
      <stop offset="1" stop-color="${color.night}"/>
    </radialGradient>
  </defs>
  <rect width="512" height="512" rx="96" fill="url(#glow)"/>
  ${bars(small)}
  ${butterfly(small)}`;
}

function icon(small = false) {
  return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
  ${iconBody(small)}
</svg>
`;
}

const font = (() => {
  const file = join(root, "assets/node_modules/@fontsource/almendra-display/files/almendra-display-latin-400-normal.woff");
  const buf = readFileSync(file);
  return opentype.parse(buf.buffer.slice(buf.byteOffset, buf.byteOffset + buf.byteLength));
})();

// 大文字は清書の綴りの目印なので深紅にする。合字にすると字と字形が1対1にならないので切る。
function wordmark(text, x, y, size, ink) {
  const options = { features: { liga: false, rlig: false } };
  const paths = font.getPaths(text, x, y, size, options);
  if (paths.length !== [...text].length) {
    throw new Error(`字形の数が文字の数と合わない: ${paths.length} / ${text.length}`);
  }
  const groups = { [ink]: "", [color.crimson]: "" };
  [...text].forEach((ch, i) => {
    const fill = /[A-Z]/.test(ch) && ch !== "S" ? color.crimson : ink;
    groups[fill] += paths[i].toPathData(2);
  });
  const advance = font.getAdvanceWidth(text, size, options);
  // 細い字形なので、同じ色の線でなぞって少し太らせる (小さく表示しても消えないように)。
  const svg = Object.entries(groups)
    .map(([fill, d]) => `<path d="${d}" fill="${fill}" stroke="${fill}" stroke-width="1.6" stroke-linejoin="round"/>`)
    .join("\n  ");
  return { svg, advance };
}

// 文字の下に、ひし形で区切った細い線を引く。
function rule(x1, x2, y, ink) {
  const mid = (x1 + x2) / 2;
  const diamond = (cx, r) => `M${cx} ${y - r}L${cx + r} ${y}L${cx} ${y + r}L${cx - r} ${y}Z`;
  return `<path d="M${x1 + 12} ${y}H${mid - 18}M${mid + 18} ${y}H${x2 - 12}" stroke="${ink}" stroke-width="2"/>
  <path d="${diamond(x1 + 4, 6)}${diamond(x2 - 4, 6)}" fill="${ink}"/>
  <path d="${diamond(mid, 11)}" fill="${color.crimson}"/>`;
}

function logo(ink) {
  const iconSize = 200;
  const gap = 36;
  const textX = iconSize + gap;
  const { svg: text, advance } = wordmark("fiXmArge Script", textX, 128, 112, ink);
  const width = Math.ceil(textX + advance + 8);
  return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${width} ${iconSize}" role="img" aria-label="fiXmArge Script">
  <g transform="scale(${iconSize / 512})">
  ${iconBody()}
  </g>
  ${text}
  ${rule(textX, textX + advance, 170, ink)}
</svg>
`;
}

// ファイルのアイコン。エクスプローラーでは 16px なので小さい版を使い、パッケージとロックのファイルは右下に印を大きく重ねる。
// 印は白で塗り、黒い縁で囲んで、赤い蝶の上でも形が分かるようにする。
const badges = {
  // パッケージのファイル (fixmarge.mod) は箱。
  box: `<path d="M0 -62L56 -32V34L0 64L-56 34V-32Z" fill="none" stroke="${color.night}" stroke-width="24" stroke-linejoin="round"/>
    <path d="M0 -62L56 -32L0 -2L-56 -32Z" fill="${color.bone}"/>
    <path d="M-56 -32L0 -2V64L-56 34Z" fill="#c9bfc2"/>
    <path d="M56 -32L0 -2V64L56 34Z" fill="#978b90"/>
    <path d="M-56 -32L0 -2L56 -32M0 -2V64" fill="none" stroke="${color.night}" stroke-width="6" stroke-linejoin="round"/>`,
  // ロックファイル (fixmarge.sum) は鍵。斜めにすると先がタイルの丸い角からはみ出すので横に寝かせる。
  key: `<g fill="none" stroke-linecap="round" transform="translate(-12 0) scale(0.83)">
      <g stroke="${color.night}" stroke-width="36"><circle cx="-50" cy="0" r="30"/><path d="M-28 0H76M46 0V28M70 0V22"/></g>
      <g stroke="${color.bone}" stroke-width="18"><circle cx="-50" cy="0" r="30"/><path d="M-28 0H76M46 0V28M70 0V22"/></g>
    </g>`,
};

function fileIcon(badge) {
  const mark = badge ? `<g transform="translate(384 384) scale(1.75)">${badges[badge]}</g>` : "";
  return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
  ${iconBody(true)}
  ${mark}
</svg>
`;
}

// VS Code の補助のアイコンテーマ。ほかのファイルには何も付けないので、fiXmArge Script のファイルだけを見分けたいときに選ぶ。
function iconTheme(files) {
  const theme = { iconDefinitions: {}, languageIds: {}, fileExtensions: {}, fileNames: {} };
  for (const f of files) {
    theme.iconDefinitions[f.id] = { iconPath: `./${f.id}.svg` };
    theme.languageIds[f.language] = f.id;
    for (const ext of f.extensions ?? []) theme.fileExtensions[ext] = f.id;
    for (const name of f.fileNames ?? []) theme.fileNames[name] = f.id;
  }
  return JSON.stringify(theme, null, 2) + "\n";
}

function png(svg, size) {
  return new Resvg(svg, { fitTo: { mode: "width", value: size } }).render().asPng();
}

function write(path, data) {
  const full = join(root, path);
  mkdirSync(dirname(full), { recursive: true });
  writeFileSync(full, data);
  console.log(path);
}

const iconSvg = icon();
const smallIconSvg = icon(true);
write("assets/icon.svg", iconSvg);
write("assets/icon-small.svg", smallIconSvg);
write("assets/logo-dark.svg", logo(color.bone));
write("assets/logo-light.svg", logo(color.ink));
write("editors/vscode/icon.png", png(iconSvg, 256));
write("playground/src/favicon.svg", smallIconSvg);
const sizes = [16, 24, 32, 48, 64, 128, 256];
write("installer/fixmarge.ico", await pngToIco(sizes.map((size) => png(size <= 32 ? smallIconSvg : iconSvg, size))));

// fixmarge.exe に埋め込むアイコン。go build が cmd/fixmarge の *_windows_<arch>.syso を拾う。
// 版は入れない (fixmarge neyar は Go が埋め込む版を読むので、ここで入れるとリリースのたびに作り直すことになる)。
// go-winres は版がないと製品名などの文字の情報も書かないので、アイコンだけを入れる。
// マニフェストも付けない。付けると Windows での動き方が変わりうるため。
execFileSync(
  "go",
  [
    "run",
    "github.com/tc-hib/go-winres@v0.3.3",
    "simply",
    "--arch=amd64,arm64",
    "--out=cmd/fixmarge/rsrc",
    "--manifest=none",
    "--icon=installer/fixmarge.ico",
  ],
  { cwd: root, stdio: "inherit" },
);
console.log("cmd/fixmarge/rsrc_windows_{amd64,arm64}.syso");

// 言語の ID とファイルの名前は editors/vscode/package.json の languages と合わせる。
const files = [
  { id: "fixmarge", language: "fixmarge", extensions: ["fxma"] },
  { id: "fixmarge", language: "fixmargec", extensions: ["fxmac"] },
  { id: "fixmarge-mod", language: "fixmarge-mod", fileNames: ["fixmarge.mod"], badge: "box" },
  { id: "fixmarge-sum", language: "fixmarge-sum", fileNames: ["fixmarge.sum"], badge: "key" },
];
for (const f of new Map(files.map((f) => [f.id, f])).values()) {
  write(`editors/vscode/icons/${f.id}.svg`, fileIcon(f.badge));
}
write("editors/vscode/icons/fixmarge-icon-theme.json", iconTheme(files));
