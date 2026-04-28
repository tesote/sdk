import { defineConfig } from 'tsup';

export default defineConfig({
  entry: {
    index: 'src/index.ts',
    'v1/index': 'src/v1/index.ts',
    'v2/index': 'src/v2/index.ts',
  },
  format: ['esm', 'cjs'],
  target: 'es2020',
  outDir: 'dist',
  dts: true,
  sourcemap: true,
  clean: true,
  minify: false,
  splitting: false,
  treeshake: true,
  outExtension({ format }) {
    return { js: format === 'cjs' ? '.cjs' : '.js' };
  },
});
