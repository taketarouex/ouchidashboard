module.exports = {
  preset: "jest-puppeteer",
  testRegex: "./*\\.e2e\\.ts$",
  setupFilesAfterEnv: ["./setupTests.ts"],
  transform: {
    "^.+\\.(ts|tsx)$": "ts-jest"
  }
}