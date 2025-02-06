import typescript from '@rollup/plugin-typescript';
import { nodeResolve } from '@rollup/plugin-node-resolve';
import terser from '@rollup/plugin-terser';
import CleanCSS from 'clean-css';
import commonjs from '@rollup/plugin-commonjs';
import concat from 'rollup-plugin-concat';
import copy from 'rollup-plugin-copy';

const debug = process.env.DEBUG;

export default {
  input: 'client/main.ts',
  output: {
    name: 'webtimer',
    dir: 'dist/frontend',
    format: 'iife',
    ...debug ? {
      sourcemap: true,
      sourcemapPathTransform: (relativeSourcePath, sourcemapPath) => {
        return `${relativeSourcePath.replace(new RegExp('^\.\./'), '')}`;
      },
    } : {},
  },
  plugins: [
    typescript({
      'compilerOptions': {
        'sourceMap': debug ? true : false,
        'outDir': 'dist/frontend/build-tsc',
      },
    }),
    commonjs(),
    nodeResolve(),
    ...debug ? [ terser() ] : [],
    concat({
      groupedFiles: [
        {
          outputFile: 'dist/frontend/style.css',
          files: [
            'node_modules/normalize.css/normalize.css',
            'html/style.css',
          ],
        }
      ],
    }),
    copy({
      targets: [
        { src: 'html/index.html', dest: 'dist/frontend' },
        { src: 'html/inter.woff2', dest: 'dist/frontend' },
        ...debug ? [
          {
            src: 'dist/frontend/style.css',
            dest: 'dist/frontend',
            transform: (contents) => new CleanCSS().minify(contents).styles,
          },
        ] : [],
      ]
    }),
  ],
};
