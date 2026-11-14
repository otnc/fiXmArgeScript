// install.sh はインストーラと一緒に installer/ で管理し、このサイトの /install.sh で配る (curl ... | sh)
import type { APIRoute } from "astro";

import script from "../../../installer/public/install.sh?raw";

export const GET: APIRoute = () => new Response(script, { headers: { "Content-Type": "text/plain; charset=utf-8" } });
