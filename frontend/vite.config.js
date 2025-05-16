import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
      plugins: [react()],
      server: {
            host: true,
            port: 3000,
            watch: {
                  usePolling: true
            }
      },
      build: {
            outDir: 'dist',
            emptyOutDir: true
      }
}); 