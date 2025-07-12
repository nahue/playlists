// @ts-check
import { defineConfig, envField } from 'astro/config';

import alpinejs from '@astrojs/alpinejs';

import tailwindcss from '@tailwindcss/vite';

// https://astro.build/config
export default defineConfig({
  integrations: [alpinejs()],

  vite: {
    plugins: [tailwindcss()]
  },
  env: {
    schema: {
      PUBLIC_API_URL: envField.string({default: "http://localhost:8080", context: "client", access: "public"}),
    }
  }
});