import { defineConfig } from 'astro/config';
import tailwind from '@astrojs/tailwind';
import mdx from '@astrojs/mdx';

// https://astro.build/config
export default defineConfig({
  site: 'https://cfn.williamsjokvist.se/',
  base: '/',
  integrations: [tailwind(), mdx()]
});