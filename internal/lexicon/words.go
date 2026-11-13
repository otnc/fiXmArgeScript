package lexicon

// words は fxma が使う辞書の語の一覧。綴りは Asankilp/fiXmArge-dictionary の
// fixmarge-EN.csv に従う。CSV はライセンスがないので同梱せず、使う語だけをここに書く。
var words = []Word{
	// 文と制御
	{"Iz'u", "〜である", Keyword | Part},
	{"showence", "証拠", Keyword},
	{"ReDogh'", "伝える", Keyword | Key},
	{"CaNoSiri", "罵り、叫び", Keyword | Key},
	{"zanfession", "懺悔", Keyword | Key},
	{"Mof", "もし", Keyword},
	{"SiVa", "しかし", Keyword},
	{"aZaGara", "〜しながら", Keyword},
	{"aSte", "〜として", Keyword},
	{"Kei-P", "継続", Keyword | Installer},
	{"KetZen", "結末", Keyword},
	{"FineNd", "終わり", Keyword | Console | Installer},
	{"We-Ho", "方法", Keyword},
	{"itsWA", "以上だ", Keyword},
	{"AbiRyoku", "能力", Keyword},

	// 世界観由来の構文
	{"meKaMo", "〜かも", Keyword},
	{"ShouBeK", "〜するべき", Keyword},
	{"sent-ky", "裁き", Keyword},
	{"AccUhEarts", "告発", Keyword | Error},
	{"Tei-Tance", "抵抗", Keyword},
	{"yOuLugiVe", "許し", Keyword | ModField},
	{"En-NOte", "成れの果て", Keyword},

	// 値と演算
	{"TadACt", "正しさ", Keyword},
	{"ULieSo", "嘘", Keyword},
	{"Vani·tA", "虚無", Keyword},
	{"N'o", "の", Keyword},
	{"eXi'", "存在する", Keyword | Part},
	{"Nii", "否定", Keyword | Console},
	{"S~ond'", "そして", Keyword},
	{"Owk", "または", Keyword},
	{"SeNaji", "同じ", Keyword | Part},
	{"Batte", "〜によって、〜を通して", Keyword | Part},

	// 非同期・マクロ・パッケージ
	{"itsver", "いつか", Keyword | Installer},
	{"pray·lee", "祈る", Keyword},
	{"cUrOi", "呪い", Keyword | Option},
	{"Oqun", "開く", Keyword},
	{"weRer", "迎える、歓迎", Keyword | Command | Installer},

	// 組み込み
	{"nAng", "長い", Builtin | Key},
	{"sAnAnd", "砂", Builtin | Part},

	// 標準パッケージ
	{"MeMe", "数学", Package},
	{"Coword", "言葉", Package | Key},
	{"VariZaMa", "様々", Package},
	{"Toim", "時間", Package | Key},
	{"dEstiMay", "運命", Package},
	{"KeyGa", "鍵", Package},
	{"ievil JIO", "魔女", Package | Key},
	{"GE-mon", "門", Package},
	{"Memoid", "記憶", Package},
	{"i·zana-vited", "導かれた", Package | Key | Part},

	// 標準パッケージの関数
	{"RacDA-m", "落ちる", Function},
	{"Hir", "空を飛ぶもの", Function},
	{"Litto", "ちょうど", Function},
	{"PyusenS'e", "純粋な", Function},
	{"mOurce", "元、根源", Function},
	{"tSuYoNg", "強い", Function},
	{"YowEAk", "弱い", Function},
	{"Zelln", "全員", Function},
	{"SaLook", "姿、形", Function},
	{"conf-Soy", "争い", Function},
	{"rEcanRey", "関係", Function},
	{"Lai Nei", "〜のように", Function | ModField | Part},
	{"ShoJikA", "短い", Function},
	{"gDie", "偉大な", Function},
	{"Ligh-Luw", "軽い", Function},
	{"Eda", "そこ", Function},
	{"Kast", "過去", Function | Part},
	{"gu-n", "今", Function | Part},
	{"KolSu", "音", Function | Part},
	{"nAddY", "仲間", Function},
	{"mu·Yon", "向こう側", Function},
	{"oRei", "命令、順序", Function},
	{"Tarip", "旅", Function},
	{"oNkhe", "〜だけ", Function},
	{"Vani·Ell", "消える", Function},
	{"LEA", "さあ、来い", Function},
	{"Matill", "まだ", Function},
	{"ieyWOO", "今すぐ", Function},
	{"Divil JIO", "大魔女", Function},
	{"sAiMoku", "沈黙", Function | LintRule},
	{"FinecloWt", "見つけ出す", Function | Part},
	{"Sen-choose", "選択肢", Function},
	{"Commity", "たぶん", Function},
	{"madRazy", "狂った", Function | Part},
	{"GoMi", "神", Function},
	{"SurC", "生き残る", Function},
	{"sAgAin", "また", Function},
	{"GoKuu", "行く", Function | Installer},
	{"Edo", "ここ", Function | Key},
	{"nOzuki", "気付き", Function | Key},
	{"deFuTai", "絶対に", Function | Part},
	{"KuwPace", "空間", Function | Part | Installer},
	{"Ado", "これ", Function | Command | ModField},
	{"ChiDd", "血", Function | Key},
	{"eATah·Guy", "お互い", Function},
	{"Say mAl", "正常", Function | Command},

	// 2語の名前の部品だけに使う語
	{"Soda", "その程度まで", Part},
	{"Alte", "全て", Part | Key},
	{"fOre", "正確な", Part},
	{"MEif", "命", Part | Key | Command},

	// 辞書のキー
	{"LyuzUn", "理由", Key},
	{"muRm", "殺人", Key},
	{"ConWa", "会話", Key},
	{"LowhEll", "牢屋", Key},
	{"shid", "死んだ", Key},
	{"WiD", "世界", Key},
	{"nEYar", "年", Key | Command},
	{"hAte", "手", Key},
	{"eMoJoh", "感情", Key},
	{"FayO", "顔", Key},

	// エラーと警告
	{"AyAmis", "過ち", Error},
	{"imJune", "不純", Error},
	{"kYoNAihim", "虚無", Error},
	{"Bavellabion", "呪いの言葉があふれる", Error},
	{"blasTOac", "冒涜的", Error},
	{"kiNbOO", "禁忌", Error},
	{"Co-Jundic", "矛盾", Error},
	{"KeI·tion", "処刑", Error},
	{"Su·Ready", "既に", Error},
	{"CraeyTA'", "壊れた", Error},
	{"SaiVI", "災い", Error},
	{"Ca-nay", "できない", Error},
	{"DowGuil", "疑惑", Error | LintRule},
	{"LaVI", "法の悪魔", Error | Command},

	// fixmarge lavi の規則
	{"Was·Da", "無駄", LintRule},
	{"shi-des", "死", LintRule},
	{"emBasy", "恥ずかしい", LintRule},
	{"anxAn", "不安", LintRule},
	{"SuZai", "怪しい", LintRule},
	{"Meh·BeRac", "たぶん", LintRule},
	{"JokeDan", "冗談", LintRule},
	{"dEs·lee", "絶望", LintRule},

	// コンソールの文言
	{"hung gement", "判決", Console},
	{"Guil-me", "罪", Console},

	// コマンドとオプション
	{"kee’on", "続ける", Command},
	{"sinruits", "現実", Command},
	{"KaRaDy", "体", Option},
	{"Coa Tain", "法廷", Command},
	{"beAukShi", "美しい", Command},
	{"hA-k", "話す", Command},
	{"NAi", "なぜ", Command},
	{"Per·SemOn", "見世物", Keyword | Command | Console},
	{"Trench", "友達", Command},

	// fixmarge.mod の項目
	{"MokpOSe", "目的", ModField | Option},
	{"neYO", "必要", ModField},
	{"UnneYO", "要らない", ModField},
	{"forYulvu", "許可する", ModField},

	// インストーラの文言
	{"ShiFey", "失敗", Installer},
	{"Syu Tain", "休廷", Installer},
}

// compounds は、JavaScript で複数語の名前を辞書の語の camelCase で組んだもの (RFC の表記の規則 6)。
var compounds = []string{
	"laineiAlte",
	"sodaKast",
	"sodaGun",
	"batteKolsu",
	"foreSenaji",
	"madrazySanand",
	"izanavitedExi",
	"fineclowtKuwpace",
	"meifKuwpace",
	"izuDefutai",
	"izuKuwpace",
}

// dialects は、標準語と違う綴りの骨格から標準語の骨格への対応。
var dialects = map[string]string{
	"lowhel": "lowhell",
}
