import path from 'path';
import MiniCssExtractPlugin from 'mini-css-extract-plugin'
import CssMinimizerPlugin from 'css-minimizer-webpack-plugin';

const __dirname = import.meta.dirname;

export default {
  mode: 'production',
  context: __dirname,
  entry: [
    './html/index.html',
    './html/style.css',
    './client/main.ts',
  ],
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
        exclude: /node_modules/,
      },
      {
        test: /\.css$/i,
        use: [MiniCssExtractPlugin.loader, 'css-loader'],
      },
      {
        test: /\.(woff|woff2|eot|ttf|otf)$/i,
        type: 'asset/resource',
      },
      {
        test: /\.html$/i,
        type: 'asset/resource',
      },
    ],
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js'],
    extensionAlias: {
      '.js': ['.ts', '.js'],
    }
  },
  optimization: {
    minimizer: [
      new CssMinimizerPlugin(),
      '...',
    ],
  },
  output: {
    filename: 'main.js',
    assetModuleFilename: '[name][ext]',
    path: path.resolve(__dirname, 'dist/frontend'),
    clean: true,
    library: 'webtimer',
  },
  plugins: [
    new MiniCssExtractPlugin({
      filename: 'style.css',
    }),
  ],
};
