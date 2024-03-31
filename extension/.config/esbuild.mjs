import * as esbuild from 'esbuild'

await esbuild.build({
  entryPoints: ['src/extension.js'],
  bundle: true,
  format: 'cjs',
  platform: 'node',
  external: ['vscode'],
  outfile: 'build/extension.js',
})