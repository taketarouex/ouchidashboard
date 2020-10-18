module.exports = {
  preset: "jest-puppeteer",
  testRegex: "./*\\.e2e\\.ts$",
  setupFilesAfterEnv: ["<rootDir>/setupTests.ts"],
  transform: {
    "^.+\\.(ts|tsx)$": "ts-jest"
  }
}