const path = require ('path')
const PurgecssPlugin = require('purgecss-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const glob = require("glob")

const PATHS = {
    src: path.join(__dirname, 'src'),
    dist: path.join(__dirname, 'dist')
}


module.exports = {
    mode: 'production',
    entry: ['./src/index.js', './src/style.scss'],
    output: {
        filename: 'assets/js/main.js',
        path: path.resolve (__dirname, 'dist'),
    },
    module: {
        rules: [
            {
                test: /\.scss$/,
                use: [
                    {
                        loader: MiniCssExtractPlugin.loader
                    },
                    {
                        loader: 'css-loader',
                        options: {
                            url: false,
                        }
                    },
                    {
                        loader: 'postcss-loader'
                    },
                    {
                        loader: 'sass-loader'
                    }
                ]
            }
        ]
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: "assets/css/main.css",
        }),
        new PurgecssPlugin({
            paths: glob.sync(`{${PATHS.dist},${PATHS.src}}/**/*`,  { nodir: true }),
        }),
    ],
    optimization: {
        usedExports: true,
    }
};