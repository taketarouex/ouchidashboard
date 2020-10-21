module.exports = {
  preset: "jest-puppeteer",
  testRegex: "./*\\.e2e\\.ts$",
  setupFilesAfterEnv: ["<rootDir>/setUpTests.ts"],
  transform: {
    "^.+\\.(ts|tsx)$": "ts-jest"
  }
}