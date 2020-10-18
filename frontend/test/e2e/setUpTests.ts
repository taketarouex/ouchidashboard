import { toMatchImageSnapshot } from "jest-image-snapshot";
import "regenerator-runtime/runtime";

require("expect-puppeteer");

expect.extend({ toMatchImageSnapshot });

jest.setTimeout(50000);