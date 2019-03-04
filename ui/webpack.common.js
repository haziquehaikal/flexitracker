const path = require("path");
const CleanWebpackPlugin = require("clean-webpack-plugin");
const HtmlWebpackPlugin = require("html-webpack-plugin");

const HtmlWebpackPluginConfig = new HtmlWebpackPlugin({
  template: "./public/index.html",
  filename: "index.html",
  inject: "body"
});

module.exports = {
  entry: {
    vendor: ["react", "react-dom", "prop-types"],
    app: "./src/index.js"
  },
  output: {
    path: path.resolve("dist"),
    filename: "index_bundle.js"
  },
  resolve: {
    extensions: [".js", ".jsx"]
  },
  module: {
    rules: [
      {
        test: /\.css$/,
        use: [{ loader: "style-loader" }, { loader: "css-loader" }]
      },
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: "babel-loader"
      },
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        use: "babel-loader"
      }
    ]
  },

  plugins: [HtmlWebpackPluginConfig, new CleanWebpackPlugin(["dist"])],

  output: {
    filename: "[name].[chunkhash:8].js",
    path: path.resolve(__dirname, "dist")
  }
};
