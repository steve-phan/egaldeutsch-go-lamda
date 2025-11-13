module.exports = {
  siteMetadata: {
    title: "EgalDeutsch - Learn German Through Stories",
    description:
      "Learn German through engaging stories and interactive quizzes. Practice reading comprehension with content from A1 to C2 level.",
    siteUrl: "https://egaldeutsch.com",
  },
  plugins: [
    "gatsby-plugin-typescript",
    "gatsby-plugin-postcss",
    {
      resolve: "gatsby-plugin-alias-imports",
      options: {
        alias: {
          "@": "src",
          "@/components": "src/components",
          "@/lib": "src/lib",
          "@/utils": "src/utils",
          "@/types": "src/types",
        },
        extensions: ["js", "jsx", "ts", "tsx"],
      },
    },
    // Note: favicon.ico will be served from static folder
  ],
};
