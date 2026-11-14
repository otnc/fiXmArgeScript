import { defineCollection } from "astro:content";
import { docsLoader } from "@astrojs/starlight/loaders";
import { docsSchema } from "@astrojs/starlight/schema";

// ドキュメントは src/content/docs/docs/ に置き、/docs/ の下に出す (/ は紹介のページ)
export const collections = {
  docs: defineCollection({ loader: docsLoader(), schema: docsSchema() }),
};
