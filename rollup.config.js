import typescript from '@rollup/plugin-typescript';
import { nodeResolve } from '@rollup/plugin-node-resolve';
import terser from '@rollup/plugin-terser';
import CleanCSS from 'clean-css';
import commonjs from '@rollup/plugin-commonjs';
import copy from 'rollup-plugin-copy';

export default {
  input: 'client/main.ts',
  output: {
    name: 'webtimer',
    dir: 'dist/frontend',
    format: 'iife',
    sourcemap: true,
    globals: {
      '@textfit': 'textfit',
    },
    sourcemapPathTransform: (relativeSourcePath, sourcemapPath) => {
      return `${relativeSourcePath.replace(new RegExp('^\.\./'), '')}`;
    },
  },
  plugins: [
    typescript({
      "compilerOptions": {
        "outDir": "dist/frontend/build-tsc"
      },
    }),
    commonjs(),
    nodeResolve(),
    terser(),
    copy({
      targets: [
        { src: 'html/index.html', dest: 'dist/frontend' },
        { src: 'html/inter.woff2', dest: 'dist/frontend' },
        {
          src: 'html/style.css',
          dest: 'dist/frontend',
          transform: (contents) => new CleanCSS().minify(contents).styles
        },
      ]
    })
  ],
  external: ['@textfit']
};
